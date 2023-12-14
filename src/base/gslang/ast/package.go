// -------------------------------------------
// @file      : package.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/14 下午2:23
// -------------------------------------------

package ast

type Package struct {
	BaseNode
	Scripts map[string]*Script // 脚本列表
	Types   map[string]Expr
}
