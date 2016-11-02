# Orca Ping

[![Build Status](https://travis-ci.org/bbengfort/orca.svg?branch=master)](https://travis-ci.org/bbengfort/orca)
[![Coverage Status](https://coveralls.io/repos/github/bbengfort/orca/badge.svg?branch=master)](https://coveralls.io/github/bbengfort/orca?branch=master)
[![GoDoc Reference](https://godoc.org/github.com/bbengfort/orca?status.svg)](https://godoc.org/github.com/bbengfort/orca)
[![Go Report Card](https://goreportcard.com/badge/github.com/bbengfort/orca)](https://goreportcard.com/report/github.com/bbengfort/orca)
[![Stories in Ready](https://badge.waffle.io/bbengfort/orca.png?label=ready&title=Ready)](https://waffle.io/bbengfort/orca)

**This is the Ping branch of the Orca project**

[![Orca][orca.jpg]][orca_flickr]

**Echolocation of device with static nodes and network latency.**

Orca is a ping/listener utility that is intended to measure network latency of gRPC requests from mobile devices (laptops) to storage locations in home and work networks. Orca was built to get some baseline metrics for latencies in distributed storage systems and is purely experimental, not intended for large scale use.

Orca specifies two primary commands: _listen_ and _ping_. The listener is a Protocol Buffer service that accepts echo requests and returns replies. The ping utility sends echo requests to the specified IP address and measures the aggregate latency.

## Usage

To run a listener on a server use the `orca listen` command:

```
$ orca listen --help
```

To ping a listener use the `orca ping` command:

```
$ orca ping --help
```

If you use it like the ping utility that implements the ICMP protocol, then you're doing it right!

## Acknowledgements

Orca is an open source project built to obtain metrics about mobile distributed systems and various latencies. If you'd like to contribute, I'd love some help, but no current plans are underway for future development.

### Attribution

The photo used in this README, &ldquo;[Orca Sighting][orca_flickr]&rdquo; by [Jay Cox](https://www.flickr.com/photos/jaycoxfilm/) is used under a [CC BY-NC-ND 2.0](https://creativecommons.org/licenses/by-nc-nd/2.0/) creative commons license.

[orca.jpg]: fixtures/orca.jpg
[orca_flickr]: https://flic.kr/p/4Nkop2
