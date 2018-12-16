package main

import (
  "fmt"
  "log"
  "net/http"
  "net/url"
  
  "github.com/getlantern/systray"
  
  "./src/icon"
	"./src/proxy"
)

func onReady() {
	systray.SetIcon(icon.Data)
  systray.SetTitle("Auto Auth Proxy")
  systray.SetTooltip("localhost:8480")
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
  log.Fatal(http.ListenAndServe(":8480", handler))
}

func onExit() {}

func main() {
	systray.Run(onReady, onExit)
}