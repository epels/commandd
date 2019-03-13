# uptimed

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
