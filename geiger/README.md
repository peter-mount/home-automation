# geiger-counter

This small app runs on a Raspberry PI with a GMC-320+ Geiger Counter connected over USB.

This has been tested only with the GMC-320+ but in theory it should work with other models.

## Installation

### Build

Build the binary for the Raspberry PI.
The following will work on any machine as it will cross-compile to the Arm architecture.

     CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o geiger-counter geiger/bin/main.go

### Configuration

create a `config.yml` file with the following entries:

    mq:
        url: amqp://user:password@hostname
        connectionName: "House Automation Dev"
        product: "Area51 House Automation"
        version: "0.2"

    graphitePublisher:
        exchange: "graphite"
        debug: false
        disabled: false

    geigerPort: "/dev/ttyUSB0"
    geigerPrefix: "home.house.ground.kitchen.geiger"

You will need to ensure that:

* mq.url has the correct entries for your local rabbitMQ instance
* geigerPort is the port the Geiger counter is visible on the PI.

### Installation

Copy the following to the PI:
* geiger-counter binary to /usr/local/bin
* config.yml to /usr/local/etc
* geiger-counter.service to /etc/systemctl/system

You might need to edit `geiger-counter.service` if you have named config.yml differently.

Run:

    sudo systemctl daemon-reload
    sudo systemctl start geiger-counter

If it's running then run the following so that it starts on reboot:

    sudo systemctl enable geiger-counter
