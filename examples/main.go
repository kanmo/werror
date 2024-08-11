package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/bugsnag/bugsnag-go/v2"
	"time"
	"werror"
	// If you are using this library in your own project,
	// make sure to import it using the full module path:
	// "github.com/kanmo/werror"
)

// Insert your API key
const apiKey = "Your API Jey Comes Here"

func main() {
	if len(apiKey) != 32 {
		fmt.Println("Please set your API key in main.go before running example.")
		return
	}

	bugsnag.Configure(
		bugsnag.Configuration{
			APIKey:       apiKey,
			ReleaseStage: "production",
			//ProjectPackages: []string{"your-project-package"},
		})

	werr := createWrappedError()
	err := bugsnag.Notify(werr, context.Background())
	if err != nil {
		fmt.Println(err)
	}

	// Wait for the error to be sent
	time.Sleep(2 * time.Second)

}

func createWrappedError() error {
	return werror.Wrap(errors.New("wrapped error"))
}
