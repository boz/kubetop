package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/boz/kubetop/backend"
	"github.com/boz/kubetop/backend/client"
	"github.com/boz/kubetop/ui"
)

func main() {
	client, err := client.NewClient()
	if err != nil {
		panic(err)
	}

	backend := backend.NewBackend(client.Clientset())

	app := ui.NewApp(backend)

	donech := make(chan bool)
	go watchSignals(app, donech)

	if err := app.Run(); err != nil {
		panic(err)
	}

	backend.Stop()
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
