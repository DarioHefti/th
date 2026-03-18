package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/DarioHefti/th/internal/config"
	"github.com/DarioHefti/th/internal/detect"
	"github.com/DarioHefti/th/internal/llm"
	"github.com/DarioHefti/th/internal/output"
	"github.com/spf13/cobra"
)

var (
	copyToClipboard bool
	configFlag      bool
	Version         = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "th [query]",
	Short: "Get shell commands from an LLM",
	Long: `th (Terminal Help) - Get shell commands using OpenCode Zen free models

Examples:
  th "list all files modified today"
  th "find large files over 100MB"
  th --config  # Re-run setup wizard
`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Usage()
		}

		ctx := context.Background()

		cfg, err := config.Load()
		if err != nil {
			if config.IsConfigNotFound(err) || configFlag {
				if err := runSetupWizard(); err != nil {
					return err
				}
				cfg, err = config.Load()
				if err != nil {
					return fmt.Errorf("loading config after setup: %w", err)
				}
			} else {
				return fmt.Errorf("loading config: %w", err)
			}
		}

		llmClient, err := llm.NewClient(cfg.Endpoint, cfg.Model, "")
		if err != nil {
			return fmt.Errorf("creating LLM client: %w", err)
		}

		env, err := detect.Detect()
		if err != nil {
			return fmt.Errorf("detecting environment: %w", err)
		}

		query := args[0]

		command, err := llmClient.GetCommand(ctx, env.SystemPrompt(), query)
		if err != nil {
			return fmt.Errorf("getting command: %w", err)
		}

		output.PrintCommand(command, copyToClipboard)

		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("th version %s\n", Version)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		output.PrintError(err)
		os.Exit(1)
	}
}

func runSetupWizard() error {
	output.PrintSetupRequired()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("OpenCode Zen - Free Models")
	fmt.Println("Available models:")
	fmt.Println("  - minimax-m2.5-free (MiniMax M2.5)")
	fmt.Println("  - big-pickle (Stealth model)")
	fmt.Println("  - mimo-v2-flash-free (MiMo V2 Flash)")
	fmt.Println("  - nemotron-3-super-free (Nemotron 3 Super)")
	fmt.Println()

	fmt.Print("Model name (press Enter for mimo-v2-flash-free): ")
	model, _ := reader.ReadString('\n')
	model = strings.TrimSpace(model)
	if model == "" {
		model = "mimo-v2-flash-free"
	}

	cfg := &config.Config{
		Provider: "zen",
		Endpoint: "https://opencode.ai/zen/v1",
		Model:    model,
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	output.PrintSuccess(fmt.Sprintf("Configuration saved to %s", config.ConfigPath()))

	return nil
}

func init() {
	rootCmd.Flags().BoolVar(&copyToClipboard, "c", false, "Copy result to clipboard")
	rootCmd.Flags().BoolVar(&configFlag, "config", false, "Run setup wizard")
	rootCmd.AddCommand(versionCmd)
}
