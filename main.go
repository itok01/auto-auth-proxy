package main

import (
  "fmt"
  "log"
  "net/http"
  "net/url"
  "os"

  "github.com/getlantern/systray"

  "github.com/itok01/auto-auth-proxy/src/icon"
  "github.com/itok01/auto-auth-proxy/src/proxy"
)

var localProxy = os.Args[5]

func onReady() {
  systray.SetIcon(icon.Data)
  systray.SetTitle("Auto Auth Proxy")
  systray.SetTooltip(fmt.Sprintf("localhost:%s", localProxy))
  mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

  go func() {
    <-mQuit.ClickedCh
    systray.Quit()
  }()

  proxyUrl, err := url.Parse(fmt.Sprintf("http://%s:%s@%s", proxy.User, proxy.Pass, proxy.Proxy))
  if err != nil {
    log.Fatal(err)
  }
  http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

  handler := http.HandlerFunc(proxy.HandleHttp)
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", localProxy), handler))
}

func onExit() {}

func main() {
  systray.Run(onReady, onExit)
}
