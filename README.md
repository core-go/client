# client
- This library is used to communicate with other services by http

## Installation
Please make sure to initialize a Go module before installing core-go/client:

```shell
go get -u github.com/core-go/client
```

Import:
```go
import "github.com/core-go/client"
```
## Features
### Initialize http client from config
- Initialize client for http
- Initialize client for https
### Log request, response at client
Support to turn on, turn off
- request
- response
- duration
- http response status code
- response content length

### Benefits
- Do not need to re-compile the service, user can switch client from http to https
- Do not need to re-compile the service, user can turn on, turn off the log (request, response, duration, response content length...)
### Microservice Architect
![Microservice Architect](https://cdn-images-1.medium.com/max/800/1*vKeePO_UC73i7tfymSmYNA.png)

### A typical micro service
When you zoom one micro service, the flow is as below
![A typical micro service](https://cdn-images-1.medium.com/max/800/1*d9kyekAbQYBxH-C6w38XZQ.png)
