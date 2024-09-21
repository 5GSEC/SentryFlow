package k8s

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewClient returns a new Client using the provided scheme to map go structs to
// GroupVersionKinds.
func NewClient(scheme *runtime.Scheme, kubeConfig string) (client.Client, error) {
	config, err := getConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %v", err)
	}
	return client.New(config, client.Options{
		Scheme: scheme,
	})
}

func getConfig(kubeConfig string) (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil && errors.Is(err, rest.ErrNotInCluster) {
		if kubeConfig == "" {
			kubeConfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}
