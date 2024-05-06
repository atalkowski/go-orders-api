package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/dreamsofcode-io/orders-api/application"
)

func main() {
	app := application.New()
	// Set up a handler if we are interrupted by SIGENT (control-C); when interrupted we vill call the cancel method.
	var ctx, cancel = signal.NotifyContext(context.Background(), os.Interrupt) // See https://phk.go.dev/context#Background
	defer cancel()

	err := app.Start(ctx) // Note we removed the temporary context.TODO() see above url.
	if err != nil {
		fmt.Println("failed to start app:", err)
	} else {
		fmt.Println("App started OK")
	}
}
