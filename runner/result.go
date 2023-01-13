package runner

import (
	util "github.com/hktalent/go-utils"
)

// 记录日志到 大数据搜索引擎
func SendEsLog(o interface{}) {
	if 0 == len(util.EsUrl) {
		return
	}

	var m1 = map[string]interface{}{}
	if data, err := util.Json.Marshal(o); nil == err {
		if nil == util.Json.Unmarshal(data, &m1) {
			m1["tags"] = "subdomain"
			m1["tools"] = "ksubdomain"
			szId := util.GetSha1(&m1)
			util.SendReq(&m1, szId, "osint")
		}
	}
}

func (r *runner) handleResult() {
	for result := range r.recver {
		x1 := &result
		util.DefaultPool.Submit(func() {
			SendEsLog(x1)
		})
		for _, out := range r.options.Writer {
			_ = out.WriteDomainResult(result)
		}
		r.printStatus()
	}
}
