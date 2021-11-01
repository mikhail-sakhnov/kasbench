package main

import (
	"fmt"
	"log"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/scheme"
)

func PrintGETTargets(clientset kubernetes.Interface) {
	newReq := func() *rest.Request {
		return clientset.CoreV1().RESTClient().Get().
			VersionedParams(&metav1.GetOptions{}, scheme.ParameterCodec).
			Timeout(time.Second)

	}
	reqs := []*rest.Request{
		newReq().Namespace("kube-system").Resource("services"),
		newReq().Namespace("kube-system").Resource("pods"),
		newReq().Resource("nodes"),
		newReq().Resource("persistentvolumes"),
		newReq().Namespace("kube-system").Resource("configmaps"),
	}
	for _, req := range reqs {
		fmt.Println("GET", req.URL())
	}
}

func main() {
	// Create an InClusterConfig and use it to create a client for the controller
	// to use to communicate with Kubernetes
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))

	if err != nil {
		log.Fatalf("Failed to create config: %v", err)
	}
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err)
	}

	PrintGETTargets(clientset)

}
