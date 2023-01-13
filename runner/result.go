package runner

import (
	util "github.com/hktalent/go-utils"
)

// 记录日志到 大数据搜索引擎
func SendEsLog(m1 interface{}) {
	if 0 == len(util.EsUrl) {
		return
	}
	szId := "xxx"
	util.SendReq(&m1, szId, "ksubdomain")

}

var bOk = make(chan struct{})
var bDo = make(chan struct{})
var oR = make(chan interface{}, 5000)

func DoSaves() {
	var n = len(oR)
	var oS = make([]interface{}, n)
	for n > 0 {
		oS = append(oS, <-oR)
		n--
	}
	util.DoSyncFunc(func() {
		SendEsLog(&oS)
	})
}

func DoRunning() {
	defer DoSaves()
	for {
		select {
		case <-bOk:
			return
		case <-bDo:
			DoSaves()
		}
	}
}

func (r *runner) handleResult() {
	go DoRunning()

	for result := range r.recver {
		var m1 = map[string]interface{}{"ip": result.Answers, "subdomain": result.Subdomain, "tags": "subdomain"}
		oR <- &m1
		if 5000 <= len(oR) {
			bDo <- struct{}{}
		}
		for _, out := range r.options.Writer {
			_ = out.WriteDomainResult(result)
		}
		r.printStatus()
	}
	close(bOk)
}
