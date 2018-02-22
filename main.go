package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

const (
	version     = "0.1.1"
	ifconfigUrl = "http://ifconfig.co/ip"
)

var (
	conMeta meta
	counter int
)

type meta struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip",omitempty`
	ExtIP    string `json:"extip",omitempty`
}

func (m *meta) newMeta() {
	m.Hostname, _ = os.Hostname()
}

func (m *meta) getIp() {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("Error getting interfaces: %s", err)
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Printf("Error getting ip from %v: %s", i, err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// fmt.Printf("%s\n", ip)
			if !ip.IsLoopback() {
				m.IP = fmt.Sprintf("%s", ip)
				break
			}
		}
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	counter += 1
	fmt.Fprintf(w, "%-15s %s\n", "Hostname:", conMeta.Hostname)
	fmt.Fprintf(w, "%-15s %s\n", "IP:", conMeta.IP)
	fmt.Fprintf(w, "%-15s %s\n", "External IP:", conMeta.ExtIP)
	fmt.Fprintf(w, "%-15s %d\n", "Request count:", counter)
	fmt.Fprintf(w, "%-15s %s\n", "Version:", version)
}

func (m *meta) getExtIp() (string, error) {
	resp, err := http.Get(ifconfigUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ip := net.ParseIP(strings.TrimSpace(string(body)))
	ipstr := fmt.Sprintf("%s", ip)
	if ip != nil {
		fmt.Println(ipstr)
		m.ExtIP = ipstr
		return ipstr, nil
	} else {
		return "", nil
	}
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	counter += 1
	cj, _ := json.Marshal(conMeta)
	fmt.Fprintf(w, "%s", string(cj))
}

func main() {
	conMeta.newMeta()
	conMeta.getIp()
	conMeta.getExtIp()
	counter = 0
	http.HandleFunc("/", handler)
	http.HandleFunc("/json", jsonHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
