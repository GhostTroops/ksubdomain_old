package ksubdomain

import (
	"github.com/hktalent/ksubdomain/core/options"
	"github.com/hktalent/ksubdomain/runner"
	"github.com/urfave/cli/v2"
)

var testCommand = &cli.Command{
	Name:  runner.TestType,
	Usage: "测试本地网卡的最大发送速度",
	Action: func(c *cli.Context) error {
		ether := options.GetDeviceConfig()
		runner.TestSpeed(ether)
		return nil
	},
}
