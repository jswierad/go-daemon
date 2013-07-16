package daemon_test

import (
	"fmt"
	"go-daemon"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func ExampleReborn() {

	err := daemon.Reborn(027, "/")
	if err != nil {
		os.Exit(1)
	}

	signals := make(chan os.Signal, 8)
	signal.Notify(signals, syscall.SIGKILL, syscall.SIGTERM)
	for sig := range signals {
		if sig == syscall.SIGTERM {
			return
		}
	}
}

func ExampleRedirectStream() {

	file, err := os.OpenFile("/tmp/daemon-log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		os.Exit(1)
	}

	if err = daemon.RedirectStream(os.Stdout, file); err != nil {
		os.Exit(2)
	}
	if err = daemon.RedirectStream(os.Stderr, file); err != nil {
		os.Exit(2)
	}
	file.Close()

	fmt.Println("some message")
	log.Println("some message")
}

func ExampleServeSignals() {

	TermHandler := func(sig os.Signal) (stop bool, err error) {
		log.Println("SIGTERM:", sig)
		stop = true
		return
	}

	HupHandler := func(sig os.Signal) (stop bool, err error) {
		log.Println("SIGHUP:", sig)
		stop = false
		return
	}

	daemon.SignalsHandler(TermHandler, syscall.SIGTERM, syscall.SIGKILL)
	daemon.SignalsHandler(HupHandler, syscall.SIGHUP)

	err := daemon.ServeSignals()
	if err != nil {
		log.Println("Error:", err)
	}
}
