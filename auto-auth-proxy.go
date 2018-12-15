package main

import (
  "encoding/base64"
  "fmt"
  "io"
  "log"
  "net"
  "net/http"
  "net/url"
)

var (
  user = ""
  pass = ""
  host = ""
  port = ""
  
  auth = "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+pass))
  proxy = fmt.Sprintf("%s:%s", host, port)
)

func handleHttp(w http.ResponseWriter, r *http.Request) {
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
  proxyConn, err := net.Dial("tcp", proxy)
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


func main() {
  proxyUrl, err := url.Parse(fmt.Sprintf("http://%s:%s@%s", user, pass, proxy))
  if err != nil {
    log.Fatal(err)
  }
  http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

  handler := http.HandlerFunc(handleHttp)
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), handler))
}