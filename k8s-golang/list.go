package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfigPath string
	labelSelector  string
)

type Worker struct {
	ID      string
	StopCh  chan struct{}
	Stopped chan struct{}
}

func init() {
	homeDir, _ := os.LookupEnv("HOME")

	flag.StringVar(&kubeconfigPath, "kubeconfig", fmt.Sprintf("%s/.kube/config", homeDir), "Path to the kubeconfig file")
	flag.StringVar(&labelSelector, "label-selector", "app=myapp", "Label selector for ConfigMap filtering")
	flag.Parse()
}

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Printf("Error building kubeconfig: %v", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %v", err)
		os.Exit(1)
	}

	resp, err := clientset.CoreV1().ConfigMaps("default").List(context.TODO(), metav1.ListOptions{
		LabelSelector: "component=my-config",
	})
	if err != nil {
		fmt.Printf("Error listing ConfigMaps: %v", err)
		os.Exit(1)
	}

	for _, cm := range resp.Items {
		fmt.Println(cm)
	}
}
