// -------------------------------------------
// @file      : idgen.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/3 下午8:36
// -------------------------------------------

package misc

import (
	"gogs/base/gserrors"
	"gogs/base/idgen/snowflake_v2"
	log "gogs/base/logger"
)

var sf *snowflake_v2.SnowflakeV2 // 默认1

// InitIDGen 初始化ID生成器,必须调用
func InitIDGen(serverID int64) {
	sf = snowflake_v2.NewSnowflakeV2(serverID)
}

// NewID 生成新的ID 上层考虑到返回值可能为0
func NewID() int64 {
	if sf == nil {
		gserrors.Panic("init snowflake id generator first")
	}
	id, err := sf.Next()
	if err != nil {
		log.Errorf("CAUTION! create new uuid failed: %v", err)
	}
	return id
}
