package cli

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/y0s3ph/gitops-bootstrap/internal/wizard"
)

var successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#04B575"))
var dimStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#999999"))

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Bootstrap a GitOps workflow on a Kubernetes cluster",
	Long: `Interactive wizard to set up a production-ready GitOps workflow.

Guides you through selecting a GitOps controller, secrets management,
environments, and repository structure. Produces a complete, ready-to-use
GitOps repository and installs the controller on your cluster.`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	cfg, err := wizard.Run()
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(successStyle.Render("✓ Configuration complete"))
	fmt.Println()

	out, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshalling config: %w", err)
	}
	fmt.Println(dimStyle.Render("Generated configuration:"))
	fmt.Println(string(out))

	fmt.Println(dimStyle.Render("Next: scaffolder and installer will use this configuration."))
	fmt.Println(dimStyle.Render("(not yet implemented — coming in Phase 1)"))

	return nil
}
