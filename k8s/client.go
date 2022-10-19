package k8s

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type Client struct {
	clientSet        *kubernetes.Clientset
	metricsClientSet *metrics.Clientset
}

func NewClient() *Client {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	mcs, err := metrics.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &Client{
		clientSet:        cs,
		metricsClientSet: mcs,
	}
}
