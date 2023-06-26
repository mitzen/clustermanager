package controller

import (
	"fmt"

	"k8s.io/client-go/rest"
)

type ClientConfig struct {
}

func (c *ClientConfig) GetConfig() *rest.Config {
	// InClusterConfig
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Print("Error getting client config")
		//r.log.Error(err, "Fault in rest.InClusterConfig")
	}

	// home, _ := os.UserHomeDir()
	// kubeConfigPath := filepath.Join(home, ".kube", "config")
	// config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	// if err != nil {
	// 	panic(err)
	// }

	return config
}
