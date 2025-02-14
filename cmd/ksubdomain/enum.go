package ksubdomain

import (
	"bufio"
	"context"
	_ "embed"
	util "github.com/hktalent/go-utils"
	"github.com/hktalent/ksubdomain/core"
	"github.com/hktalent/ksubdomain/core/dns"
	"github.com/hktalent/ksubdomain/core/gologger"
	"github.com/hktalent/ksubdomain/core/options"
	"github.com/hktalent/ksubdomain/runner"
	"github.com/hktalent/ksubdomain/runner/outputter"
	"github.com/hktalent/ksubdomain/runner/outputter/output"
	"github.com/hktalent/ksubdomain/runner/processbar"
	"github.com/urfave/cli/v2"
	"log"
	"math/rand"
	"os"
	"strings"
)

// https://raw.githubusercontent.com/publicsuffix/list/master/public_suffix_list.dat
// https://github.com/publicsuffix/list/
// https://data.iana.org/TLD/tlds-alpha-by-domain.txt
//
//go:embed tldsList.txt
var tldsList string
var tlds []string

func init() {
	// https://www.computerhope.com/jargon/num/domains.htm
	// https://www.namecheap.com/legal/domains/icann-fee/
	// https://www.namecheap.com/domains/full-tld-list/
	tlds = append(tlds, strings.Split(tldsList, "\n")...)
	tlds = util.RemoveDuplication_mapNoEmpy(tlds)
}

/*
兼容，前后有 * 情况
*/
func doName(s string) []string {
	var a []string
	if strings.HasSuffix(s, ".*") {
		s = s[0 : len(tlds)-2]
		for _, x := range tlds {
			a = append(a, s+x)
		}
	} else if strings.HasPrefix(s, "*.") {
		a = append(a, s[2:])
	} else {
		a = append(a, s)
	}
	return a
}

