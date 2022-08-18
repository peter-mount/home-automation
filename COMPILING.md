# Compiling

Before you build any application run the following command to ensure you have all the dependencies
available to you:

    go mod download

To compile each service you can use one of the following commands:

| application     | go build                                              |
|-----------------|-------------------------------------------------------|
| automation      | go build -o automation-service automation/bin/main.go |
| geiger          | go build -o geiger-counter geiger/bin/main.go         |
| pimoroni-enviro | go build -o pimoroni-enviro enviro/bin/main.go        |

## For the Raspberry PI

You don't need to have the go development tools installed on your PI.
You can cross compile from any machine that has the development tools installed.

To compile for the Raspberry PI 4B with the default Raspberry PI OS installed prefix the `go build` command with the following:

    CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7

For example, to compile `geiger-counter` for the PI use:

    CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o geiger-counter geiger/bin/main.go

Then copy the binary to the PI.
