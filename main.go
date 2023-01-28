package main

import (
	util "github.com/hktalent/go-utils"
	myCmd "github.com/hktalent/ksubdomain/cmd/ksubdomain"
	"os"
)

func main() {
	os.RemoveAll("ksubdomain.yaml")
	util.DoInitAll()
	myCmd.Main()
	util.Wg.Wait()
	util.CloseAll()
}
