package common

import (
	"crypto/tls"
	"encoding/json"
	"github.com/4dogs-cn/TXPortMap/pkg/output"
	"github.com/4dogs-cn/TXPortMap/pkg/ping"
	"golang.org/x/net/context"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

func defaultCheckRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 1 {
		return http.ErrUseLastResponse
	}
	return nil
}

var dialer = &net.Dialer{
	Timeout:   3 * time.Second,
	KeepAlive: 1 * time.Second,
	// DualStack: true, // this is deprecated as of go 1.16
	// or create your own transport, there's an example on godoc.
}

func checkAvailability(domain string, maddr string) *output.ResultSuccess {
	even := &output.ResultSuccess{
		Target:  maddr,
		StepIP:  "",
		Country: "",
		Status:  false,
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if strings.Index(addr, domain) != -1 {
				addr = maddr
			}
			return dialer.DialContext(ctx, network, addr)
		},
	}

	client := &http.Client{Transport: tr}
	client.CheckRedirect = defaultCheckRedirect

	maddr_arr := strings.Split(maddr, ":")
	ip := maddr_arr[0]
	port := maddr_arr[1]
	resp, err := client.Get("https://" + domain + ":" + port + "/ip")
	if err == nil {
		body, _ := ioutil.ReadAll(resp.Body)
		if len(body) > 0 {
			var text = string(body)
			if strings.Index(text, "request success your ip is") != -1 {
				even.StepIP = strings.Split(text, "your ip is ")[1]
				even.Status = true
				resp, err = http.Get("http://geoip.apie.cc/index.php?security=CUe36wCk28cVw2&ip=" + ip)
				if err == nil {
					body, _ := ioutil.ReadAll(resp.Body)
					if len(body) > 0 {
						var country CountryInfo
						json.Unmarshal(body, &country)
						even.Country = country.Iso_code
					}
				}
				even.Ping = ping.Ping(ip, 2)
			}
		}
	}
	even.Time = time.Now()
	return even
}
