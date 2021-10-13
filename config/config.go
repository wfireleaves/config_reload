package config

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"sync"
	"sync/atomic"
)

type Cfg struct {
	ServiceCfg *Service
	health     bool
}

var (
	configList []*Cfg
	curCfgIdx  int32
	reloadLock sync.Mutex
)

func init() {
	configList = make([]*Cfg, 2)
	configList[0] = &Cfg{}
	configList[1] = &Cfg{}
	curCfgIdx = 0
	if err := ReLoadCfg(); err != nil {
		panic(fmt.Sprintf("init load config failed, %s", err.Error()))
	}
	fmt.Println(GetCurCfg().ServiceCfg.String())
}

func getCurIdx() int32 {
	return atomic.LoadInt32(&curCfgIdx)
}

func getNextCfg() *Cfg {
	return configList[(getCurIdx()+1)%2]
}

func GetCurCfg() *Cfg {
	cfg := configList[getCurIdx()]
	if cfg.health {
		return cfg
	}
	fmt.Println("cfg not ready", getCurIdx())
	if err := ReLoadCfg(); err != nil {
		fmt.Println("get config reload failed", err.Error())
		panic("reload failed")
	}
	cfg = configList[getCurIdx()]
	if !cfg.health {
		fmt.Println("reload config not working")
		panic("reload not working")
	}
	return cfg
}

func ReLoadCfg() error {
	reloadLock.Lock()
	defer reloadLock.Unlock()

	fmt.Println("start reload config, cur idx:", getCurIdx())

	nextCfg := getNextCfg()
	nextCfg.ServiceCfg = &Service{}
	xmlFile, err := ioutil.ReadFile("cfg_xml/service.xml")
	if err != nil {
		return fmt.Errorf("read xml file failed,err:%s", err.Error())
	}
	err = xml.Unmarshal(xmlFile, nextCfg.ServiceCfg)
	if err != nil {
		return fmt.Errorf("unmarshal xml file failed, err:%s", err.Error())
	}
	nextCfg.health = true

	atomic.StoreInt32(&curCfgIdx, (getCurIdx()+1)%2)
	configList[(getCurIdx()+1)%2].health = false

	fmt.Println("finish reload config, cur idx:", getCurIdx())

	return nil
}
