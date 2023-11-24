package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfigPath string
	labelSelector  string
)

func init() {
	flag.StringVar(&kubeconfigPath, "kubeconfig", "/home/adamw/.kube/config", "Path to the kubeconfig file")
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

	// resp, err := clientset.CoreV1().ConfigMaps("default").List(context.TODO(), metav1.ListOptions{})
	// if err != nil {
	// 	fmt.Printf("Error listing ConfigMaps: %v", err)
	// 	os.Exit(1)
	// }

	// Create an errgroup
	ctx, _ := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	labelOptions := informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
		opts.LabelSelector = "component=my-config"
	})
	sif := informers.NewSharedInformerFactoryWithOptions(clientset, time.Second*10, labelOptions)
	cminf := sif.Core().V1().ConfigMaps()

	inf := cminf.Informer()
	inf.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			cm, _ := obj.(*v1.ConfigMap)
			g.Go(func() error {
				return processWork(ctx, cm)
			})
		},
		UpdateFunc: func(oldCM, newCM interface{}) {
			oldConfigMap := oldCM.(*v1.ConfigMap)
			newConfigMap := newCM.(*v1.ConfigMap)

			oldFieldValue := oldConfigMap.Data["config"]
			newFieldValue := newConfigMap.Data["config"]

			// Only trigger function if config field has changed
			if strings.Compare(oldFieldValue, newFieldValue) != 0 {
				g.Go(func() error {
					return processWork(ctx, newConfigMap)
				})
			}
		},
		DeleteFunc: func(obj interface{}) {
			// cm, _ := obj.(*v1.ConfigMap)
			// outputConfigMap(cm)
			fmt.Println("removing capacity")
		},
	})

	g.Go(func() error {
		go inf.Run(ctx.Done())
		<-ctx.Done()
		return nil
	})

	// g.Go(func() error {
	// 	return triggerError()
	// })

	// Print information about each ConfigMap
	// fmt.Printf("ConfigMaps in namespace %s:\n", "default")
	// for _, configMap := range resp.Items {
	// 	fmt.Printf("Name: %s\n", configMap.Name)
	// 	fmt.Printf("Namespace: %s\n", configMap.Namespace)
	// 	fmt.Printf("Creation Timestamp: %v\n", configMap.CreationTimestamp.Time)
	// 	fmt.Println("Data:")
	// 	for key, value := range configMap.Data {
	// 		fmt.Printf("  %s: %s\n", key, value)
	// 	}
	// 	fmt.Println("-----")
	// }
	// Wait for either the informer or the errgroup function to return
	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Finished")
}

func outputConfigMap(cm *v1.ConfigMap) {
	fmt.Printf("Name: %s\n", cm.Name)
	fmt.Printf("Namespace: %s\n", cm.Namespace)
	fmt.Printf("Creation Timestamp: %v\n", cm.CreationTimestamp.Time)
	fmt.Println("-----")
}

func triggerError() error {
	time.Sleep(15 * time.Second)

	return errors.New("Program failed")
}

func processWork(ctx context.Context, cm *v1.ConfigMap) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	fmt.Printf("Processing work for %s\n", cm.Name)

	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(10)

	fmt.Println(randomNumber)
	if randomNumber > 6 {
		return errors.New("Oh no a random error occured")
	}

	// outputConfigMap(cm)
	return nil
}
