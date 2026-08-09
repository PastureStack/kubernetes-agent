package config

import (
	"fmt"
	"os"

	"github.com/rancher/go-rancher/v3"
	"github.com/urfave/cli"
)

const (
	kuberentesHostEnv = "KUBERNETES_SERVICE_HOST"
	kuberentesPortEnv = "KUBERNETES_SERVICE_PORT"
)

type Config struct {
	KubernetesURL     string
	PlatformURL       string
	PlatformAccessKey string
	PlatformSecretKey string
	WorkerCount       int
	HealthCheckPort   int
	Locale            string
}

func Conf(context *cli.Context) Config {
	kubernetesURL := context.String("kubernetes-url")
	if host, port := os.Getenv(kuberentesHostEnv), os.Getenv(kuberentesPortEnv); host != "" && port != "" {
		kubernetesURL = fmt.Sprintf("https://%s:%s", host, port)
	}
	config := Config{
		KubernetesURL:     kubernetesURL,
		PlatformURL:       context.String("platform-url"),
		PlatformAccessKey: context.String("platform-access-key"),
		PlatformSecretKey: context.String("platform-secret-key"),
		WorkerCount:       context.Int("worker-count"),
		HealthCheckPort:   context.Int("health-check-port"),
		Locale:            context.String("locale"),
	}

	return config
}

func GetPlatformClient(conf Config) (*client.RancherClient, error) {
	return client.NewRancherClient(&client.ClientOpts{
		Url:       conf.PlatformURL,
		AccessKey: conf.PlatformAccessKey,
		SecretKey: conf.PlatformSecretKey,
	})
}
