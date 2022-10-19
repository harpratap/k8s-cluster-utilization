package k8s

// func (c *Client) cpuRequestSumPerNamespace() {
// 	cidrCollection := map[string]int{}
// 	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	for _, node := range nodes.Items {
// 		// pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
// 		// 	FieldSelector: fmt.Sprintf("spec.nodeName=%s,status.phase=Running", node.Name),
// 		// })
// 		// if err != nil {
// 		// 	panic(err.Error())
// 		// }
// 		cidr := strings.Split(node.Spec.PodCIDR, "/")[1]
// 		cidrCollection[cidr] += 1
// 		fmt.Printf("%s - /%s : %d\n", node.Name, cidr, cidrCollection[cidr])
// 	}
// 	fmt.Println(cidrCollection)
// 	// for k, v := range cidrCollection {
// 	// 	d := stats.LoadRawData(v)
// 	// 	mean, err := stats.Mean(d)
// 	// 	if err != nil {
// 	// 		panic(err.Error())
// 	// 	}
// 	// 	median, err := stats.Median(d)
// 	// 	if err != nil {
// 	// 		panic(err.Error())
// 	// 	}
// 	// 	fmt.Printf("/%s, %f, %f\n", k, mean, median)
// 	// }
// }
