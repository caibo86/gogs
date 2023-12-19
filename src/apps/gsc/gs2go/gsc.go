// -------------------------------------------
// @file      : gsc.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/19 下午3:53
// -------------------------------------------

package main

import (
	"flag"
	"gogs/base/gslang"
	log "gogs/base/logger"
)

func main() {
	log.Init(
		log.SetIsOpenFile(true),
	)
	// 程序完成后关闭全局日志服务
	defer func() {
		e := recover()
		if e != nil {
			log.Error("inner error\n\t%s", e)
		}

		if err := log.Close(); err != nil {
			panic(err)
		}
	}()
	// 解析命令行参数
	flag.Parse()
	// packages := []string{"yf/platform/yfnet", "yf/platform/yfdocker"}
	var packages []string
	compiler := gslang.NewCompiler()
	packages = append(packages, flag.Args()...)
	// 编译默认的两个包及命令行提供的目标包
	for _, name := range packages {
		log.Debugf("开始编译包:%s", name)
		_, err := compiler.Compile(name)
		if err != nil {
			log.Errorf("compile package %s failed\n\t%s", name, err)
			return
		}
	}

	// 访问者
	gen, err := NewGen4Go()
	if err != nil {
		log.Error("inner error\n\t%s", err)
		return
	}
	log.Debug("生成器")
	err = compiler.Accept(gen)
	log.Debug("完成")
	if err != nil {
		log.Error("inner error\n\t%s", err)
		return
	}
}