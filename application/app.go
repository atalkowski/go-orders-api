package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
}

func New() *App {
	app := &App{
		rdb: redis.NewClient(&redis.Options{}),
	}
	app.loadRoutes()
	return app
}

// This is defining a method on the "class" App by using the following
func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3001",
		Handler: a.router,
	}
	fmt.Println("Connecting to redis ...")
	err := a.rdb.Ping(ctx).Err() // Start up the redis client connection
	if err != nil {
		return fmt.Errorf("failed to connect to redis service (is it running?) : %w", err)
	}
	// Set up an anonymous function which is deferred ... to handle the redis close call:
	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Println("failed to close redis", err)
		}
	}()

	fmt.Println("Starting server...")
	// Set up a channel to pass the results of the go routine below (I.E. a separte thread.
	ch := make(chan error, 1) // chan, slice and map are the allowed types for make; sets up a buffered channel of size 1
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			// Thread the needle ... pass the result of threaded go routine using the channel we created
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch) // Tells anyone (i.e. App?) listening to this channel that we are done
	}()
	fmt.Println("Server listening on port ", server.Addr)
	// The select below is like a switch case but for process events or something along those lines.
	// So when our channel is done we return
	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		fmt.Println("Shutting down gracefully...")
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout) // I.e. you cannot use the ctx above .. that is now defunct
	}
	//  err, open := <-ch
	// if !open { blah blah  } ... but let's ot do that.
	//return nil // Success!!
}
