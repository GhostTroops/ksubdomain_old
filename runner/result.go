package runner

import (
	util "github.com/hktalent/go-utils"
)

// 记录日志到 大数据搜索引擎
func SendEsLog(m1 map[string]interface{}) {
	if 0 == len(util.EsUrl) {
		return
	}

	m1["tags"] = "subdomain"
	szId := util.GetSha1(&m1)
	util.SendReq(&m1, szId, "ksubdomain")

}

func (r *runner) handleResult() {
	for result := range r.recver {
		var m1 = map[string]interface{}{"ip": result.Answers, "subdomain": result.Subdomain}
		util.DoSyncFunc(func() {
			SendEsLog(m1)
		})
		for _, out := range r.options.Writer {
			_ = out.WriteDomainResult(result)
		}
		r.printStatus()
	}
}
