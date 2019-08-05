# API Server documentation

This documentation is a quick overview of the CRUD operations that take place within `plunder`, this should be a living document as the various endpoint mature over time. 

## Using the API Server

The API Server now starts as default and listens on a different port to HTTP services used for deployment, by default the `plunder` API server will listen on port `60443` however the `-p` `--port` flag can be used to specify a specific port. Currently the API server will bind to all interfaces.

### Starting the API Server

The below example will start the API server on a custom port.

```
plunder server -p 12345
```

## Accessing the API Server

The API Endpoints should be accessed using REST methodologies and JSON payloads, the API Endpoints should **always** be defined in `endpoints.go` (this may change later). 

## Current issues

### Server configuration

Currently DHCP can be stopped and started but logging output is buggy, HTTP/TFTP Can be started but can't be stopped or restarted.
