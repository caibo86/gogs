// -------------------------------------------
// @file      : node.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/29 下午6:52
// -------------------------------------------

package gsnet

import (
	"fmt"
)

// IChannel 通道接口
type IChannel interface {
	fmt.Stringer
	Write(*Message) error
	Status() Status //
}
