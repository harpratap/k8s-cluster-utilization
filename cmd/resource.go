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

var resourceNodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Calculate resource usage stats per Node",
	Long:  `Segregates info on Allocatable, daemonsets and pods CPU requests and usage`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		nUtil, err := kc.NodeUsage(context.Background(), args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf(`Total CPU: %.2f
Total Memory: %.2f
Allocatable CPU: %.2f (%.2f%%)
Allocatable Memory: %.2f (%.2f%%)
Requested CPU: %.2f (%.2f%%)
Requested Memory: %.2f (%.2f%%)
Used CPU:
Used Memory:
Daemonset CPU: %.2f (%.2f%%)
Daemonset Memory: %.2f (%.2f%%)
Workload CPU: %.2f (%.2f%%)
Workload Memory: %.2f (%.2f%%)
Allocatable Pods: %d
Running Pods: %d (%d%%)
`,
			nUtil.TotalCPU,
			nUtil.TotalMemory,
			nUtil.AllocatableCPU, nUtil.AllocatableCPU*100/nUtil.TotalCPU,
			nUtil.AllocatableMemory, nUtil.AllocatableMemory*100/nUtil.TotalMemory,
			nUtil.RequestedCPU, nUtil.RequestedCPU*100/nUtil.AllocatableCPU,
			nUtil.RequestedMemory, nUtil.RequestedMemory*100/nUtil.AllocatableMemory,
			nUtil.DaemonsetCPU, nUtil.DaemonsetCPU*100/nUtil.AllocatableCPU,
			nUtil.DaemonsetMemory, nUtil.DaemonsetMemory*100/nUtil.AllocatableMemory,
			nUtil.WorkloadCPU, nUtil.WorkloadCPU*100/nUtil.AllocatableCPU,
			nUtil.WorkloadMemory, nUtil.WorkloadMemory*100/nUtil.AllocatableMemory,
			nUtil.AllocatablePods,
			nUtil.RunningPods, nUtil.RunningPods*100/nUtil.AllocatablePods)
	},
}

func init() {
	rootCmd.AddCommand(resourceCmd)
	resourceCmd.Flags().StringVarP(&exceptionNS, "exception-namespaces", "x", "", "comma separated names of namespaces to exclude")
	resourceCmd.PersistentFlags().BoolVarP(&skipBestEffort, "skip-best-effort-pods", "s", false, "skip pods with no cpu OR memory requests")

	resourceCmd.AddCommand(resourceNamespaceCmd)
	resourceCmd.AddCommand(resourceNodeCmd)
}
