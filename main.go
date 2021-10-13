package main

import (
	"config_reload/config"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func reloadSingal() {
	err := config.ReLoadCfg()
	if err != nil {
		fmt.Println("reload failed", err.Error())
	}
}

func getConfig() {
	fmt.Println(config.GetCurCfg().ServiceCfg.String())
}

func main() {
	ch := make(chan os.Signal, 5)
	signal.Notify(ch)
	ticker := time.NewTicker(time.Second * 3)
	var wg sync.WaitGroup
	for {
		select {
		case sig := <-ch:
			switch sig {
			case syscall.SIGINT:
				fmt.Println("service will eixt")
				goto END
			case syscall.SIGHUP:
				fmt.Println("service will reload config")
				reloadSingal()
			default:
				fmt.Println("not catch")
			}
		case <-ticker.C:
			for i := 0; i < 5; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					getConfig()
				}()
			}
		}
	}

END:
	wg.Wait()
	fmt.Println("progress end")
}
