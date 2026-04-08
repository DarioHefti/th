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
	Args: validateRootArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if configFlag {
			return runSetupWizard()
		}

		if len(args) == 0 {
			return cmd.Usage()
		}

		ctx := context.Background()

		cfg, err := config.Load()
		if err != nil {
			if config.IsConfigNotFound(err) {
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
	fmt.Println("  - minimax-m2.5 (MiniMax M2.5)")
	fmt.Println("  - big-pickle (Stealth model)")
	fmt.Println("  - nemotron-3-super-free (Nemotron 3 Super)")
	fmt.Println()

	fmt.Print("Model name (press Enter for minimax-m2.5): ")
	model, _ := reader.ReadString('\n')
	model = strings.TrimSpace(model)
	if model == "" {
		model = "big-pickle"
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

func validateRootArgs(cmd *cobra.Command, args []string) error {
	cfgMode, err := cmd.Flags().GetBool("config")
	if err != nil {
		return err
	}

	if cfgMode {
		if len(args) > 0 {
			return fmt.Errorf("--config does not accept a query argument")
		}
		return nil
	}

	return cobra.MaximumNArgs(1)(cmd, args)
}

func init() {
	rootCmd.Flags().BoolVar(&copyToClipboard, "c", false, "Copy result to clipboard")
	rootCmd.Flags().BoolVar(&configFlag, "config", false, "Run setup wizard")
	rootCmd.AddCommand(versionCmd)
}
