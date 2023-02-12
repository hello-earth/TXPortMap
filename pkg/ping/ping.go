package ping

import (
	"github.com/go-ping/ping"
	"os"
	"os/signal"
	"time"
)

func Ping(dstIP string, count int) int64 {
	if dstIP == "0.0.0.0" {
		return 9999
	}
	pinger, err := ping.NewPinger(dstIP)
	if err != nil {
		panic(err)
	}
	pinger.Count = count
	pinger.SetPrivileged(true)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			pinger.Stop()
		}
	}()
	pinger.Timeout = 2 * time.Second
	err = pinger.Run()
	if err != nil {
		panic(err)
		return 9999
	}
	stats := pinger.Statistics()
	return stats.AvgRtt.Milliseconds()
}
