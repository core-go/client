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
![Microservice Architect](https://camo.githubusercontent.com/cf46a1780520d3612f1d81b219b56a14428fc24bb4ae9f4eede169aa9c58bee8/68747470733a2f2f63646e2d696d616765732d312e6d656469756d2e636f6d2f6d61782f3830302f312a764b6565504f5f5543373369377466796d536d594e412e706e67)

### A typical micro service
When you zoom one micro service, the flow is as below
![A typical micro service](https://camo.githubusercontent.com/581033268b9152e7ea8881904f533a51a29eeb3a63e8d6478540668c6e422ce3/68747470733a2f2f63646e2d696d616765732d312e6d656469756d2e636f6d2f6d61782f3830302f312a64396b79656b416251594278482d4336773338585a512e706e67)
