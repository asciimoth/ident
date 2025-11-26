# ident
[RFC 1413](https://datatracker.ietf.org/doc/html/rfc1413) implementation

```go
package main

import (
	"fmt"
	"github.com/asciimoth/ident"
	"net"
)

func main() {
	listener, err := net.Listen("tcp4", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		_, lport, err := net.SplitHostPort(conn.LocalAddr().String())
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		rhost, rport, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		rhost = net.JoinHostPort(rhost, ident.DEFAULT_PORT)
		resp, err := ident.Query(conn.RemoteAddr().Network(), rhost, rport, lport)
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		fmt.Println("Connn from", conn.RemoteAddr(), resp)
		conn.Close()
	}
}
```

## License
Files in this repository are distributed under the CC0 license.  

<p xmlns:dct="http://purl.org/dc/terms/">
  <a rel="license"
     href="http://creativecommons.org/publicdomain/zero/1.0/">
    <img src="http://i.creativecommons.org/p/zero/1.0/88x31.png" style="border-style: none;" alt="CC0" />
  </a>
  <br />
  To the extent possible under law,
  <a rel="dct:publisher"
     href="https://github.com/asciimoth">
    <span property="dct:title">ASCIIMoth</span></a>
  has waived all copyright and related or neighboring rights to
  <span property="dct:title">ident</span>.
</p>

