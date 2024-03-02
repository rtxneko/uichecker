package main

import (
	"encoding/json"
	"github.com/phuslu/iploc"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	up             = "0"
	down           = "0"
	total          = "0"
	remark         = ""
	enable         = "true"
	expiryTime     = "0"
	listen         = ""
	protocol       = "dokodemo-door"
	streamSettings = `{
  "network": "tcp",
  "security": "none",
  "tcpSettings": {
    "header": {
      "type": "none"
    }
  }
}`
	sniffing = "{}"
)

func addProxy(ip string, session *http.Cookie, of *os.File, remoteAddress string, remotePort string) {
	settings := `{
  "address": "` + remoteAddress + `",
  "port": ` + remotePort + `,
  "network": "tcp,udp"
}`
	port := remotePort
	endpoint := "http://" + ip + ":54321/xui/inbound/add"
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	payload := "up=" + up + "&down=" + down + "&total=" + total + "&remark=" + remark + "&enable=" + enable + "&expiryTime=" + expiryTime + "&listen=" + listen + "&port=" + port + "&protocol=" + protocol + "&settings=" + url.QueryEscape(settings) + "&streamSettings=" + url.QueryEscape(streamSettings) + "&sniffing=" + url.QueryEscape(sniffing)
	var data = strings.NewReader(payload)
	req, err := http.NewRequest("POST", endpoint, data)
	if err != nil {
		log.Println("Could not compose request: ", err)
	}
	req.AddCookie(session)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	input := new(xResponse)
	err = json.NewDecoder(resp.Body).Decode(input)
	if input.Success {
		log.Println("[Success]", ip, iploc.Country(net.ParseIP(ip)))
		_, err := of.WriteString(ip + "\n")
		if err != nil {
			log.Panicln(err)
		}
		err = of.Sync()
		if err != nil {
			log.Panicln(err)
		}
	}
}
