# Mockzilla
A simple Go program that receives data over HTTP/S and consumes it by displaying the content

### Run mockzilla
```bash
./mockzilla -response-message='{"status":"Success"}' -content-type="application/json" -response-code=201 -cert=server.crt -key=server.key -port=8443
```
