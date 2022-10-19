# k8s-cluster-utilization

This tool will help you aggregate resource related data from your cluster
- Total sum of all CPU and Memory requests of pods in a namespace
- Total sum of all CPU and Memory requests in a cluster (can provide exception list for system namespaces like `kube-system`). It will be printed in CSV format so you can load it in Excel and run some calculations on top of it
- Total usage of all CPU and Memory requests of pods in a namespace (relies on metrics.k8s.io API, same as `kubectl top pods`)
- Total usage of all CPU and Memory requests
- `--skip-best-effort` flag allows you to remove pods that do not specify any CPU or Memory, this is useful because such pods may not affect your scheduling and cluster size much but end up adding noise to your Utilization ratio

Try out using
```
go build .
./k8s-cluster-utilization help
```
