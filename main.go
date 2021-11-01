package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/scheme"
)

const caCertPath = "ca.cert"
const certPath = "cert.cert"
const certKeyPath = "cert.key"

func PrintGETTargets(file io.Writer, clientset kubernetes.Interface) error {
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
		if _, err := fmt.Fprintf(file, "%s %s\n", "GET", req.URL()); err != nil {
			return err
		}
	}

	return nil
}

func generateSuiteForKubeconfig(config *rest.Config) error {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("can't create k8s client: %w", err)
	}
	tmpDir, err := os.MkdirTemp(".", "bench")
	if err != nil {
		return fmt.Errorf("can't init tmp dir: %w", err)
	}

	if err := saveCertificates(tmpDir, config.CAData, config.CertData, config.KeyData); err != nil {
		return fmt.Errorf("can't save certificates: %w", err)
	}

	if err := generateTargets(clientset, tmpDir, "GET"); err != nil {
		return fmt.Errorf("can't prepare targets: %w", err)
	}
	fmt.Println("Suite generated")
	fmt.Printf("Manual run: \n%s\n", generateRunCommand(tmpDir))
	fmt.Println(tmpDir)
	return nil
}

func generateRunCommand(pathToSuite string) string {
	return fmt.Sprintf("vegeta attack -duration=60s -cert %s/%s -key %s/%s -root-certs=%s/%s -targets %s/%s -output %s/report",
		pathToSuite,
		certPath,
		pathToSuite,
		certKeyPath,
		pathToSuite,
		caCertPath,
		pathToSuite,
		targetFilePath,
		pathToSuite,
	)
}

const targetFilePath = "targets"

func generateTargets(clientset kubernetes.Interface, dir string, method string) error {
	targetsFile, err := os.Create(filepath.Join(dir, targetFilePath))
	if err != nil {
		return err
	}
	switch method {
	case "GET":
		return PrintGETTargets(targetsFile, clientset)
	default:
		panic("Unknown method: " + method)
	}
}

func saveCertificates(dir string, caData, certData, keyData []byte) error {
	if err := os.WriteFile(filepath.Join(dir, caCertPath), caData, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, certPath), certData, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, certKeyPath), keyData, 0755); err != nil {
		return err
	}
	return nil
}

func main() {

	var kubeconfig string
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to target cluster kubeconfig")

	flag.Parse()
	if kubeconfig == "" {
		panic("-kubeconfig argument is required")
	}
	// Create an InClusterConfig and use it to create a client for the controller
	// to use to communicate with Kubernetes
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	// config.
	if err != nil {
		log.Fatalf("Failed to create config: %v", err)
	}

	generateSuiteForKubeconfig(config)

}
