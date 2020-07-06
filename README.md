## Premise

We want to create 2 micro services, both http services and both listen to `/hello`:

- service 1
  - `/hello` => `1 here`

- service 2
  - `/hello` => `2 here`

Is it possible with go-micro to do such thing?

Answer: it appears to be possible.

## Run this:

- go run main.go --msg=1

on another session:
- go run main.go --msg=2