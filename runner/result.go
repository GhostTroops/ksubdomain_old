package runner

import (
	util "github.com/hktalent/go-utils"
)

func (r *runner) handleResult() {
	go util.DoRunning()
	defer util.CloseLogBigDb()

	for result := range r.recver {
		var m1 = map[string]interface{}{"ip": result.Answers, "subdomain": result.Subdomain, "tags": "subdomain"}
		util.PushLog(&m1)
		for _, out := range r.options.Writer {
			_ = out.WriteDomainResult(result)
		}
		r.printStatus()
	}
}
