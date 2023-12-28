// -------------------------------------------
// @file      : gsc.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/19 下午3:53
// -------------------------------------------

package main

import (
	"flag"
	"gogs/base/gserrors"
	"gogs/base/gslang"
	log "gogs/base/logger"
)

const ignoreErr = "sync /dev/stdout: invalid argument"

func main() {
	log.Init(
		log.SetFilename("log/gsc.log"),
		log.SetIsOpenFile(true),
		log.SetIsAsync(true),
	)
	// 程序完成后关闭全局日志服务
	defer func() {
		e := recover()
		if e != nil {
			log.Errorf(gserrors.New(e.(error)).Error())
		}
		if err := log.Close(); err != nil {
			panic(gserrors.New(err))
		}
	}()
	// 解析命令行参数
	flag.StringVar(&moduleName, "module", "gogs", "golang module name")
	flag.Parse()
	log.Infof("Set module name: %s", moduleName)
	// packages := []string{"yf/platform/yfnet", "yf/platform/yfdocker"}
	var packages []string
	compiler := gslang.NewCompiler()
	log.Info("Start compiling packages: ", flag.Args())
	packages = append(packages, flag.Args()...)
	// 编译默认的两个包及命令行提供的目标包
	for _, name := range packages {
		log.Info("Compiling package: ", name)
		_, err := compiler.Compile(name)
		if err != nil {
			log.Errorf("compile package %s failed\n\t%s", name, err.Error())
			return
		}
	}

	// 访问者
	gen, err := NewGen4Go()
	if err != nil {
		log.Errorf("%s", err)
		return
	}
	err = compiler.Accept(gen)
	if err != nil {
		log.Errorf("%s", err)
		return
	}
	log.Info("Successfully compiled package: ", flag.Args())
	return
}
