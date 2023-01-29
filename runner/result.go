package runner

import (
	util "github.com/hktalent/go-utils"
	"strings"
)

func (r *runner) handleResult() {
	go util.DoRunning()
	defer util.CloseLogBigDb()
	var szSkp = "0.0.0.1"
	for result := range r.recver {
		if 1 == len(result.Answers) && szSkp == result.Answers[0] || -1 < strings.Index(result.Subdomain, szSkp) {
			continue
		}
		var m1 = map[string]interface{}{"ip": result.Answers, "subdomain": result.Subdomain, "tags": "subdomain"}
		go util.PushLog(&m1)
		for _, out := range r.options.Writer {
			_ = out.WriteDomainResult(result)
		}
		r.printStatus()
	}
}
