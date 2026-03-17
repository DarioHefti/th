package cmd

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/terminal-help/th/internal/auth"
	"github.com/terminal-help/th/internal/config"
	"github.com/terminal-help/th/internal/detect"
	"github.com/terminal-help/th/internal/llm"
	"github.com/terminal-help/th/internal/output"
)

var (
	copyToClipboard bool
	configFlag      bool
	Version         = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "th [query]",
	Short: "Get shell commands from an LLM",
	Long: `th (Terminal Help) - Get shell commands from Azure AI Foundry

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

		azureAuth, err := auth.NewAzureAuth(cfg.TenantID, cfg.ClientID)
		if err != nil {
			return fmt.Errorf("creating auth: %w", err)
		}

		output.PrintAuthPrompt()

		llmClient, err := llm.NewClient(cfg.Endpoint, cfg.Deployment, cfg.APIVersion, azureAuth.GetToken)
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

	fmt.Print("Azure Tenant ID (press Enter for common tenants): ")
	tenantID, _ := reader.ReadString('\n')
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		tenantID = "common"
	}

	fmt.Print("Azure Client ID (app registration): ")
	clientID, _ := reader.ReadString('\n')
	clientID = strings.TrimSpace(clientID)
	if clientID == "" {
		return fmt.Errorf("client ID is required")
	}

	fmt.Print("Azure AI Foundry Endpoint (e.g., https://my-resource.openai.azure.com/): ")
	endpoint, _ := reader.ReadString('\n')
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return fmt.Errorf("endpoint is required")
	}
	parsedURL, err := url.Parse(endpoint)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("invalid endpoint URL: must be a valid URL (e.g., https://my-resource.openai.azure.com/)")
	}
	if parsedURL.Scheme != "https" {
		return fmt.Errorf("endpoint must use HTTPS")
	}

	fmt.Print("Deployment/Model name (e.g., gpt-4o): ")
	deployment, _ := reader.ReadString('\n')
	deployment = strings.TrimSpace(deployment)
	if deployment == "" {
		return fmt.Errorf("deployment name is required")
	}

	fmt.Print("API Version (press Enter for 2024-02-15-preview): ")
	apiVersion, _ := reader.ReadString('\n')
	apiVersion = strings.TrimSpace(apiVersion)
	if apiVersion == "" {
		apiVersion = "2024-02-15-preview"
	}

	cfg := &config.Config{
		TenantID:   tenantID,
		ClientID:   clientID,
		Endpoint:   endpoint,
		Deployment: deployment,
		APIVersion: apiVersion,
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	output.PrintSuccess(fmt.Sprintf("Configuration saved to %s", config.ConfigPath()))

	return nil
}

func init() {
	rootCmd.Flags().BoolVar(&copyToClipboard, "clipboard", false, "Copy result to clipboard")
	rootCmd.Flags().BoolVar(&configFlag, "config", false, "Run setup wizard")
	rootCmd.AddCommand(versionCmd)
}
