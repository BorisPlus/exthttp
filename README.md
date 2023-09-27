# External HTTP Server

Just for special purpose HTTP server:

At HTTP-headers:

* "/image.jpg" - return green image;
* "/" - return "I receive teapot-status code!" with "Teapot" HTTP status (418, RFC 9110).

Directories at files:

* /log
  
## Debug

```bash
go clean -testcache && go test -v ./
```
