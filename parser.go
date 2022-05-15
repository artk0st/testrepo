package parser

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Proxy struct {
	ip        string
	port      string
	geo       string
	speed     string
	protocol  string
	anonymity string
	update    string
}

var countParameter int = 0
var pages []Proxy
var hidemy Proxy

func UpdateProxies() {
	lines := getHidemy("https://hidemy.name/ru/proxy-list/?start=1")
	lines = append(lines, getHidemy("https://hidemy.name/ru/proxy-list/?start=65")...)
	lines = append(lines, getHidemy("https://hidemy.name/ru/proxy-list/?start=129")...)
	lines = append(lines, getHidemy("https://hidemy.name/ru/proxy-list/?start=193")...)
	lines = append(lines, getHidemy("https://hidemy.name/ru/proxy-list/?start=257")...)
	lines = append(lines, getHidemy("https://hidemy.name/ru/proxy-list/?start=321")...)
	lines = append(lines, getHidemy("https://hidemy.name/ru/proxy-list/?start=385")...)
	lines = append(lines, getHidemy("https://hidemy.name/ru/proxy-list/?start=449")...)
	lines = append(lines, getHidemy("https://hidemy.name/ru/proxy-list/?start=513")...)
	lines = append(lines, getHidemy("https://hidemy.name/ru/proxy-list/?start=567")...)

	f, err := os.OpenFile("proxies.list", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Truncate(0)

	for _, line := range lines {
		if rawConnect(line.ip, line.port) {
			if line.protocol == "HTTPS" {
				if _, err = f.WriteString("https://" + line.ip + ":" + line.port + "\n"); err != nil {
					panic(err)
				}
			} else {
				if _, err = f.WriteString("http://" + line.ip + ":" + line.port + "\n"); err != nil {
					panic(err)
				}
			}
		}
	}
	// if _, err = f.WriteString("http://52.183.8.192:3128\n"); err != nil {
	// 	panic(err)
	// }
}

func getHidemy(CollyDocsLink string) []Proxy {
	collector := colly.NewCollector()

	pages := make([]Proxy, 64)

	collector.OnHTML("td", func(e *colly.HTMLElement) {
		switch countParameter {
		case 0:
			hidemy.ip = e.Text
		case 1:
			hidemy.port = e.Text
		case 2:
			hidemy.geo = e.Text
		case 3:
			hidemy.speed = e.Text
		case 4:
			hidemy.protocol = e.Text
		case 5:
			hidemy.anonymity = e.Text
		case 6:
			hidemy.update = e.Text
		}
		countParameter++
		if countParameter == 7 {
			if hidemy.ip != "" && !strings.Contains(hidemy.ip, "IP") {
				if hidemy.protocol == "HTTP" || hidemy.protocol == "HTTPS" || hidemy.protocol == "HTTP, HTTPS" {
					pages = append(pages, hidemy)
				}
			}
			countParameter = 0
		}
	})
	collector.Visit(CollyDocsLink)
	return pages
}

func rawConnect(host string, port string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		fmt.Println("Connecting error:", err)
		return false
	}
	if conn != nil {
		defer conn.Close()
		fmt.Println("Opened", net.JoinHostPort(host, port))
	}
	return true
}
