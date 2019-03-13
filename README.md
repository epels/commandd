# uptimed

[![Build Status](https://travis-ci.com/epels/uptimed.svg?token=fFCsEza59TasaQuy6qDV&branch=master)](https://travis-ci.com/epels/uptimed)

Daemon exposing the output of the `uptime` command over HTTP.

## Configuration
The uptimed daemon allows basic configuration through flags. Comes with reasonable defaults.
 
```
$ ./uptimed -help
Usage of ./uptimed:
  -addr string
    	Address to listen on (default ":8080")
  -pattern string
    	Pattern to respond to. Set to / for any path (default "/uptime")
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
