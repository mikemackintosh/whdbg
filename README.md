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


## Nginx
The following nginx configuration will allow you to use WSS correctly.

```nginx
server {
  listen 443 ssl;
  ssl_certificate /etc/ssl/certs/whdbg.crt;
  ssl_certificate_key /etc/ssl/private/whdbg.key;

  root /var/www/html;
  index index.html;;
  server_name _;

  location / {
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_pass http://localhost:8080;
  }
}
```
