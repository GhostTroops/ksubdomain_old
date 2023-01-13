package main

import (
	util "github.com/hktalent/go-utils"
	myCmd "github.com/hktalent/ksubdomain/cmd/ksubdomain"
)

func main() {
	util.DoInitAll()
	myCmd.Main()
	util.Wg.Wait()
	util.CloseAll()
}
