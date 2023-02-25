package main

import (
	"flag"
	"fmt"
	"github.com/4dogs-cn/TXPortMap/pkg/common"
	_ "github.com/projectdiscovery/fdmax/autofdmax" //Add automatic increase of file descriptors in linux
	"os"
	"time"
)

func init() {
	flag.Parse()

	// fmt.Println("threadnum: ", common.NumThreads)

	if common.NumThreads < 1 || common.NumThreads > 3000 {
		fmt.Println("number of goroutine must between 1 and 3000")
		os.Exit(-1)
	}
}

// 建议扫描top100或者top1000端口时使用顺序扫描，其它情况使用随机扫描
func main() {

	// trace追踪文件生产，调试时打开注释即可
	/*
		f1, err := os.Create("scan.trace")
		if err != nil {
			log.Fatal(err)
		}
		trace.Start(f1)
		defer trace.Stop()
	*/
	//common.ArgsPrint()
	engine := common.CreateEngine()

	// 命令行参数错误
	if err := engine.Parser(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// common.ArgsPrint()
	//engine.Wg.Add(engine.WorkerCount)
	//go engine.Scheduler()
	engine.Run()

	// 等待扫描任务完成
	engine.Wg.Wait()
	println("scan task finish, waiting for cdn test task complete")
	close(engine.ProxyChan)

	timer := time.NewTimer(time.Second * 300)

	select {
	case <-timer.C:
		panic("it take too long time for complete, i can not wait anymore!!!")
	}
	engine.PWg.Wait()
	timer.Stop()

	if common.Writer != nil {
		common.Writer.Close()
	}
}
