package k8s

import (
	"context"
	"fmt"
	"os"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var TIMEOUT int64 = 600

type nsResourceUtilization struct {
	Namespace     string
	RequestCPU    float64
	RequestMemory float64
	UsedCPU       float64
	UsedMemory    float64
}

func (c *Client) requestsAndUsageSumPerNamespace(ctx context.Context, nsUtil *nsResourceUtilization, skipBestEffort bool) error {
	podList, err := c.clientSet.CoreV1().Pods(nsUtil.Namespace).List(ctx, metav1.ListOptions{
		Limit:          99999,
		TimeoutSeconds: &TIMEOUT,
		FieldSelector:  "status.phase=Running",
	})
	if err != nil {
		return err
	}
	for _, p := range podList.Items {
		for _, con := range p.Spec.Containers {
			nsUtil.RequestCPU = nsUtil.RequestCPU + con.Resources.Requests.Cpu().AsApproximateFloat64()
			nsUtil.RequestMemory = nsUtil.RequestMemory + con.Resources.Requests.Memory().AsApproximateFloat64()
		}
		err = c.usageSumPerNamespace(ctx, nsUtil, &p, skipBestEffort)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching usage: %s", err.Error())
			continue
		}
	}
	return nil
}

func isContainerCpuRequestZero(containerName string, pod *v1.Pod) bool {
	for _, c := range pod.Spec.Containers {
		if containerName != c.Name {
			continue
		}
		return c.Resources.Requests.Cpu().IsZero()
	}
	return true
}

func isContainerMemoryRequestZero(containerName string, pod *v1.Pod) bool {
	for _, c := range pod.Spec.Containers {
		if containerName != c.Name {
			continue
		}
		return c.Resources.Requests.Memory().IsZero()
	}
	return true
}

func (c *Client) usageSumPerNamespace(ctx context.Context, nsUtil *nsResourceUtilization, pod *v1.Pod, skipBestEffort bool) error {
	podMetrics, err := c.metricsClientSet.MetricsV1beta1().PodMetricses(nsUtil.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	// SkipBestEffort will not add usage for containers
	// that do not specify requests. This will prevent
	// skewing the calculations in case you have too many containers
	// without any requests but they do not require reserving
	// node resources for scheduling
	for _, c := range podMetrics.Containers {
		if !(isContainerCpuRequestZero(c.Name, pod) && skipBestEffort) {
			nsUtil.UsedCPU = nsUtil.UsedCPU + c.Usage.Cpu().AsApproximateFloat64()
		}
		if !(isContainerMemoryRequestZero(c.Name, pod) && skipBestEffort) {
			nsUtil.UsedMemory = nsUtil.UsedMemory + c.Usage.Memory().AsApproximateFloat64()
		}
	}
	return nil
}

func (c *Client) NamespaceResource(ctx context.Context, ns string, skipBestEffort bool) (*nsResourceUtilization, error) {
	nsUtil := nsResourceUtilization{
		Namespace: ns,
	}
	err := c.requestsAndUsageSumPerNamespace(ctx, &nsUtil, skipBestEffort)
	if err != nil {
		return nil, err
	}
	return &nsUtil, nil
}

func (c *Client) NamespacesResource(ctx context.Context, exceptionNS string, skipBestEffort bool) {
	nsList, err := c.clientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	exceptionNS_list := strings.Split(exceptionNS, ",")
	exceptionNS_map := make(map[string]interface{})
	for _, ns := range exceptionNS_list {
		exceptionNS_map[ns] = nil
	}
	fmt.Printf("NAMESPACE,CPU REQUESTS,CPU USED,MEMORY REQUESTS,MEMORY USED\n")
	allNsUtil := nsResourceUtilization{}
	for _, ns := range nsList.Items {
		if _, ok := exceptionNS_map[ns.Name]; ok {
			continue
		}
		nsUtil, err := c.NamespaceResource(ctx, ns.Name, skipBestEffort)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Printf("%s,%.2f,%.2f,%.0f,%.0f\n", ns.Name, nsUtil.RequestCPU, nsUtil.UsedCPU, nsUtil.RequestMemory, nsUtil.UsedMemory)
		allNsUtil.RequestCPU = allNsUtil.RequestCPU + nsUtil.RequestCPU
		allNsUtil.RequestMemory = allNsUtil.RequestMemory + nsUtil.RequestMemory
		allNsUtil.UsedCPU = allNsUtil.UsedCPU + nsUtil.UsedCPU
		allNsUtil.UsedMemory = allNsUtil.UsedMemory + nsUtil.UsedMemory
	}
	fmt.Printf("\n\nTotal Cluster Utilization\n")
	fmt.Printf("Requested CPUs: %.2f, Used CPUs: %.2f (%.2f%%)\n", allNsUtil.RequestCPU, allNsUtil.UsedCPU, (allNsUtil.UsedCPU / allNsUtil.RequestCPU * 100))
	fmt.Printf("Requested Memory: %.2f, Used Memory: %.2f (%.2f%%)\n", allNsUtil.RequestMemory, allNsUtil.UsedMemory, (allNsUtil.UsedMemory / allNsUtil.RequestMemory * 100))
}
