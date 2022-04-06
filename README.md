# go-crypto-dashboard

This repository contains various tools for tracking and managing your crypto assets.

## Available tools

* Scraper - scrapes your portfolios for further analysis
* Metabase - lets you analyze your PnL, balance and other metrics

## Tools in development

* [Passivbot](https://www.passivbot.com/) tools
  * Droplet - a tool with an optional API to manage passivbot instances
  * Bucket - web UI for managing droplets
  * Hosting - API that allows you to launch and configure servers on various cloud providers

## How to use

### docker-compose

1. Clone the repository
1. Run `docker-compose up -d`
1. Open your browser and go to https://localhost:3000

* Run `docker-compose -d --build` to rebuild the containers, if you change the configuration or the code

### Scraper

You'll need Go 1.18 or later to build the project.

1. Clone the repository
1. Run `go get` to download the dependeencies
1. Run `go build -o bin/[name] cmd/[name]/main.go` to build the binary of a particular tool or the entire project.
1. Run `./bin/[name] to start the binary.

## Metabase credinentials

* Login: exchanges@dashboard.com
* Password: admin01

