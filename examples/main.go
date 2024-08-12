package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/bugsnag/bugsnag-go/v2"
	"os"
	"time"
	"werror"
	// If you are using this library in your own project,
	// make sure to import it using the full module path:
	// "github.com/kanmo/werror"
)

func main() {
	// Set your API key to Environment variable
	apiKey := os.Getenv("BUGSNAG_API_KEY")
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
