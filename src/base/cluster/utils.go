// -------------------------------------------
// @file      : utils.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/8 上午11:48
// -------------------------------------------

package cluster

import (
	log "gogs/base/logger"
	"strconv"
	"strings"
)

func GetIDByName(name string) int64 {
	s := strings.Split(name, ":")
	if len(s) != 2 {
		return 0
	}
	id, err := strconv.ParseInt(s[1], 10, 64)
	if err != nil {
		log.Errorf("invalid format for name(type:id): %s", name)
		return 0
	}
	return id
}
