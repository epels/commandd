# commandd

[![Build Status](https://travis-ci.com/epels/commandd.svg?token=fFCsEza59TasaQuy6qDV&branch=master)](https://travis-ci.com/epels/commandd)

Daemon exposing the output of any arbitrary command over HTTP.

## Requirements

* `Go 1.13`

## Configuration

The commandd daemon allows basic configuration through flags. Comes with reasonable defaults.
 
```
$ ./commandd -help
Usage of ./commandd:
  -addr string
    	Address to listen on (default ":8080")
  -pattern string
    	Pattern to serve to (default "/run")
  -timeout duration
    	Timeout for command (default 10s)
```

Anything after the flags is the command to execute on requesting pattern. A typical invocation looks like this:

```bash
$ ./commandd -addr=":8080" -timeout="2s" echo -n foo bar baz
```

```bash
$ curl http://localhost:8080/run
foo bar baz
```

## Docker

Running in Docker is easy. First, build an image.

```bash
docker build -t some-tag .
```

To run a container from the freshly built image:

```bash
docker run -p 8080:8080 --rm -t some-tag
```

The `-p` flag publishes the container's port to the host.
