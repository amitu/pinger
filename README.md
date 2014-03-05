# Pinger

## What is this?

Lets say you are building Skype. And some customer complains that call quality was bad. So you should store some kind of ping rtt time for that client.

This is where this service comes in. This service listens on HTTP and UDP ports. So you send the IP address of customer, along with some other meta data, and a destination. And this service pings that IP, sends 10 packets, computes the average, min, max etc, and posts the data on the destination.

As destination it can use either HTTP or UDP, or both. 

## Using pinger over HTTP

Pinger by default listens on http://localhost:9918. This can be configured by command line parameter `--http=localhost:8816` for example. 

Pinger exposese `/pingit` where GET or POST requests can be made. Following parameters are expected:

http
: optional, http endpoint to post the response to. http has to full url. can include GET parameters. 

udp
: optional, udp server on which response should be posted

ip
: required, the IP to ping

data
: optional, some extra data you may want to send, to identify the client etc. 

/pingit returns "ok" and does not wait for pings to finish. 

**Note**: If both http and udp are provided data would be sent to both destination, if none is provided, it just logs information on stdout. 

## Using pinger over UDP

Pinger by default listens on localhost:9918 UDP port. This can be configured by command line parameter `--udp=localhost:8816` for example.

On UDP pinger expects packets, which is json containing the same parameters above, http, udp, ip and data. 

## Data to HTTP Destination

pinger pings the passed IP ten times, and computes the following details:

Success
: Out of 10 pings, how many got replies. 

Timeouts
: How many timed out. The reason both is provided is because pinger only allows 10 seconds for one IP test, and if ping latencies are high, it will not send 10 ping requests, only as many as that can fit in 10 seconds. 

Avg
: Average rtt over all successful replies. All rtts are in nanoseconds.

Min
: Minimum rtt among successful ones.

Max
: Maximum rtt among successful ones. 

Data
: A copy of data sent in original request. 

All these parameters are POSTed to the http parameter passed. 

## Data to UDP Destination.

Pinger sends the same information as a single UDP packet to destination UDP server, it encodes them as keys in a JSON encoded string. 

## Installation

If you have go setup with GOROOT set, do "go install github.com/amitu/pinger" and you are all set, `pinger` named binary would be created in GOROOT/bin. `pinger --help`. 