// -------------------------------------------
// @file      : a.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/20 下午9:07
// -------------------------------------------

package gs

// +k8s:deepcopy-gen=true
type ACar struct {
	ID      int64
	Name    string
	Price   float32
	Owner   *AStudent
	Code    AErr
	Size    [4]int32
	Drivers []*AStudent
	Attrs   map[string]string
	Color   AColor
}

// +k8s:deepcopy-gen=true
type AStudent struct {
	ID   int64
	Name string
	Age  int16
}

// +k8s:deepcopy-gen=true
type AErr uint16

// +k8s:deepcopy-gen=true
type AColor uint32
