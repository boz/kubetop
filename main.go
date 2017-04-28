package main

import (
	"os"
	"os/signal"
	"syscall"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/boz/kubetop/backend"
	"github.com/boz/kubetop/backend/client"
	"github.com/boz/kubetop/ui"
	"github.com/boz/kubetop/util"
)

var (
	logFile = kingpin.Flag("--log-file", "log file output").
		Short('l').
		Default("kubetop.log").
		OpenFile(os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	logLevel = kingpin.Flag("--log-level", "log level").
			Short('v').
			Default("debug").
			Enum("debug", "info", "warn", "error")
)

func main() {
	kingpin.Parse()

	env, err := util.NewEnv(*logFile, *logLevel)
	kingpin.FatalIfError(err, "opening logs")

	env = env.ForComponent("main")

	client, err := client.NewClient(env)
	kingpin.FatalIfError(err, "creating client")

	backend := backend.NewBackend(env, client.Clientset())

	app := ui.NewApp(env, backend)

	donech := make(chan bool)
	go watchSignals(env, app, donech)

	if err := app.Run(); err != nil {
		env.Log().WithError(err).
			Fatal("error running app")
	}

	backend.Stop()
}

func watchSignals(env util.Env, app *ui.App, donech chan bool) {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGQUIT)

	select {
	case <-donech:
	case <-sigch:
		app.Stop()
		<-donech
	}
}
