package k8s

import "context"

type NodeUtilization struct {
	AllocatableCPU    float64
	AllocatableMemory float64
	RequestedCPU      float64
	RequestedMemory   float64
	UsedCPU           float64
	UsedMemory        float64
	PodsPerNode       int
}

func (c *Client) NodeUsage(ctx context.Context, name string) {

}
