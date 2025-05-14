package kub_api

import (
	"context"
	"flag"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubAPI struct {
	Kubeconfig *string
}

func (kapi *KubAPI) List() {
	namespace := flag.String("namespace", "default", "namespace to list pods in")
	flag.Parse()

	// Use in-cluster config if running inside a pod, otherwise use kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", *kapi.Kubeconfig)
	if err != nil {
		fmt.Printf("Error building kubeconfig: %v\n", err)
		//config, err = clientcmd.InClusterConfig()
		if err != nil {
			fmt.Printf("Error building in-cluster config: %v\n", err)
			os.Exit(1)
		}
	}

	// Create a Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating clientset: %v\n", err)
		os.Exit(1)
	}

	// List pods in the specified namespace
	pods, err := clientset.CoreV1().Pods(*namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing pods in namespace '%s': %v\n", *namespace, err)
		os.Exit(1)
	}

	fmt.Printf("Pods in namespace '%s':\n", *namespace)
	for _, pod := range pods.Items {
		fmt.Printf("- Name: %s, Status: %s\n", pod.Name, pod.Status.Phase)
	}
}
