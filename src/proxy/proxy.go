package proxy

import (
  "encoding/base64"
  "fmt"
  "io"
  "log"
  "net"
  "net/http"
)

var (
  User = ""
  Pass = ""
  Host = ""
  port = ""

  auth = "Basic " + base64.StdEncoding.EncodeToString([]byte(User+":"+Pass))
  Proxy = fmt.Sprintf("%s:%s", Host, port)
)

func HandleHttp(w http.ResponseWriter, r *http.Request) {
  log.Printf("[%s] %s", r.Method, r.URL)

  // get connection
  hj, ok := w.(http.Hijacker)
  if !ok {
    panic("server doesn't support hijacking")
  }

  conn, _, err := hj.Hijack()
  if err != nil {
    log.Fatal(err)
  }

  // connect proxy
  proxyConn, err := net.Dial("tcp", Proxy)
  if err != nil {
    proxyConn.Close()
    log.Fatal(err)
  }

  // auth proxy
  r.Header.Set("Proxy-Authorization", auth)
  r.Write(proxyConn)

  // localhost to proxy
  go func() {
    io.Copy(conn, proxyConn)
    proxyConn.Close()
  }()

  // proxy to localhost
  go func() {
    io.Copy(proxyConn, conn)
    conn.Close()
  }()
}