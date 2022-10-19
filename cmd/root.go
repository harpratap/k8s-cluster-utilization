package cmd

import (
	"os"

	"github.com/harpratap/k8s-cluster-utilization/k8s"
	"github.com/spf13/cobra"
)

var (
	kc             *k8s.Client
	exceptionNS    string
	skipBestEffort bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8s-cluster-utilization",
	Short: "Calculate stats for your K8s cluster usage",
	Long:  `Calculate stats for your K8s cluster usage`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		kc = k8s.NewClient()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.k8s-cluster-utilization.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
