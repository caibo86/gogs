// -------------------------------------------
// @file      : flag.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午9:24
// -------------------------------------------

package config

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

var (
	revision  = "???"
	version   = "???"
	buildTime = "???"

	showVersion bool
	ServerType  string
	ServerID    int64
)

func Version() string {
	return version
}

func Revision() string {
	return revision
}

func BuildTime() string {
	return buildTime
}

func parseFlags() {
	flag.BoolVar(&showVersion, "version", false, "show version info")
	flag.StringVar(&ServerType, "serverType", "", "server type")
	flag.Int64Var(&ServerID, "serverID", 0, "server id")
	flag.Parse()

	if showVersion {
		fmt.Printf("Version: %s\n", Version())
		fmt.Printf("Revision: %s\n", Revision())
		fmt.Printf("Built at: %s\n", BuildTime())
		fmt.Printf("Powered by: %s\n", runtime.Version())
		os.Exit(0)
	}
}
