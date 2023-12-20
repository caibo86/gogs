package gg

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
