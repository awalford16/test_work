package main

import (
	"context"
	"errors"
	"fmt"

	v1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

// KubernetesClient implements ConfigMapOperations using the actual Kubernetes API.
type KubernetesClient struct {
	Clientset *kubernetes.Clientset
}

func main() {
	// Example using the actual Kubernetes API
	kubeClient := fake.NewSimpleClientset()

	// Creating a ConfigMap with a label
	defer tmpCreateConfigMap(kubeClient, "test-config")()

	cm, err := kubeClient.CoreV1().ConfigMaps("default").Get(context.TODO(), "test-config", metav1.GetOptions{})
	if err != nil {
		return
	}

	fmt.Println("ConfigMaps with label selector Data:", cm.Data)
}

// Create a config map while the program is running and then delete it
func tmpCreateConfigMap(kubeClient kubernetes.Interface, name string) func() {
	_, err := kubeClient.CoreV1().ConfigMaps("default").Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		// Create the config map if it does not already exist
		if k8sErrors.IsNotFound(err) {
			fmt.Println("Creating a new configmap")

			_, _ = kubeClient.CoreV1().ConfigMaps("default").Create(context.TODO(), &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-config",
					Namespace: "default",
					Labels:    map[string]string{"app": "myapp"},
				},
				Data: map[string]string{"key": "value"},
			}, metav1.CreateOptions{})
		} else {
			errors.New("Failed to get Configmap")
		}
	}

	return func() {
		fmt.Printf("Deleting CM %s\n", name)
		kubeClient.CoreV1().ConfigMaps("default").Delete(context.Background(), name, metav1.DeleteOptions{})
	}
}
