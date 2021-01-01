# HTTP Relay
HTTP client to client (web browser to web browser) communication

## Features
- Communication using HTTP requests (supported by the all web browsers and most platforms)
- Synchronous bidirectional data exchange
- Synchronous one directional data transfer
- Asynchronous and buffered one to many data transfer
- Client as a server communication method
- Communication channel ownership

## Installation
### Download
[Download executable](https://gitlab.com/jonas.jasas/httprelay/-/releases) (Linux, Mac, Windows)

### Docker
Docker image is built without additional layers. Final image size is less than 3Mb.

- Latest image: `jonasjasas/httprelay`
- [Image list](https://hub.docker.com/r/jonasjasas/httprelay)
- Run: `docker run -p 8080:8080 jonasjasas/httprelay`

### Build
Install the package to your [$GOPATH](https://github.com/golang/go/wiki/GOPATH "GOPATH") with the [go tool](https://golang.org/cmd/go/ "go command") from shell:

```bash
go get gitlab.com/jonas.jasas/httprelay
cd ~/go/src/gitlab.com/jonas.jasas/httprelay
go run ./cmd/...
```

Make sure [Git is installed](https://git-scm.com/downloads) on your machine and in your system's `PATH`.

### Test installation

Go to http://localhost:8080/health should display version number. 

## Communication methods

### Sync (Alice <-> Bob)
Sync communication method provides HTTP client to HTTP client (web browser to web browser) synchronous data exchange.
Two requests can exchange data on any HTTP method.

- Alice: `GET https://demo.httprelay.io/sync/your_secret_channel_id?msg=Hello-Bob`
- Bob: `GET https://demo.httprelay.io/sync/your_secret_channel_id?msg=Hello-Alice`

URL query data is placed in `httprelay-query` response header field.

- **[Text message exchange example](https://jsfiddle.net/jasajona/y35rLnd9/)** Exchange text messages using just GET requests and query parameters
- **[GPS tracker example](https://jsfiddle.net/jasajona/cgaju9o8/)** Exchange coordinates and track each others location on map in real time


If the method supports content transfer (e.g. POST, PUT etc.) data is going to be received as a response body by the counterpart.
```shell script
curl -v https://demo.httprelay.io/sync/your_secret_channel_id?msg=I-love-you-Bob    
...
> GET /sync/your_secret_channel_id?msg=I-love-you-Bob HTTP/2
> Host: demo.httprelay.io
> User-Agent: curl/7.58.0
> Accept: */*
...
< HTTP/2 200 
< access-control-expose-headers: X-Real-IP, X-Real-Port, Httprelay-Time, Httprelay-Your-Time, Httprelay-Method, Httprelay-Query
< content-type: text/plain
< httprelay-method: POST
< httprelay-time: 1572596282762
< httprelay-your-time: 1572596271247
< x-real-ip: 1.2.3.4
< x-real-port: 31153
< content-length: 16
< 
I love you Alice
```

```shell script
curl -X POST -v -H "Content-Type: text/plain" --data "I love you Alice" https://demo.httprelay.io/sync/your_secret_channel_id
...
> POST /sync/your_secret_channel_id HTTP/2
> Host: demo.httprelay.io
> User-Agent: curl/7.58.0
> Accept: */*
> Content-Type: text/plain
> Content-Length: 16
... 
< HTTP/2 200 
< access-control-expose-headers: X-Real-IP, X-Real-Port, Httprelay-Time, Httprelay-Your-Time, Httprelay-Method, Httprelay-Query
< httprelay-method: GET
< httprelay-query: msg=I-love-you-Bob
< httprelay-time: 1572596271247
< httprelay-your-time: 1572596282762
< x-real-ip: 1.2.3.4
< x-real-port: 31152
< content-type: text/html
< content-length: 0
< 
```

### Link (Alice -> Bob)
Link communication method provides HTTP client to HTTP client (web browser to web browser) synchronous one directional data transfers.
Link communication method implements producer -> consumer pattern.
Producer must use `POST` method, consumer must use `GET` method.  

- Producer: `POST https://demo.httprelay.io/link/your_secret_channel_id`
- Consumer: `GET https://demo.httprelay.io/link/your_secret_channel_id`

Producer's request will be finished when consumer makes the request.
If consumer makes request prior producer, receiver request will wait till producer makes the request.

- **[Text message transfer example](https://jsfiddle.net/jasajona/q6uhLuqf/)**
- **[GPS tracker example](https://jsfiddle.net/jasajona/mjrwLc3d/)**

### Mcast (Alice -> Bob, Carol)
Mcast communication method provides one to many buffered and asynchronous HTTP client to HTTP client (web browser to web browser) data transfers.
Mcast communication method must be used when there are multiple consumers and producer don't need to know when or if receivers received it's data.

- Producer (Alice): `POST https://httprelay.io/mcast/your_secret_channel_id`
- Consumer (Bob): `GET https://httprelay.io/mcast/your_secret_channel_id`
- Consumer (Carol): `GET https://httprelay.io/mcast/your_secret_channel_id`

Producers's request will finish as soon as all data is transferred to the server.
Currently data is buffered in memory for 20 minutes (next Httprelay versions are going to support more data storage options).
If consumer makes request prior producer, consumer request will wait till producer makes the request.
Each producer request receives `httprelay-query` header field with the currently sent data sequence number.
Each consumer request receives `httprelay-query` header field with the currently received data sequence number and cookie `SeqId` with the next sequence number.
Cookies must be enabled or `SeqId` query parameter must be provided in consumer's request.
If there is no `SeqId` provided, consumer will receive most recent data.
If consumer provides `SeqId` greater than most recent `SeqId`. Request will wait till new data received from producer.

- **[Message transfer example](https://jsfiddle.net/jasajona/ntwmheaf/)**
- **[Multi-user painting example](https://jsfiddle.net/jasajona/ky0cLgf9/)**
- **[Location sharing example](https://jsfiddle.net/jasajona/5ks1y3nL/)**
- **[Image exchange example](https://jsfiddle.net/jasajona/f2an7tjh/)**


### Proxy (Server <-> Client)
Proxy communication method allows an HTTP client to act as a server.
Using this method you can turn web browser or any other HTTP client into a server.

[HTTP Relay JavaScript library](https://gitlab.com/jonas.jasas/httprelay-js) 
is a small framework that abstracts communication with the HTTP Relay server and lets you feel like your web browser is a server accessible online.

- **[Basic usage](https://codesandbox.io/s/ik6w1)**
- **[Interactive page](https://codesandbox.io/s/gzwuv)**
- **[File sharing](https://codesandbox.io/s/mrhsf)**
- **[REST API](https://codesandbox.io/s/uu9l3)**
- **[Assets](https://codesandbox.io/s/90t1d)**

 
## Writing permission
Producer can take ownership of `channel id` by providing `wsecret` query parameter.
Channel ownership can be set in Link and Mcast communication methods.

`POST https://httprelay.io/mcast/mychan?wsecret=qwerty`

after setting `wsecret=qwerty` on `mychan` channel, subsequent POST requests must provide `wsecret=qwerty` query parameter to write new data in to channel.
If there are no successful POST requests, channel ownership expires in 20 minutes.



## Response headers
All communication methods receives response header fields with the following information:

- **httprelay-method**: counterpart request HTTP method (GET, POST, PUT, etc.).
- **httprelay-query**: counterpart request URL query (everything what follows after `?`).
- **httprelay-time**: counterpart request timestamp. Moment when counterpart request was received.
- **httprelay-your-time**: current request timestamp. Moment when current request was received.
- **x-real-ip**: counterpart IP address.
- **x-real-port**: counterpart outbound port.

## Command line arguments
- **-a** Bind to IP address (default 0.0.0.0)
- **-p** Bind to port (default 8080)
- **-u** Bind to Unix socket path
- **-h** Print help