var enumCommand = &cli.Command{
	Name:    runner.EnumType,
	Aliases: []string{"e"},
	Usage:   "枚举域名",
	Flags: append(commonFlags, []cli.Flag{

		&cli.StringFlag{
			Name:     "domainList",
			Aliases:  []string{"dl"},
			Usage:    "从文件中指定域名",
			Required: false,
			Value:    "",
		},
		&cli.BoolFlag{
			Name:    "json",
			Aliases: []string{"j"},
			Usage:   "输出格式为json",
			Value:   false,
		},
		&cli.StringFlag{
			Name:     "filename",
			Aliases:  []string{"f"},
			Usage:    "字典路径",
			Required: false,
			Value:    "config/subdomain.txt", // $HOME/MyWork/scan4all/config/database/subdomain.txt
		},
		&cli.BoolFlag{
			Name:  "skip-wild",
			Usage: "跳过泛解析域名",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "ns",
			Usage: "读取域名ns记录并加入到ns解析器中",
			Value: false,
		},
		&cli.IntFlag{
			Name:    "level",
			Aliases: []string{"l"},
			Usage:   "枚举几级域名，默认为2，二级域名",
			Value:   2,
		},
		&cli.StringFlag{
			Name:    "level-dict",
			Aliases: []string{"ld"},
			Usage:   "枚举多级域名的字典文件，当level大于2时候使用，不填则会默认",
			Value:   "",
		},
	}...),
	Action: func(c *cli.Context) error {
		if c.NumFlags() == 0 {
			cli.ShowCommandHelpAndExit(c, "enum", 0)
		}
		var domains []string
		var writer []outputter.Output
		var processBar processbar.ProcessBar = &processbar.ScreenProcess{}
		var err error
		var domainTotal int = 0

		// handle domain
		if c.String("domain") != "" {
			if util.FileExists(c.String("domain")) {
				if data, err := os.ReadFile(c.String("domain")); nil == err {
					domains = append(domains, strings.Split(strings.TrimSpace(string(data)), "\n")...)
				}
			} else {
				domains = append(domains, c.String("domain"))
			}
		}
		if c.String("domainList") != "" {
			dl, err := core.LinesInFile(c.String("domainList"))
			if err != nil {
				gologger.Fatalf("读取domain文件失败:%s\n", err.Error())
			}
			domains = append(dl, domains...)
		}
		if c.Bool("stdin") {
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				domains = append(domains, scanner.Text())
			}
		}
		if c.Bool("skip-wild") {
			tmp := domains
			domains = []string{}
			for _, sub := range tmp {
				if !core.IsWildCard(sub) {
					domains = append(domains, sub)
				} else {
					gologger.Infof("域名:%s 存在泛解析,已跳过", sub)
				}
			}
		}

		var subdomainDict []string
		subdomainDict, err = core.LinesInFile(c.String("filename"))
		if err != nil {
			gologger.Fatalf("打开文件:%s 错误:%s", c.String("filename"), err.Error())
		} else {
			core.DefaultDomainList = subdomainDict
		}

		levelDict := c.String("level-dict")
		var levelDomains []string
		if levelDict != "" {
			dl, err := core.LinesInFile(levelDict)
			if err != nil {
				gologger.Fatalf("读取domain文件失败:%s,请检查--level-dict参数\n", err.Error())
			}
			levelDomains = dl
		} else if c.Int("level") > 2 {
			levelDomains = core.GetDefaultSubNextData()
		}

		render := make(chan string)
		util.DefaultPool.Submit(func() {
			defer close(render)
			log.Println("domains len", len(domains))
			for _, domain1 := range domains {
				a1 := doName(domain1)
				for _, domain := range a1 {
					for _, sub := range subdomainDict {
						var dd = sub + "." + domain
						if -1 < strings.Index(domain, "*") {
							dd = strings.ReplaceAll(domain, "*", sub)
						}

						//fmt.Printf("%s\r", dd)
						render <- dd
						// 这里应该加判断，如果本级无法dns，就没有必要到第3级以上
						if len(levelDomains) > 0 {
							for _, sub2 := range levelDomains {
								dd2 := sub2 + "." + dd
								render <- dd2
							}
						}
					}

				}
			}
		})
		domainTotal = len(subdomainDict) * len(domains)
		if len(levelDomains) > 0 {
			domainTotal *= len(levelDomains)
		}

		// 取域名的dns,加入到resolver中
		specialDns := make(map[string][]string)
		defaultResolver := options.GetResolvers(c.String("resolvers"))
		if c.Bool("ns") {
			for _, domain := range domains {
				nsServers, ips, err := dns.LookupNS(domain, defaultResolver[rand.Intn(len(defaultResolver))])
				if err != nil {
					continue
				}
				specialDns[domain] = ips
				gologger.Infof("%s ns:%v", domain, nsServers)
			}

		}
		onlyDomain := c.Bool("only-domain")

		if c.Bool("csv") {
			fileWriter, err := output.NewCsvOutImp(c.String("output"), onlyDomain, true)
			if err != nil {
				gologger.Fatalf(err.Error())
			}
			writer = append(writer, fileWriter)

		}
		if c.Bool("json") {
			fileWriter, err := output.NewJsonOutImp(c.String("output"), onlyDomain)
			if err != nil {
				gologger.Fatalf(err.Error())
			}
			writer = append(writer, fileWriter)
		}

		if c.String("output") != "" && 0 == len(writer) {
			fileWriter, err := output.NewFileOutput(c.String("output"), onlyDomain)
			if err != nil {
				gologger.Fatalf(err.Error())
			}
			writer = append(writer, fileWriter)
		}

		if c.Bool("not-print") {
			processBar = nil
		}

		screenWriter, err := output.NewScreenOutput(onlyDomain)
		if err != nil {
			gologger.Fatalf(err.Error())
		}
		writer = append(writer, screenWriter)

		opt := &options.Options{
			Rate:             options.Band2Rate(c.String("band")),
			Domain:           render,
			DomainTotal:      domainTotal,
			Resolvers:        defaultResolver,
			Silent:           c.Bool("silent"),
			TimeOut:          c.Int("timeout"),
			Retry:            c.Int("retry"),
			Method:           runner.VerifyType,
			DnsType:          c.String("dns-type"),
			Writer:           writer,
			ProcessBar:       processBar,
			SpecialResolvers: specialDns,
		}
		opt.Check()
		opt.EtherInfo = options.GetDeviceConfig()
		//fmt.Printf("%+v\n", opt.EtherInfo)
		ctx := context.Background()
		r, err := runner.New(opt)
		if err != nil {
			gologger.Fatalf("%s\n", err.Error())
			return nil
		}
		defer r.Close()
		r.RunEnumeration(ctx)

		return nil
	},
}
