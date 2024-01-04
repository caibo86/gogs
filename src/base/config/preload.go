// -------------------------------------------
// @file      : preload.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午9:26
// -------------------------------------------

package config

import (
	"fmt"
	"os"
	"strconv"
)

func Preload() {
	parseFlags()
	err := os.Setenv("SERVER_ID", strconv.Itoa(int(ServerID)))
	if err != nil {
		panic(fmt.Errorf("setenv err:%s", err))
	}
}
