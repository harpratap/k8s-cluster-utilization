package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "Calculate resource usage stats",
	Long:  `Uses resource.request to calculate utilization`,
	Run: func(cmd *cobra.Command, args []string) {
		kc.NamespacesResource(context.Background(), exceptionNS, skipBestEffort)
	},
}

var resourceNamespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "Calculate resource usage stats per Namespace",
	Long:  `Uses resource.request to calculate utilization of a namespace`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		nsUtil, err := kc.NamespaceResource(context.Background(), args[0], skipBestEffort)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("NAMESPACE\t\tCPU USED/REQUESTS \tMEMORY USED/REQUESTS\n")
		fmt.Printf("%s\t%.2f/%.2f (%.2f%%)\t%.0f/%.0f (%.2f%%)\n", args[0], nsUtil.UsedCPU, nsUtil.RequestCPU, (nsUtil.UsedCPU / nsUtil.RequestCPU * 100),
			nsUtil.UsedMemory, nsUtil.RequestMemory, (nsUtil.UsedMemory / nsUtil.RequestMemory * 100))
	},
}

func init() {
	rootCmd.AddCommand(resourceCmd)
	resourceCmd.Flags().StringVarP(&exceptionNS, "exception-namespaces", "x", "", "comma separated names of namespaces to exclude")
	resourceCmd.PersistentFlags().BoolVarP(&skipBestEffort, "skip-best-effort-pods", "s", false, "skip pods with no cpu OR memory requests")

	resourceCmd.AddCommand(resourceNamespaceCmd)
}
