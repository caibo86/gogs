// -------------------------------------------
// @file      : const.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午6:27
// -------------------------------------------

package ast

// String 字面量字符串常量
type String struct {
	BaseExpr // 内嵌基本表达式实现
	Value    string
	Origin   string
}

func (s *String) OriginName() string {
	return s.Origin
}

// NewString 在代码节点内新建字符串常量
func (script *Script) NewString(value, origin string) *String {
	s := &String{
		Value:  value,
		Origin: origin,
	}
	s.Init("string", script)
	return s
}

// Float 字面量浮点数常量
type Float struct {
	BaseExpr // 内嵌基本表达式实现
	Value    float64
	Origin   string
}

func (f *Float) OriginName() string {
	return f.Origin
}

// NewFloat 在代码节点内新建浮点数常量
func (script *Script) NewFloat(value float64, origin string) *Float {
	f := &Float{
		Value:  value,
		Origin: origin,
	}
	f.Init("float", script)
	return f
}

// Int 字面量整数常量
type Int struct {
	BaseExpr // 内嵌基本表达式实现
	Value    int64
	Origin   string
}

func (i *Int) OriginName() string {
	return i.Origin
}

// NewInt 在代码节点内新建整数常量
func (script *Script) NewInt(value int64, origin string) *Int {
	i := &Int{
		Value:  value,
		Origin: origin,
	}
	i.Init("int", script)
	return i
}

// Bool 字面量布尔常量
type Bool struct {
	BaseExpr // 内嵌基本表达式实现
	Value    bool
	Origin   string
}

func (b *Bool) OriginName() string {
	return b.Origin
}

// NewBool 在代码节点内新建布尔常量
func (script *Script) NewBool(value bool, origin string) *Bool {
	b := &Bool{
		Value:  value,
		Origin: origin,
	}
	b.Init("bool", script)
	return b
}
