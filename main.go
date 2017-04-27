package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/boz/kubetop/ui"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, err
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	).ClientConfig()

	return config, err
}

func getClient() (*kubernetes.Clientset, *rest.Config, error) {
	config, err := getConfig()
	if err != nil {
		return nil, nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	return clientset, config, nil
}

func main() {

	/*
		_, _, err := getClient()
		if err != nil {
			panic(err)
		}
	*/

	app := ui.NewApp()

	donech := make(chan bool)
	go watchSignals(app, donech)

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func watchSignals(app *ui.App, donech chan bool) {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGQUIT)
	select {
	case <-donech:
	case <-sigch:
		app.Stop()
		<-donech
	}
}
