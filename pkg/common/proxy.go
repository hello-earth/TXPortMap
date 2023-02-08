package common

import (
	"bytes"
	"github.com/4dogs-cn/TXPortMap/pkg/output"
	"golang.org/x/net/context"
	"io/ioutil"
	_ "io/ioutil"
	"log"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func getGid() (gid uint64) {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		panic(err)
	}
	return n
}

var domainDic = map[uint64]map[string]string{}

var mutex sync.Mutex

func resolve(domain string, ip string) {
	var gid = getGid()
	domain = strings.Split(domain, ":")[0] + strings.Split(ip, ":")[1]
	mutex.Lock()
	domainDic[gid] = map[string]string{domain: ip}
	mutex.Unlock()
}

func getKeys(m map[string]string) []string {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率很高
	j := 0
	keys := make([]string, len(m))
	for k := range m {
		keys[j] = k
		j++
	}
	return keys
}

func contains(list []string, pep string) bool {
	for _, v := range list {
		if v == pep {
			return true
		}
	}
	return false
}

func check(ip string) *output.ResultSuccess {
	dialer := &net.Dialer{
		Timeout:   3 * time.Second,
		KeepAlive: 30 * time.Second,
		// DualStack: true, // this is deprecated as of go 1.16
	}
	// or create your own transport, there's an example on godoc.
	http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		var guid = getGid()
		var domains = getKeys(domainDic[guid])
		var rDomain = strings.Split(addr, ":")[0]
		if contains(domains, rDomain) {
			addr = domainDic[getGid()][rDomain]
			log.Println(addr)
		}
		return dialer.DialContext(ctx, network, addr)
	}
	even := &output.ResultSuccess{
		Target: ip,
		Info:   "",
		Status: false,
	}

	resp, err := http.Get("https://worker.aproxy.tk:" + strings.Split(ip, ":")[1] + "/ip")
	if err == nil {
		body, _ := ioutil.ReadAll(resp.Body)
		if len(body) > 0 {
			var text = string(body)
			even.Info = text
			even.Status = true
		}
	} else {
		log.Println(err)
	}
	even.Time = time.Now()
	return even
}
