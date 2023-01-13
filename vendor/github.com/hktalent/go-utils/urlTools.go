package go_utils

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var mHttp = regexp.MustCompile(`(http[s]?:\/\/[^; $]+)`)

func GetUrlInfo(u string) {
	c1 := PipE.GetClient4Http2()
	c1.CheckRedirect = nil
	PipE.DoGetWithClient4SetHd(c1, u, "GET", nil, func(resp *http.Response, err error, szU string) {
		if nil == err {
			for _, x := range strings.Split("X-Cache-Hits,X-Cache,Via,Traceparent,Server-Timing,Strict-Transport-Security,Date,Paypal-Debug-Id,Set-Cookie,Etag,Content-Type,X-Timer,Accept-Ranges,Cache-Control,X-Xss-Protection,Vary,content-type,etag,paypal-debug-id,set-cookie,traceparent,X-Content-Type-Options,accept-ranges,via,date,strict-transport-security,x-served-by,x-cache,x-cache-hits,x-timer,server-timing,content-length", ",") {
				delete(resp.Header, x)
			}

			fmt.Printf("%+v", strings.Join(mHttp.FindAllString(resp.Header.Get("Content-Security-Policy"), -1), "\n"))

			fmt.Printf("\n\n%+v", resp.Header.Get("Content-Security-Policy"))
		}
	}, func() map[string]string {
		return map[string]string{}
	}, true)
}
