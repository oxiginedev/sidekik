# Sidekik

A useful set of reusable functions across golang projects

## Installation

```bash
go get github.com/oxiginedev/sidekik
```

## Usage

### sidekik.Retry

Retry retries a given function based on the provided options. If no options are provided, it will retry the function once, with a one second delay in-between.

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/oxiginedev/sidekik"
)

func main() {
	err := sidekik.Retry(sidekik.RetryOptions{}, func() error {
		res, err := http.Get("https://catfacts.ninja")
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to get catfacts: %v", err)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
```