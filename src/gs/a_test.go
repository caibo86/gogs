// -------------------------------------------
// @file      : a_test.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/20 下午10:39
// -------------------------------------------

package gs

import (
	"encoding/gob"
)

func init() {
	gob.Register(&Car{})
	gob.Register(&Student{})
}
