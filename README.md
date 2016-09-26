# whdbg
Webhook Debugger

### Summary

Webhook Debugger is used to parse formal HTTP requests being received from external API's, also know as Webhooks.

## Example

You can run the server without building, locally:

    $ go run main.go

When you make a request to `:8080`, the request details will be displayed in the console via `stdout`.

Response:

    GET / HTTP/1.1
    Host: localhost:8080
    Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
    Accept-Encoding: gzip, deflate, sdch
    Accept-Language: en-US,en;q=0.8
    Connection: keep-alive
    Upgrade-Insecure-Requests: 1
    User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.116 Safari/537.36


    GET /this/is?some=path HTTP/1.1
    Host: localhost:8080
    Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
    Accept-Encoding: gzip, deflate, sdch
    Accept-Language: en-US,en;q=0.8
    Connection: keep-alive
    Upgrade-Insecure-Requests: 1
    User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.116 Safari/537.36

