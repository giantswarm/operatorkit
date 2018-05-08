package main

import (
	"flag"
	"fmt"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/operatorkit/example/memcached-operator/memcached"
)

func main() {
	err := mainWithError()
	if err != nil {
		panic(fmt.Sprintf("%#v\n", err))
	}
}

func mainWithError() (err error) {
	c := parseFlags()
	_, err = memcached.New(c)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func parseFlags() memcached.Config {
	var config memcached.Config

	flag.StringVar(&config.K8sAddress, "kubernetes.address", "", "Address used to connect to Kubernetes.")
	flag.BoolVar(&config.K8sInCluster, "kubernetes.incluster", true, "Whether to use the in-cluster config to authenticate with Kubernetes.")
	flag.StringVar(&config.K8sCAFile, "kubernetes.ca", "", "Certificate authority file path to use to authenticate with Kubernetes.")
	flag.StringVar(&config.K8sCrtFile, "kubernetes.crt", "", "Certificate file path to use to authenticate with Kubernetes.")
	flag.StringVar(&config.K8sKeyFile, "kubernetes.key", "", "Key file path to use to authenticate with Kubernetes.")
	flag.Parse()

	return config
}
