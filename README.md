# counter

* **HTTP counter:** ![counter](http://162.243.131.99/http.png)
* **HTTPS counter:** ![counter](https://162.243.131.99/https.png)

## SSL cert generation

Generate self-signed SSL certs for the HTTPS listener by running:

```
go run $GOROOT/src/pkg/crypto/tls/generate_cert.go --host="localhost"
```