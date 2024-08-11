# werror is Custom Error Library with Bugsnag Integration

This library provides a customizable error type that integrates with Bugsnag for capturing and reporting errors with detailed stack traces. It implements the ErrorWithCallers interface required by Bugsnag, allowing you to easily trace errors back to their origin.

## Features

	•	Stack Trace Integration: Errors created with this library automatically include stack traces that can be sent to Bugsnag for easy debugging.
	•	Error Annotation: You can add additional context to your errors using WithCode and WithReason functions.
	•	Selective Reporting: Control which errors get reported to Bugsnag using the WithIgnoreReport function.

## Usage

### Creating an Error

To create a custom error that includes a stack trace, use the NewError function provided by the library.

```go
import "github.com/kanmo/werror"

// Create a basic error
err := werror.New("Something went wrong")
```

### Annotating Errors

You can annotate your errors with additional context, such as a code or a reason.  
reason needs implementation of Stringer interface.

```go
err := errors.New("Something went wrong")
err = werror.Wrap(err, WithCode(codes.InvalidArgument), WithReason("The operation is not allowed"))
```

### Ignoring Errors for Bugsnag Reporting

If you want to prevent certain errors from being reported to Bugsnag, use the WithIgnoreReport function.

```go
err := errors.New("Something went wrong")
err = werror.Wrap(err, WithIgnoreReport())
```

### Integration with Bugsnag

This library is designed to integrate seamlessly with Bugsnag. The custom error type implements the ErrorWithCallers interface, allowing Bugsnag to capture and report stack traces.

Simply pass the error to Bugsnag as usual:
(see: https://github.com/bugsnag/bugsnag-go) 

```go
import (
	"context"
	"github.com/bugsnag/bugsnag-go/v2"
)	

bugsnag.Notify(err, context.Background())
```

## Example

Here’s a full example of how to use the library:

```go
package main

import (
	"context"
	"github.com/kanmo/werror"
	"github.com/bugsnag/bugsnag-go/v2"
	"google.golang.org/grpc/codes"
)

func main() {
	// Initialize Bugsnag
	bugsnag.Configure(bugsnag.Configuration{
		APIKey: "your-api-key",
	})

	// Create an annotated error
	err := errors.New("Something went wrong")
	err = werror.Wrap(err,
		WithCode(codes.InvalidArgument),
		WithReason(pb.ErrorReason_INVALID_REQUEST), // your application specific reason
	)

	// Notify Bugsnag of the error
	bugsnag.Notify(err, context.Background())
}
```

## License

This library is licensed under the MIT License. See the [LICENSE.txt](LICENSE.txt) file for details.
```