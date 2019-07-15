# topz

List running processes like `top` but in browser. Stolen from [Brendan Burns](https://github.com/brendandburns/topz) :-).

## Run

```
$ go run cmd/server/main.go
$ curl localhost:8080/topz
```

# Container version

You can use topz as a sidecar container to see what's happening inside an application container.

```
# run the main application container
$ docker run -d <my-app-image> # e.g. nginx
<container-hash-value>

# stash that value in an envvar
$ export APP_ID=<container-hash-value>

# run topz sidecar in the same PID namespace
$ docker run --pid=container:${APP_ID} \
  -p 8080:8080 \
  brendanburns/topz:db0fa58 \
  /server --addr=0.0.0.0:8080
```
