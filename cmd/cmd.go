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
		metricPort, _ := cmd.Flags().GetInt("metric_port")
		if metricPort == 0 {
			metricPort = 4200
		}
		sync(metricPort)
	},
}

var processTasksCmd = &cobra.Command{
	Use:   "process-task",
	Short: "Processing ",
	Run: func(cmd *cobra.Command, args []string) {
		metricPort, _ := cmd.Flags().GetInt("metric_port")
		if metricPort == 0 {
			metricPort = 4201
		}
		performTasks(metricPort)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
	syncCmd.Flags().Int("metric_port", 4200, "metric_port")
	RootCmd.AddCommand(syncCmd)
	processTasksCmd.Flags().Int("metric_port", 4201, "metric_port")
	RootCmd.AddCommand(processTasksCmd)
}
