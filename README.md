# Sidekik

A useful set of reusable functions across golang projects

## Installation

```bash
go get github.com/oxiginedev/sidekik
```

## Usage

### sidekik.Retry

Retry retries a given function based on the provided options. If no options are provided, it will retry the function once, with a one second delay in-between.
