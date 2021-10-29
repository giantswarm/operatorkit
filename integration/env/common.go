package env

import (
	"fmt"
	"os"
)

const (
	EnvVarE2EKubeconfig = "E2E_KUBECONFIG"
)

var (
	kubeconfig string
)

func init() {
	kubeconfig = os.Getenv(EnvVarE2EKubeconfig)
	if kubeconfig == "" {
		panic(fmt.Sprintf("env var %#q must not be empty", EnvVarE2EKubeconfig))
	}
}

func KubeConfigPath() string {
	return kubeconfig
}
