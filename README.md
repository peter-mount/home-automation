# home-automation

This repository contains the backend for my Home Automation which run's on a Raspberry PI 4B.

## Overview

The automation consists of a collection of Zigbee and WiFi devices which feed into the system via MQTT.

### Applications in this Repository

* **geiger-counter** is a small service which allows a Geiger counter to feed current radiation levels into the system.
* **pimoroni-enviro** is a small bridge to import the metrics from the Pimoroni Enviro suite of Pico-W based sensors
  into
  the system
    * [Enviro-Urban](https://shop.pimoroni.com/products/enviro-urban?variant=40056508252243) for Temperature, Humidity, Noise and Particle pollution readings
    * [Enviro-Weather](https://shop.pimoroni.com/products/enviro-weather?variant=40056776917075) for Temperature, Humidity, Pressure, Wind, Rain & Light levels

The following exist but have yet to be ported to this repository as they require some cleaning up first:

* **automation** is the core of the system. This allows switches or PIR sensors to trigger scenes like turning on or off
  lights, smoke detectors etc.
* **graphite-bridge** takes the raw metrics received from zigbee sensors and transforms them into a form that can be
  imported into Graphite/Carbon

### Third Party Applications

* [RabbitMQ](https://www.rabbitmq.com/) is the MQTT message broker
* [zigbee2mqtt](https://www.zigbee2mqtt.io/) is used to bridge the Zigbee network to MQTT

Optional:

* [Graphite & Carbon](https://graphiteapp.org/) is used to store the metrics.
* [Grafana](https://grafana.com/) is used for Dashboards and Alerts

All of the above other than the optional third-party applications Graphite & Grafana runs on the central Raspberry PI 4B.
The one I'm using has 4GB of Ram which is more than enough to run everything.
I suspect you might be able to get away with a 1GB or 2GB safely.

## Notes

1. I may move away from Graphite as Carbon cannot handle sensors that do not publish often, like the weather sensors.
2. My own PicoW based sensors will have a common firmware. That will be available in its own repository.
3. The Pimoroni Enviro-Weather PicoW board's [firmware](https://github.com/pimoroni/enviro) currently does not support the Rain Gauge. I've made modifications
   to its firmware. I'm intending to clean that code up or rewrite it to my own PicoW firmware.

