// -------------------------------------------
// @file      : lexer_test.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/15 下午6:15
// -------------------------------------------

package cblang

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewLexer(t *testing.T) {
	Convey("创建词法分析器", t, func() {
		var buff bytes.Buffer

		lexer := NewLexer("test", &buff)
		So(lexer, ShouldNotBeNil)
	})
}
