package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version string
	commit  string
)

// serveCmd represents the serve command
var RootCmd = &cobra.Command{
	Use:   "",
	Short: "",
	Long:  "",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print version of corev3-explorer",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("VERSION : %s - COMMIT : %s\n", version, commit)
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync block from blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		sync()
	},
}

var processTasks = &cobra.Command{
	Use:   "process-task",
	Short: "Processing ",
	Run: func(cmd *cobra.Command, args []string) {
		sync()
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(syncCmd)
	RootCmd.AddCommand(processTasks)
}
