package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/PastureStack/kubernetes-agent/config"
	"github.com/PastureStack/kubernetes-agent/healthcheck"
	"github.com/PastureStack/kubernetes-agent/hostwatch"
	"github.com/PastureStack/kubernetes-agent/kubernetesclient"
	"github.com/PastureStack/kubernetes-agent/platformevents"
	"github.com/PastureStack/kubernetes-agent/truststore"
	"github.com/PastureStack/kubernetes-agent/watchevents"
)

var VERSION = "0.0.0"

func main() {
	app := cli.NewApp()
	app.Name = "kubernetes-agent"
	app.Usage = "Start the PastureStack Kubernetes compatibility agent"
	app.Version = VERSION
	app.Action = launch

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "kubernetes-url",
			Value:  "http://localhost:8080",
			Usage:  "URL for kubernetes API",
			EnvVar: "KUBERNETES_URL",
		},
		cli.StringFlag{
			Name:   "platform-url,cattle-url",
			Usage:  "URL for the control-platform API",
			EnvVar: "PLATFORM_URL,CATTLE_URL",
		},
		cli.StringFlag{
			Name:   "platform-access-key,cattle-access-key",
			Usage:  "Control-platform API access key",
			EnvVar: "PLATFORM_ACCESS_KEY,CATTLE_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "platform-secret-key,cattle-secret-key",
			Usage:  "Control-platform API secret key",
			EnvVar: "PLATFORM_SECRET_KEY,CATTLE_SECRET_KEY",
		},
		cli.IntFlag{
			Name:   "worker-count",
			Value:  50,
			Usage:  "Number of workers for handling events",
			EnvVar: "WORKER_COUNT",
		},
		cli.StringFlag{
			Name:   "locale",
			Value:  "en-US",
			Usage:  "Operator message locale: en-US or zh-TW",
			EnvVar: "PASTURESTACK_LOCALE",
		},
		cli.IntFlag{
			Name:   "health-check-port",
			Value:  10240,
			Usage:  "Port to configure an HTTP health check listener on",
			EnvVar: "HEALTH_CHECK_PORT",
		},
		cli.IntFlag{
			Name:  "host-update-interval",
			Value: 5,
			Usage: "The frequency at which host labels should be updated",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func launch(c *cli.Context) error {
	conf := config.Conf(c)

	resultChan := make(chan error)

	if conf.Locale != "en-US" && conf.Locale != "zh-TW" {
		return fmt.Errorf("unsupported locale %q; use en-US or zh-TW", conf.Locale)
	}
	customCAConfigured, err := truststore.ConfigurePlatformCA()
	if err != nil {
		return err
	}
	if customCAConfigured {
		log.Info("Using the mounted control-platform CA certificate.")
	}

	platformClient, err := config.GetPlatformClient(conf)
	if err != nil {
		return fmt.Errorf("initialize control-platform API client: %w", err)
	}

	kClient, err := kubernetesclient.NewClient(conf.KubernetesURL)
	if err != nil {
		return err
	}

	svcHandler := watchevents.NewServiceHandler(platformClient, kClient)

	nsHandler := watchevents.NewNamespaceHandler(platformClient, kClient)

	if err := svcHandler.Start(); err != nil {
		return fmt.Errorf("start Kubernetes service watch: %w", err)
	}
	defer svcHandler.Stop()

	if err := nsHandler.Start(); err != nil {
		return fmt.Errorf("start Kubernetes namespace watch: %w", err)
	}
	defer nsHandler.Stop()

	go func(rc chan error) {
		err := platformevents.ConnectToEventStream(platformClient, kClient, conf)
		log.Errorf("%s: %s", operatorMessage(conf.Locale, "stream-exit"), err)
		rc <- err
	}(resultChan)

	go func(rc chan error) {
		err := healthcheck.StartHealthCheck(conf.HealthCheckPort)
		log.Errorf("%s: %s", operatorMessage(conf.Locale, "health-exit"), err)
		rc <- err
	}(resultChan)

	go func(rc chan error) {
		err := hostwatch.StartHostSync(c.Int("host-update-interval"), kClient)
		log.Errorf("%s: %s", operatorMessage(conf.Locale, "host-sync-exit"), err)
		rc <- err
	}(resultChan)

	err = <-resultChan
	log.Info(operatorMessage(conf.Locale, "exit"))
	return err
}

func operatorMessage(locale, key string) string {
	messages := map[string]map[string]string{
		"en-US": {
			"stream-exit":    "Platform stream listener exited with an error",
			"health-exit":    "Health check exited with an error",
			"host-sync-exit": "Host synchronization exited with an error",
			"exit":           "Exiting.",
		},
		"zh-TW": {
			"stream-exit":    "平台事件串流因錯誤而停止",
			"health-exit":    "健康檢查因錯誤而停止",
			"host-sync-exit": "主機同步因錯誤而停止",
			"exit":           "正在結束。",
		},
	}
	return messages[locale][key]
}
