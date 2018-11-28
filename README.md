# goxy - A TCP tunnel over HTTP


## Bonus Features

* Allow multiple TCP connections to the same remote host.
* Local TCP and HTTP ports are customizable.
* Remote TCP and HTTP ports are customizable too.


## Obfuscation methods

### Random path (URI)

In order to allow multiple connections from the same client to the same remote Goxy server, an ID has to be passed to each HTTP request to allow the server to redirect the content to the correct server-side TCP connection.

When establishing a new connection between a Goxy client and a Goxy server, a handshake is made. The client asks for a id to the server, and the server return an unused id. This id is an array of at least three (classic) alphabetic letters.

Then, for each following request initiated, the client picks for each letter in the id, a random word highly used in the english language having its second to last letter equal to this letter.

These words are then joined with a '/' and used as an URI.

Example:

If the id given by the server is 'ALY', the client may make a request using the following URL: 

```text
'http://remoteserver.com/local/world/plays'
                            ^     ^     ^
                            A     L     Y
```


### Hide in common files 

To bypass proxies that try to detect the content type of files downloaded (or uploaded) to a server, each request append a random file extension to the URI (ie. '.png').

Then, depending on the file extension presents in the URI, the server adds the magic bytes corresponding to the file type at the start of the response body.

The 'Content-Type' HTTP header of the response is set to the corresponding mime-type.

Example:

If the request URI is '/local/world/plays.png', the response body will looks like this: 

```text
x89 x50 x4E x47 x0D x0A x1A x0A .. .. .. .. .. .. ..
\_______PNG Format magic______/ \__Useful_data_____/

```

When receiving a response, the client detect the file extension and ignore the N first bytes of the response. 


### Prioritize GET over POST

While browsing the web, the average ratio of GET requests vs the number of POST requests is more/less of 90% GET and 10% POST.

Assuming that an SSH connection write the same amount of data it reads from the TCP connection, making 50% POST and 50% GET may be a bad idea.

While inspecting the content length of the TCP packets while using SSH, packets size rarely exceed 128 bytes.

Encoding the content of client output data in base64 and putting this encoded data in a commonly used HTTP header allows the ratio of GET/POST request to be 95% GET for 5% POST requests. 


### Encode content

To prevent classic (but useful) string matching (like the SSH handshake header) all the content is encoded as base64.


### Long polling

To prevent unnecessary request to know if the server as available data, the server uses long polling. Doing so, if the user do not touch his terminal, no request should be completed.


### Using OTP

In order to prevent the server from asking data to the server, while the client did not ask for it, the server gives an OTP (One Time Password) for the next read and another one for the next write.

If the server receives a request with an invalid OTP, the server response with a useless 200 response so the proxy dosn't notice a mismatch in the status code.


### Use common User-Agents

To prevent proxy that blocks request based on the user agent of the request, a Google Chrome or Firefox user agent is set for each request.


## Building and starting the server and the client

### Installing

Golang 1.11.2 was used to make this project.

Be sure to have set your GOPATH variable.

```sh
mkdir -p $GOPATH/src/github.com/scotow
cp . $GOPATH/src/github.com/scotow/
```


### Running the server

```sh
go run $GOPATH/src/github.com/scotow/goxy/cmd/goxys/main.go
```

Use the `-h` options to change the default parameters


### Running the client

```sh
go run $GOPATH/src/github.com/scotow/goxy/cmd/goxyc/main.go

```

Use the `-h` options to change the default parameters


### SSH examples

Simple SSH connections

```sh
ssh -p 2222 localhost
```

SSH Reverse tunnel option can be used to connect from home to the company's computer:

```sh
ssh -p 2222 -R 2222:localhost:22 my.home.com
```

Then while at home:

```sh
ssh -p 2222 localhost
```