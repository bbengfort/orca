# Orca

[![Build Status](https://travis-ci.org/bbengfort/orca.svg?branch=master)](https://travis-ci.org/bbengfort/orca)
[![Coverage Status](https://coveralls.io/repos/github/bbengfort/orca/badge.svg?branch=master)](https://coveralls.io/github/bbengfort/orca?branch=master)
[![GoDoc Reference](https://godoc.org/github.com/bbengfort/orca?status.svg)](https://godoc.org/github.com/bbengfort/orca)
[![Go Report Card](https://goreportcard.com/badge/github.com/bbengfort/orca)](https://goreportcard.com/report/github.com/bbengfort/orca)
[![Stories in Ready](https://badge.waffle.io/bbengfort/orca.png?label=ready&title=Ready)](https://waffle.io/bbengfort/orca)

**Echolocation of device with static nodes and network latency.*

[![Orca][orca.jpg]][orca_flickr]

Orca is a reflector/ping generator service that is intended to measure network latency of gRPC requests from mobile devices (laptops) to storage locations in home and work networks. Orca was built to get some baseline metrics for latencies in distributed storage systems and is purely experimental, not intended for large scale use.

Orca specifies two primary commands: _reflect_ and _generate_. The reflector is a Protocol Buffers service that accepts echo requests and returns replies. The generator is a service that sends a request to all known reflectors on a routine interval and logs the latency in a SQLite database.

## Installation

In order to install Orca and get it running on your system, you have two options: build from source or download a prebuilt binary that may or may not exist for your system.

To download the prebuilt binary, visit the [Current Release](#) on GitHub and download the executable for your system. Extract the archive to reveal the executable, and put it on your path, making sure it has executable permissions.

Note: currently only Linux x64 and OS X (darwin) x64 binaries are prebuilt. My preference for path location is `$HOME/bin` but make sure that location is in your path. Otherwise, `/usr/local/bin` is a good choice for moving the binary to.

To build from source, make sure you have your Go environment setup and run:

```
$ go get github.com/bbengfort/orca/...
```

You can then use `go install` to install the binary to `$GOBIN` as follows:

```
$ go install cmd/orca.go
```

To contribute or modify the orca command, you'll want to go the source route; otherwise if you are using Linux or OS X, I recommend going the binary route. If there is an architecture that you'd like prebuilt, please let me know in a GitHub issue.

## Getting Started

Once installed, some configuration is necessary. This is the recommended configuration, though other configuration options are available and can be inspected via the `orca help` command.

1. Create the directory where configuration and data will reside:

        $ mkdir ~/.orca

2. Copy the example configuration file to the orca directory

        $ cp fixtures/orca-config-example.yml ~/.orca/config.yml

    Note that the first path refers to the fixtures directory inside of the
    GitHub repository. If you downloaded the binary package from the releases page, then the example configuration will be in the archive you downloaded.

3. Edit the configuration as needed, noting the defaults and comments in the example configuration file.

4. Create the orca SQLite database used to track meta information about reflection and generators.

        $ orca createdb

At this point orca is configured to begin reflecting, however generators require an extra configuration step.

### Relectors

Running orca as a reflector runs a gRPC service that listens for echo requests on the current IP address and port, increments the request sequence (to detect out of order messages), and responds with an echo reply. This process means that there is a small database transaction before reply, so latency will not be exactly the same as using the `ping` utility. If `debug: true` in the configuration, then the server will also log all incoming echo requests to `stdout`, which can be redirected to a log file if required.

Run the reflector as follows:

```
$ orca reflect
```

To run in the background on a constant server:

```
$ nohup orca reflect &
```

Upstart and LaunchAgent scripts for managing the background process are forthcoming, though I'd be happy to accept a pull request for them!

### Generators

Generators require a bit more configuration, since you'll have to add all of the reflectors that you want the generator to ping to the database. Do this via the `orca devices --add` command:

```
$ orca devices --add
Enter device name []: rogue
Enter device IP address []: 1.2.3.4:3265
Enter device domain []:
```

Do this for as many devices as you'd like to ping on each interval. To see the devices already added use `orca devices --list`. Run the generator as follows:

```
$ orca generate
```

Similar to the reflector, you'll have to nohup and background this in order to ensure it always runs. LaunchAgent and Upstart scripts are coming soon. The generator waits until the interval has passed, loads up the list of devices to ping, and sends an echo request to them, recording the request (in the case of non-connectivity) and sequence number in the database. On receipt of the reply, it measures latency and stores the information in the database.

## Running Agents

Orca is a long running process that conducts work on a routine interval. As a result, you'll want to run orca as an *agent* - a background daemon that runs on behalf of a user. Running an agent depends on the operating system environment, in this section we will present agent scripts for `launchd` on OS X and Upstart on Ubuntu.

### Launch Agent

Orca is designed to be [launchd](https://developer.apple.com/library/content/documentation/MacOSX/Conceptual/BPSystemStartup/Chapters/CreatingLaunchdJobs.html) compliant to run as a user agent in the background and be started by the operating system. In order to configure Orca to be run in the background, a property list describing the agent needs to be installed to the Library and the binary must have a file mode that is not group or world writable.

```
$ chmod 600 /usr/local/bin/orca
$ cp fixtures/com.bengfort.orca.plist ~/Library/LaunchAgents
$ chmod 600 ~/Library/LaunchAgents/com.bengfort.orca.plist
$ launchctl load ~/Library/LaunchAgents/com.bengfort.orca.plist
```

Note that this assumes that you've placed the executable onto your path at `/usr/local/bin` (modify the path as needed) and that you've cloned the Orca repository. The plist file that describes the LaunchAgent sits in the fixtures directory in the root of the Orca repository.

### Upstart

The Orca reflector is designed to be run on an Ubuntu Linux server and therefore has an Upstart script that will make sure it is loaded and always running. The Upstart configuration is in the fixtures directory of the repository. Install as follows

```
$ sudo cp fixtures/orca.conf /etc/init/orca.conf
$ sudo chown root:root /etc/init/orca.conf
$ sudo service orca start
```

This should run the orca reflector service and keep it alive even when the server reboots.

## Location Services

Orca can provide location services for mobile devices via the [MaxMind GeoIP2 Precision City Service](https://www.maxmind.com/en/geoip2-precision-city-service). In order to enable location services, you need to register for a MaxMind developer account and include your API user id and license key in the YAML configuration file. Because MaxMind is a paid service, location lookups are only made when the current IP address of the machine changes.

Note also that the granularity for this service is limited; for example, a GeoIP2 lookup from my office in the A.V. Williams Building of the University of Maryland yielded the following location via latitude and longitude:

![Map Granularity](fixtures/map.png)

This location is not centered on my office, the building, or even in the center of the university. The level of granularity probably differs based on the type of network you're connected to. If a higher level of granularity is required, then the use of GPS is recommended. Additional location inaccuracy can come about when tethering to a mobile phone. Use specified locations with care!

## Acknowledgements

Orca is an open source project built to obtain metrics about mobile distributed systems and various latencies. If you'd like to contribute, I'd love some help, but no current plans are underway for future development.

### Attribution

The photo used in this README, &ldquo;[Orca (Design Exploration)][orca_flickr]&rdquo; by [Alberto Cerriteño](https://www.flickr.com/photos/acerriteno/) is used under a [CC BY-NC-ND 2.0](https://creativecommons.org/licenses/by-nc-nd/2.0/) creative commons license.

[orca.jpg]: fixtures/orca.jpg
[orca_flickr]: https://flic.kr/p/4HDnoE
