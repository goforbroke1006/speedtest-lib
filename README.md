# speedtest-lib

Speed test lib provides mechanism for network bandwidth check

### Supported provider

* Ookla's https://www.speedtest.net/
* Netflix's https://fast.com/

### Minimal server

Check file [examples/api/main.go](./examples/api/main.go)

### How to run main example

In terminal

```shell
go run examples/api/main.go
```

And open in browser next links:

* http://localhost:8080/healthz (for health check)
* http://localhost:8080/readyz (for readiness check - all providers return data)
* http://localhost:8080/test?provider=ookla (get speed test data from Ookla provider)
* http://localhost:8080/test?provider=netflix (get speed test data from Netflix provider)
