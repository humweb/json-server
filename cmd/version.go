package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

var (
	version = "1.0.0"
)

func newVersionCmd() *cobra.Command {
	// versionCmd represents the version command.
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  "Show version information",
		RunE:  runVersion,
	}

	return versionCmd
}

func runVersion(_ *cobra.Command, _ []string) error {
	fmt.Printf("Version:\t %s\n", version)
	fmt.Printf("Go version:\t %s\n", runtime.Version())

	return nil
}
