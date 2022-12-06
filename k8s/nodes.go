package k8s

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodeUtilization struct {
	TotalCPU          float64
	TotalMemory       float64
	AllocatableCPU    float64
	AllocatableMemory float64
	RequestedCPU      float64
	RequestedMemory   float64
	UsedCPU           float64
	UsedMemory        float64
	AllocatablePods   int
	RunningPods       int
	DaemonsetCPU      float64
	DaemonsetMemory   float64
	WorkloadCPU       float64
	WorkloadMemory    float64
}

func isPodDaemonset(pod *v1.Pod) bool {
	for _, v := range pod.OwnerReferences {
		return v.Kind == "DaemonSet"
	}
	return false
}

func (c *Client) NodeUsage(ctx context.Context, name string) (*NodeUtilization, error) {
	node, err := c.clientSet.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	nUtil := NodeUtilization{
		AllocatableCPU:    node.Status.Allocatable.Cpu().AsApproximateFloat64(),
		AllocatableMemory: node.Status.Allocatable.Memory().AsApproximateFloat64(),
		AllocatablePods:   int(node.Status.Allocatable.Pods().AsApproximateFloat64()),
		TotalCPU:          node.Status.Capacity.Cpu().AsApproximateFloat64(),
		TotalMemory:       node.Status.Capacity.Memory().AsApproximateFloat64(),
	}

	pods, err := c.clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", name),
	})
	if err != nil {
		return nil, err
	}

	for _, pod := range pods.Items {
		nUtil.RunningPods++
		for _, container := range pod.Spec.Containers {
			nUtil.RequestedCPU += container.Resources.Requests.Cpu().AsApproximateFloat64()
			nUtil.RequestedMemory += container.Resources.Requests.Memory().AsApproximateFloat64()
			if isPodDaemonset(&pod) {
				nUtil.DaemonsetCPU += container.Resources.Requests.Cpu().AsApproximateFloat64()
				nUtil.DaemonsetMemory += container.Resources.Requests.Memory().AsApproximateFloat64()
			} else {
				nUtil.WorkloadCPU += container.Resources.Requests.Cpu().AsApproximateFloat64()
				nUtil.WorkloadMemory += container.Resources.Requests.Memory().AsApproximateFloat64()
			}
		}
	}
	return &nUtil, nil
}
