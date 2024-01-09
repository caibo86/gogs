// -------------------------------------------
// @file      : gen4go.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/19 下午4:00
// -------------------------------------------

package main

import (
	"bytes"
	"fmt"
	"gogs/base/cberrors"
	"gogs/base/cblang"
	"gogs/base/cblang/ast"
	log "gogs/base/logger"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var moduleName string

// 包名映射的引入包的go代码
var packageMapping = map[string]string{
	"network.":  `import "gogs/base/cluster/network"`,
	"cberrors.": `import "gogs/base/cberrors"`,
	"cluster.":  `import "gogs/base/cluster"`,
	"config.":   `import "gogs/base/config"`,
	"bytes.":    `import "bytes"`,
	"fmt.":      `import "fmt"`,
	"time.":     `import "time"`,
	"bits.":     `import "math/bits"`,
	"io":        `import "io"`,
}

// cblang内置类型对应的golang表示
var keyMapping = map[string]string{
	".cblang.Byte":    "byte",
	".cblang.Bytes":   "[]byte",
	".cblang.Int8":    "int8",
	".cblang.Uint8":   "uint8",
	".cblang.Int16":   "int16",
	".cblang.Uint16":  "uint16",
	".cblang.Int32":   "int32",
	".cblang.Uint32":  "uint32",
	".cblang.Int64":   "int64",
	".cblang.Uint64":  "uint64",
	".cblang.Float32": "float32",
	".cblang.Float64": "float64",
	".cblang.Bool":    "bool",
	".cblang.String":  "string",
	"Byte":            "byte",
	"Bytes":           "[]byte",
	"Int8":            "int8",
	"Uint8":           "uint8",
	"Int16":           "int16",
	"Uint16":          "uint16",
	"Int32":           "int32",
	"Uint32":          "uint32",
	"Int64":           "int64",
	"Uint64":          "uint64",
	"Float32":         "float32",
	"Float64":         "float64",
	"Bool":            "bool",
	"String":          "string",
}

// cblang内置类型的默认值对应的golang表示
var defaultVal = map[string]string{
	".cblang.Byte":    "byte(0)",
	".cblang.Bytes":   "[]byte(nil)",
	".cblang.Int8":    "int8(0)",
	".cblang.Uint8":   "uint8(0)",
	".cblang.Int16":   "int16(0)",
	".cblang.Uint16":  "uint16(0)",
	".cblang.Int32":   "int32(0)",
	".cblang.Uint32":  "uint32(0)",
	".cblang.Int64":   "int64(0)",
	".cblang.Uint64":  "uint64(0)",
	".cblang.Float32": "float32(0)",
	".cblang.Float64": "float64(0)",
	".cblang.Bool":    "false",
	".cblang.String":  "\"\"",
}

// writeMapping 写入方法映射
var writeMapping = map[string]string{
	"Bool":    "network.WriteBool",
	"Byte":    "network.WriteByte",
	"Bytes":   "network.WriteBytes",
	"Int8":    "network.WriteInt8",
	"Uint8":   "network.WriteUint8",
	"Int16":   "network.WriteInt16",
	"Uint16":  "network.WriteUint16",
	"Int32":   "network.WriteInt32",
	"Uint32":  "network.WriteUint32",
	"Float32": "network.WriteFloat32",
	"Int64":   "network.WriteInt64",
	"Uint64":  "network.WriteUint64",
	"Float64": "network.WriteFloat64",
	"String":  "network.WriteString",
}

// readMapping 读方法映射
var readMapping = map[string]string{
	"Bool":    "network.ReadBool",
	"Byte":    "network.ReadByte",
	"Bytes":   "network.ReadBytes",
	"Int8":    "network.ReadInt8",
	"Uint8":   "network.ReadUint8",
	"Int16":   "network.ReadInt16",
	"Uint16":  "network.ReadUint16",
	"Int32":   "network.ReadInt32",
	"Uint32":  "network.ReadUint32",
	"Float32": "network.ReadFloat32",
	"Int64":   "network.ReadInt64",
	"Uint64":  "network.ReadUint64",
	"Float64": "network.ReadFloat64",
	"String":  "network.ReadString",
}

// compareMapping 比较方法映射
var compareMapping = map[string]string{
	"Bool":    `%s`,
	"Byte":    `%s != 0`,
	"Bytes":   `len(%s) > 0`,
	"Int8":    `%s != 0`,
	"Uint8":   `%s != 0`,
	"Int16":   `%s != 0`,
	"Uint16":  `%s != 0`,
	"Int32":   `%s != 0`,
	"Uint32":  `%s != 0`,
	"Float32": `%s != 0`,
	"Int64":   `%s != 0`,
	"Uint64":  `%s != 0`,
	"Float64": `%s != 0`,
	"String":  `len(%s) > 0`,
}

// Gen4Go golang代码生成器
type Gen4Go struct {
	ast.EmptyVisitor                    // 内嵌空访问者
	buff             bytes.Buffer       // 缓冲区
	tpl              *template.Template // 模板
	gen              bool
}

// NewGen4Go 新建一个golang代码生成器
func NewGen4Go() (gen *Gen4Go, err error) {
	gen = &Gen4Go{}
	functions := template.FuncMap{
		"symbol":              strings.Title,
		"pos":                 cblang.Pos,
		"typeName":            gen.typeName,
		"readType":            gen.readType,
		"writeType":           gen.writeType,
		"defaultVal":          gen.defaultVal,
		"sovFunc":             gen.sovFunc,
		"calTypeSize":         gen.calTypeSize,
		"copyType":            gen.copyType,
		"printComments":       gen.printComments,
		"printCommentsToLine": gen.printCommentsToLine,
		"params":              gen.params,
		"paramsName":          gen.paramsName,
		"paramsW":             gen.paramsW,
		"genReadWrite":        gen.genReadWrite,
	}
	gen.tpl, err = template.New("golang").Funcs(functions).Parse(tpl4go)
	return
}

// genReadWrite 是否为Struct生成读写方法
func (gen *Gen4Go) genReadWrite(table *ast.Table) bool {
	ret := false
	if cblang.IsStruct(table) && cblang.ReadWrite(table) {
		ret = true
	}
	log.Debugf("struct: %s, gen Read/Write func? : %v", table.Name(), ret)
	return ret
}

// calTypeSize 根据类型获取计算大小的函数
func (gen *Gen4Go) calTypeSize(field ast.IField) string {
	switch field.GetType().(type) {
	case *ast.TypeRef:
		ref := field.GetType().(*ast.TypeRef)
		switch ref.Ref.Name() {
		case "Bool":
			return fmt.Sprintf(
				`if %s {
					n += 3
				}`, field.GetName(),
			)
		case "Byte", "Int8", "Uint8":
			return fmt.Sprintf(
				`if %s != 0 {
					n += 3
				}`,
				field.GetName())
		case "Int16", "Uint16":
			return fmt.Sprintf(
				`if %s != 0 {
					n += 4
				}`,
				field.GetName())
		case "Int32", "Uint32", "Float32":
			return fmt.Sprintf(
				`if %s != 0 {
					n += 6
				}`,
				field.GetName())
		case "Int64", "Uint64", "Float64":
			return fmt.Sprintf(
				`if %s != 0 {
					n += 10
				}`,
				field.GetName())
		case "String", "Bytes":
			return fmt.Sprintf(
				`l = len(%s)
					if l > 0 {
					n += 6 + l 
				}`,
				field.GetName())
		default:
			switch ref.Ref.(type) {
			case *ast.Enum:
				return fmt.Sprintf(
					`if %s != 0 {
						n += 6
					}`,
					field.GetName())
			case *ast.Table:
				return fmt.Sprintf(
					`if %s != nil {
						l = %s.Size()
						n += 6 + l
					}`,
					field.GetName(), field.GetName())
			default:
				cberrors.Panic("not here %s", field.GetType().Name())
			}
		}
	case *ast.Slice, *ast.Array:
		// 切片和数组
		var ref *ast.TypeRef
		if slice, ok := field.GetType().(*ast.Slice); ok {
			ref = slice.Element.(*ast.TypeRef)
		} else {
			ref = field.GetType().(*ast.Array).Element.(*ast.TypeRef)
		}
		switch ref.Ref.Name() {
		case "Byte", "Int8", "Uint8", "Bool":
			return fmt.Sprintf(
				`l = len(%s)
				if l > 0 {
					n += 6 + l
				}`,
				field.GetName())
		case "Uint16", "Int16":
			return fmt.Sprintf(
				`l = len(%s)
				if l > 0 {
					n += 6 + l * 2
				}`,
				field.GetName())
		case "Uint32", "Int32", "Float32":
			return fmt.Sprintf(
				`l = len(%s)
				if l > 0 {
					n += 6 + l * 4
				}`,
				field.GetName())
		case "Uint64", "Int64", "Float64":
			return fmt.Sprintf(
				`l = len(%s)
				if l > 0 {
					n += 6 + l * 8
				}`,
				field.GetName())
		case "String", "Bytes":
			return fmt.Sprintf(
				`if len(%s) > 0 {
					n += 6
					for _, s := range %s {
						l = len(s)
						n += 4 + l
					}
				}`,
				field.GetName(), field.GetName())
		default:
			switch ref.Ref.(type) {
			case *ast.Enum:
				return fmt.Sprintf(
					`l = len(%s)
					if l > 0 {
						n += 6 + l * 4
					}`,
					field.GetName())
			case *ast.Table:
				return fmt.Sprintf(
					`l = len(%s)
					if l > 0 {
						n += 6
						for _, e := range %s {
							n += 4 + e.Size()
						}
					}`,
					field.GetName(), field.GetName())
			default:
				cberrors.Panic("not here %s", field.GetType().Name())
			}
		}
	case *ast.Map:
		// 字典
		hash := field.GetType().(*ast.Map)
		var keyStr string
		var valStr string
		ref := hash.Key.(*ast.TypeRef)
		switch ref.Ref.Name() {
		case "Byte", "Int8", "Uint8", "Bool":
			keyStr = "1"
		case "Uint16", "Int16":
			keyStr = "2"
		case "Uint32", "Int32":
			keyStr = "4"
		case "Uint64", "Int64":
			keyStr = "8"
		case "String":
			keyStr = "4 + len(k)"
		default:
			switch ref.Ref.(type) {
			case *ast.Enum:
				keyStr = "4"
			default:
				cberrors.Panic("map key can only be int or string, %s not supported", hash.Key.Name())
			}
		}
		ref = hash.Value.(*ast.TypeRef)
		switch ref.Ref.Name() {
		case "Byte", "Int8", "Uint8", "Bool":
			valStr = "1"
		case "Uint16", "Int16":
			valStr = "2"
		case "Uint32", "Int32", "float32":
			valStr = "4"
		case "Uint64", "Int64", "float64":
			valStr = "8"
		case "String", "Bytes":
			valStr = "4 + len(v)"
		default:
			switch ref.Ref.(type) {
			case *ast.Enum:
				valStr = "4"
			case *ast.Table:
				valStr = "4 + v.Size()"
			default:
				cberrors.Panic("map value %s not supported", hash.Value.Name())
			}
		}
		return fmt.Sprintf(
			`if len(%s) > 0 {
						n += 6
						for k, v := range %s {
							_ = k
							_ = v
							n += %s
							n += %s
						}
					}`,
			field.GetName(),
			field.GetName(),
			keyStr,
			valStr)
	}
	cberrors.Panic("not here")
	return "unknown"
}

// writeType 根据字段类型生成写入函数 TODO
func (gen *Gen4Go) writeType(field ast.IField) string {
	switch field.GetType().(type) {
	case *ast.TypeRef:
		var str string
		var ok bool
		ref := field.GetType().(*ast.TypeRef)
		if str, ok = writeMapping[ref.Ref.Name()]; ok {
			return fmt.Sprintf(
				`if `+compareMapping[ref.Ref.Name()]+` {
					i = network.WriteFieldID(data, i, %d)
					i = %s(data, i, %s)
				}`,
				field.GetName(), field.GetID(), str, field.GetName())
		} else {
			switch ref.Ref.(type) {
			case *ast.Enum:
				return fmt.Sprintf(
					`if %s != 0 {
						i = network.WriteFieldID(data, i, %d)
						i = network.WriteEnum(data, i, int32(%s))
					}`,
					field.GetName(), field.GetID(), field.GetName())
			case *ast.Table:
				return fmt.Sprintf(
					`if %s != nil {
						i = network.WriteFieldID(data, i, %d)
						size := %s.Size()
						i = network.WriteUint32(data, i, uint32(size))
						%s.MarshalToSizedBuffer(data[i:])
						i += size
					}`,
					field.GetName(), field.GetID(), field.GetName(), field.GetName())
			default:
				cberrors.Panic("not here %s", field.GetType().Name())
			}
		}
	case *ast.Slice, *ast.Array:
		// 切片和数组
		var ref *ast.TypeRef
		isSlice := false
		if slice, ok := field.GetType().(*ast.Slice); ok {
			ref = slice.Element.(*ast.TypeRef)
			isSlice = true
		} else {
			ref = field.GetType().(*ast.Array).Element.(*ast.TypeRef)
		}

		var str string
		var ok bool
		if str, ok = writeMapping[ref.Ref.Name()]; ok {
			if isSlice && ref.Ref.Name() == "Byte" {
				return fmt.Sprintf(
					`if len(%s) > 0 {
						i = network.WriteFieldID(data, i , %d)
						i = network.WriteBytes(data, i, %s)
					}`,
					field.GetName(),
					field.GetID(),
					field.GetName())
			} else {
				str = fmt.Sprintf("i = %s(data, i, e)", str)
			}
		} else {
			switch ref.Ref.(type) {
			case *ast.Enum:
				str = `i = network.WriteEnum(data, i, int32(e))`
			case *ast.Table:
				str = `size := e.Size()
						i = network.WriteUint32(data, i, uint32(size))
						if e != nil {
							e.MarshalToSizedBuffer(data[i:])
						}
						i += size`
			default:
				cberrors.Panic("not here %s", field.GetType().Name())
			}
		}
		if isSlice {
			return fmt.Sprintf(
				`if len(%s) > 0 {
					i = network.WriteFieldID(data, i , %d)
					i = network.WriteUint32(data, i, uint32(len(%s)))
					for _, e := range %s {
						%s
					}
				}`,
				field.GetName(),
				field.GetID(),
				field.GetName(),
				field.GetName(),
				str)
		} else {
			return fmt.Sprintf(
				`i = network.WriteFieldID(data, i , %d)
					i = network.WriteUint32(data, i, uint32(len(%s)))
					for _, e := range %s {
						%s
					}`,
				field.GetID(),
				field.GetName(),
				field.GetName(),
				str)
		}

	case *ast.Map:
		// 字典
		hash := field.GetType().(*ast.Map)
		var keyStr string
		var valStr string
		ref := hash.Key.(*ast.TypeRef)
		var ok bool
		if keyStr, ok = writeMapping[ref.Ref.Name()]; ok {
			keyStr = fmt.Sprintf("i = %s(data, i, k)", keyStr)
		} else {
			switch ref.Ref.(type) {
			case *ast.Enum:
				keyStr = `i = network.WriteEnum(data, i, int32(k))`
			default:
				cberrors.Panic("map key can only be int or string, %s not supported", hash.Key.Name())
			}
		}
		ref = hash.Value.(*ast.TypeRef)
		if valStr, ok = writeMapping[ref.Ref.Name()]; ok {
			valStr = fmt.Sprintf("i = %s(data, i, v)", valStr)
		} else {
			switch ref.Ref.(type) {
			case *ast.Enum:
				valStr = `i = network.WriteEnum(data, i, int32(v))`
			case *ast.Table:
				valStr = `size := v.Size()
						i = network.WriteUint32(data, i, uint32(size))
						if v != nil {
							v.MarshalToSizedBuffer(data[i:])
						}
						i += size`
			default:
				cberrors.Panic("map value %s not supported", hash.Value.Name())
			}
		}
		return fmt.Sprintf(
			`if len(%s) > 0 {
						i = network.WriteFieldID(data, i, %d)
						i = network.WriteUint32(data, i, uint32(len(%s)))
						for k, v := range %s {
							%s
							%s
						}
					}`,
			field.GetName(),
			field.GetID(),
			field.GetName(),
			field.GetName(),
			keyStr,
			valStr)
	}
	cberrors.Panic("not here")
	return "unknown"
}

// readType 根据字段类型生成读取函数
func (gen *Gen4Go) readType(field ast.IField) string {
	switch field.GetType().(type) {
	case *ast.TypeRef:
		ref := field.GetType().(*ast.TypeRef)
		var ok bool
		var str string
		if str, ok = readMapping[ref.Ref.Name()]; ok {
			return fmt.Sprintf("i, %s = %s(data, i)", field.GetName(), str)
		} else {
			switch ref.Ref.(type) {
			case *ast.Enum:
				return fmt.Sprintf(`var v int32
						i, v = network.ReadEnum(data, i)
						%s = %s(v)`,
					field.GetName(), gen.typeName(ref))
			case *ast.Table:
				return fmt.Sprintf(`var size uint32
						i, size = network.ReadUint32(data, i)
						if %s == nil {
							%s = %s
						}
						if err = %s.Unmarshal(data[i:i+int(size)]); err != nil {
							return
						} 
						i += int(size)`,
					field.GetName(), field.GetName(), gen.defaultVal(ref), field.GetName())
			default:
				cberrors.Panic("not here %s", field.GetType().Name())
			}
		}
	case *ast.Slice, *ast.Array:
		// 切片和数组
		var ref *ast.TypeRef
		var isSlice bool
		if slice, ok := field.GetType().(*ast.Slice); ok {
			isSlice = true
			ref = slice.Element.(*ast.TypeRef)
		} else {
			ref = field.GetType().(*ast.Array).Element.(*ast.TypeRef)
		}
		var str string
		var ok bool
		if str, ok = readMapping[ref.Ref.Name()]; ok {
			if isSlice && ref.Ref.Name() == "Byte" {
				return fmt.Sprintf(`i, %s = network.ReadBytes(data, i)`,
					field.GetName())
			} else {
				str = fmt.Sprintf("i, %s[j] = %s(data, i)", field.GetName(), str)
			}
		} else {
			switch ref.Ref.(type) {
			case *ast.Enum:
				str = fmt.Sprintf(`var v int32
				i, v = network.ReadEnum(data, i)
				%s[j] = %s(v)`,
					field.GetName(), gen.typeName(ref))
			case *ast.Table:
				str = fmt.Sprintf(
					`var size uint32
						i, size = network.ReadUint32(data, i)
						if size > 0 {
							%s[j] = %s
							if err = %s[j].Unmarshal(data[i:i+int(size)]); err != nil {
								return
							}
						}
						i += int(size)`, field.GetName(), gen.defaultVal(ref), field.GetName())
			default:
				cberrors.Panic("not here %s", field.GetType().Name())
			}
		}
		if isSlice {
			return fmt.Sprintf(
				`var length uint32
				i, length = network.ReadUint32(data, i)
				%s = make([]%s, length)
				for j := uint32(0); j < length; j++ {
					%s
				}`,
				field.GetName(), gen.typeName(ref), str)
		} else {
			return fmt.Sprintf(
				`var length uint32
				i, length = network.ReadUint32(data, i)
				for j := uint32(0); j < length; j++ {
					%s
				}`,
				str)
		}
	case *ast.Map:
		// 字典
		hash := field.GetType().(*ast.Map)
		var keyStr string
		var valStr string
		var ok bool
		ref := hash.Key.(*ast.TypeRef)
		if keyStr, ok = readMapping[ref.Ref.Name()]; ok {
			keyStr = fmt.Sprintf(`var k %s
					i, k = %s(data, i)`, keyMapping[ref.Ref.Name()], keyStr)
		} else {
			switch ref.Ref.(type) {
			case *ast.Enum:
				keyStr = fmt.Sprintf(`var k1 int32
					i, k1 = network.ReadEnum(data, i)
					k := %s(k1)`, gen.typeName(hash.Key))
			default:
				cberrors.Panic("map key can only be int or string, %s not supported", hash.Key.Name())
			}
		}
		ref = hash.Value.(*ast.TypeRef)
		if valStr, ok = readMapping[ref.Ref.Name()]; ok {
			valStr = fmt.Sprintf(`var v %s
					i, v = %s(data, i)`, keyMapping[ref.Ref.Name()], valStr)
		} else {
			switch ref.Ref.(type) {
			case *ast.Enum:
				valStr = fmt.Sprintf(`var v1 int32
					i, v1 = network.ReadEnum(data, i)
					v := %s(v1)`, gen.typeName(hash.Value))
			case *ast.Table:
				valStr = fmt.Sprintf(`var size uint32
						var v %s
						i, size = network.ReadUint32(data, i)
						if size > 0 {
							v = %s
							if err = v.Unmarshal(data[i:i+int(size)]); err != nil {
								return
							}
						}
						i += int(size)`, gen.typeName(hash.Value), gen.defaultVal(hash.Value))
			default:
				cberrors.Panic("map value %s not supported", hash.Value.Name())
			}
		}
		return fmt.Sprintf(
			`var length uint32
					i, length = network.ReadUint32(data, i)
					if %s == nil{
						%s = make(map[%s]%s)
					}
					for j := uint32(0); j < length; j++ {
						%s
						%s
						%s[k] = v
					}`,
			field.GetName(), field.GetName(),
			gen.typeName(hash.Key),
			gen.typeName(hash.Value),
			keyStr, valStr, field.GetName())
	}
	cberrors.Panic("not here")
	return "unknown"
}

// copyType 根据字段类型生成复制代码
func (gen *Gen4Go) copyType(field *ast.Field) string {
	switch field.Type.(type) {
	case *ast.TypeRef:
		ref := field.Type.(*ast.TypeRef)
		if _, ok := keyMapping[ref.Ref.Name()]; ok {
			return ""
		}
		switch ref.Ref.(type) {
		case *ast.Table:
			return fmt.Sprintf(
				`if m.%s != nil {
					in, out := &m.%s, &out.%s
					*out = %s	
					(*in).CopyInto(*out)
				}`,
				field.Name(), field.Name(), field.Name(), gen.defaultVal(ref))
		}
	case *ast.Slice:
		slice := field.Type.(*ast.Slice)
		ref := slice.Element.(*ast.TypeRef)
		if _, ok := keyMapping[ref.Ref.Name()]; ok {
			return fmt.Sprintf(`if m.%s != nil {
					in, out := &m.%s, &out.%s
					*out = make([]%s, len(*in))
					copy(*out,*in)
				}`,
				field.Name(), field.Name(), field.Name(), gen.typeName(ref))
		}
		switch ref.Ref.(type) {
		case *ast.Table:
			return fmt.Sprintf(`if m.%s != nil {
					in, out := &m.%s, &out.%s
					*out = make([]%s, len(*in))
					for i:= range *in {
						if (*in)[i] != nil {
							in, out := &(*in)[i], &(*out)[i]	
							*out = %s
							(*in).CopyInto(*out)
						}
					}
				}`,
				field.Name(), field.Name(), field.Name(), gen.typeName(ref), gen.defaultVal(ref))
		default:
			return fmt.Sprintf(`if m.%s != nil {
					in, out := &m.%s, &out.%s
					*out = make([]%s, len(*in))
					copy(*out,*in)
				}`,
				field.Name(), field.Name(), field.Name(), gen.typeName(ref))
		}
	case *ast.Array:
		array := field.Type.(*ast.Array)
		ref := array.Element.(*ast.TypeRef)
		if _, ok := keyMapping[ref.Ref.Name()]; ok {
			return fmt.Sprintf(`m.%s = out.%s`,
				field.Name(), field.Name())
		}
		switch ref.Ref.(type) {
		case *ast.Table:
			return fmt.Sprintf(`	for i, v:= range m.%s {
					if v != nil {
						out.%s[i] = %s
						v.CopyInto(out.%s[i])
					} else {
						out.%s[i] = nil
					}
				}`,
				field.Name(), field.Name(),
				gen.defaultVal(ref),
				field.Name(), field.Name())
		default:
			return fmt.Sprintf(`m.%s = out.%s`,
				field.Name(), field.Name())
		}
	case *ast.Map:
		// 字典
		hash := field.Type.(*ast.Map)
		keyRef := hash.Key.(*ast.TypeRef)
		valRef := hash.Value.(*ast.TypeRef)
		if _, ok := keyMapping[valRef.Ref.Name()]; ok {
			return fmt.Sprintf(`if m.%s != nil {
					in, out := &m.%s, &out.%s
					*out = make(map[%s]%s, len(*in))
					for k, v := range *in {
						(*out)[k] = v
					}
				}`,
				field.Name(), field.Name(), field.Name(),
				gen.typeName(keyRef), gen.typeName(valRef))
		}
		switch valRef.Ref.(type) {
		case *ast.Table:
			return fmt.Sprintf(`if m.%s != nil {
					in, out := &m.%s, &out.%s
					*out = make(map[%s]%s, len(*in))
					for k, v := range *in {
						var outVal %s
						if v == nil {
							(*out)[k] = nil
						} else {
							in, out := &v, &outVal
							*out = %s
							(*in).CopyInto(*out)
						}
						(*out)[k] = outVal
					}
				}`,
				field.Name(), field.Name(), field.Name(),
				gen.typeName(keyRef), gen.typeName(valRef),
				gen.typeName(valRef), gen.defaultVal(valRef))
		default:
			return fmt.Sprintf(`if m.%s != nil {
					in, out := &m.%s, &out.%s
					*out = make(map[%s]%s, len(*in))
					for key, val := range *in {
						(*out)[key] = val
					}
				}`,
				field.Name(), field.Name(), field.Name(),
				gen.typeName(keyRef), gen.typeName(valRef))
		}

	}
	return ""
}

// defaultVal 根据类型取其默认值表达式
func (gen *Gen4Go) defaultVal(expr ast.Expr) string {
	switch expr.(type) {
	case *ast.TypeRef:
		// 内置类型
		if val, ok := defaultVal[expr.Name()]; ok {
			return val
		}
		ref := expr.(*ast.TypeRef)
		// 枚举
		if enum, ok := ref.Ref.(*ast.Enum); ok {
			if _, ok := expr.Script().Imports[ref.NamePath[0]]; ok {
				return fmt.Sprintf(
					"%s.%s%s",
					ref.NamePath[0],
					enum,
					enum.Default,
				)
			}
			return fmt.Sprintf("%s%s", enum, enum.Default)
		}
		// 自定义类型
		if _, ok := expr.Script().Imports[ref.NamePath[0]]; ok {
			return fmt.Sprintf("%s.New%s()", ref.NamePath[0], strings.Title(ref.NamePath[1]))
		}
		return fmt.Sprintf("New%s()", strings.Title(ref.NamePath[0]))
	case *ast.Array:
		// 数组
		array := expr.(*ast.Array)
		var buff bytes.Buffer
		if err := gen.tpl.ExecuteTemplate(&buff, "arrayInit", array); err != nil {
			cberrors.Panic(err.Error())
		}
		return buff.String()
	case *ast.Slice:
		return "nil"
	case *ast.Map:
		// 字典
		return fmt.Sprintf("make(%s)", gen.typeName(expr))
	}
	cberrors.Panic("not here")
	return "unknown"
}

// params 根据参数生成参数列表
func (gen *Gen4Go) params(token string, params []*ast.Param, withErr bool) string {
	if len(params) == 0 {
		if withErr {
			return "err error"
		}
		return ""
	}
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("%s%d %s", token, 0, gen.typeName(params[0].Type)))
	for i := 1; i < len(params); i++ {
		buff.WriteString(fmt.Sprintf(", %s%d %s", token, i, gen.typeName(params[i].Type)))
	}
	if withErr {
		buff.WriteString(", err error")
	}
	return buff.String()
}

// paramsName 根据参数生成名字列表
func (gen *Gen4Go) paramsName(token string, params []*ast.Param, withErr bool) string {
	if len(params) == 0 {
		if withErr {
			return "err"
		}
		return ""
	}
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("%s0", token))
	for i := 1; i < len(params); i++ {
		buff.WriteString(fmt.Sprintf(",%s%d", token, i))
	}
	if withErr {
		buff.WriteString(", err")
	}
	return buff.String()
}

// paramsW 根据参数生成函数声明的入参列表
func (gen *Gen4Go) paramsW(token string, params []*ast.Param) string {
	if len(params) == 0 {
		return ""
	}
	var buff bytes.Buffer
	for i := 0; i < len(params); i++ {
		buff.WriteString(fmt.Sprintf(",\"%s\",%s%d",
			strings.TrimLeft(params[i].Type.Name(), "."), token, i))
	}
	return buff.String()
}

// typeName 取ast中的类型对应的golang表示
func (gen *Gen4Go) typeName(expr ast.Expr) string {
	switch expr.(type) {
	case *ast.TypeRef:
		// 内置类型引用
		if key, ok := keyMapping[expr.Name()]; ok {
			return key
		}
		ref := expr.(*ast.TypeRef)
		// 枚举
		if _, ok := ref.Ref.(*ast.Enum); ok {
			return strings.TrimLeft(expr.Name(), ".")
		}
		// 自定义类型引用均用指针
		return "*" + strings.TrimLeft(expr.Name(), ".")
	case *ast.Array:
		// 数组
		array := expr.(*ast.Array)
		return fmt.Sprintf("[%d]%s", array.Length, gen.typeName(array.Element))
	case *ast.Slice:
		// 切片
		slice := expr.(*ast.Slice)
		return fmt.Sprintf("[]%s", gen.typeName(slice.Element))
	case *ast.Map:
		// 字典
		hash := expr.(*ast.Map)
		return fmt.Sprintf("map[%s]%s", gen.typeName(hash.Key), gen.typeName(hash.Value))
	}
	cberrors.Panic("unknown golang typeName: %s\n\t%s", expr, cblang.Pos(expr))
	return "unknown"
}

// sovFunc sov函数名字
func (gen *Gen4Go) sovFunc(script *ast.Script) string {
	name := script.Name()
	ss := strings.Split(name, ".")
	ret := fmt.Sprintf("sov%s", strings.Title(ss[0]))
	return ret
}

// printComments 打印注释
func (gen *Gen4Go) printComments(node ast.Node) string {
	var ret string
	comments := cblang.Comments(node)
	if len(comments) > 0 {
		ret += "\n"
		for i, comment := range comments {
			value := comment.Value.(string)
			value = strings.TrimLeft(value, " ")
			if i == len(comments)-1 {
				ret += fmt.Sprintf("//%s", comment.Value)
			} else {
				ret += fmt.Sprintf("//%s\n", comment.Value)
			}
		}
	}
	return ret
}

// printComments 打印注释到一行
func (gen *Gen4Go) printCommentsToLine(node ast.Node) string {
	var ret string
	comments := cblang.Comments(node)
	if len(comments) > 0 {
		ret = "//"
		for _, comment := range comments {
			value := comment.Value.(string)
			value = strings.TrimLeft(value, " ")
			ret += comment.Value.(string)
		}
	}
	return ret
}

// writeFile 将代码节点对应的golang代码写入到文件
func (gen *Gen4Go) writeFile(script *ast.Script, bytes []byte) {
	fullPath, ok := cblang.FilePath(script)
	if !ok {
		cberrors.Panic("compile must bind file path to script")
	}
	// 写入文件名为 源文件名+.go
	fullPath += ".go"
	err := os.WriteFile(fullPath, bytes, 0644)
	if err != nil {
		cberrors.Panic(err.Error())
	}
	log.Infof("Write to file successfully: %s success", fullPath)

	cmd := exec.Command("goimports", "-w", fullPath)
	_, err = cmd.Output()
	if err != nil {
		log.Errorf("goimports format err: %s, check if goimports installed:\n\tgo install golang.org/x/tools/cmd/goimports@latest", err)
	}
}

// VisitPackage 访问包
func (gen *Gen4Go) VisitPackage(pkg *ast.Package) ast.Node {
	// 内置cblang包则直接返回
	if pkg.Name() == "base/cblang" {
		return pkg
	}
	// 轮询访问包中代码节点
	for _, script := range pkg.Scripts {
		script.Accept(gen)
	}
	return pkg
}

// VisitScript 访问代码
func (gen *Gen4Go) VisitScript(script *ast.Script) ast.Node {
	gen.buff.Reset()
	// 默认的一些代码
	if err := gen.tpl.ExecuteTemplate(&gen.buff, "script", script); err != nil {
		cberrors.Panic(err.Error())
	}

	// 轮询访问代码中的所有类型 Enum Struct Table Service
	// 按顺序生成
	for _, t := range script.Types {
		if _, ok := t.(*ast.Enum); ok {
			t.Accept(gen)
		}
	}
	for _, t := range script.Types {
		if _, ok := t.(*ast.Table); ok {
			t.Accept(gen)
		}
	}
	var hasService bool
	for _, t := range script.Types {
		if _, ok := t.(*ast.Service); ok {
			hasService = true
			t.Accept(gen)
		}
	}
	// 代码中有类型
	if gen.buff.Len() > 0 {
		var buff bytes.Buffer
		filename, _ := cblang.FilePath(script)
		filename += ".go"
		filename = filepath.Base(filename)
		// 写入额外信息
		buff.WriteString(fmt.Sprintf(
			`// -------------------------------------------
// @file      : %s
// @author    : generated by cblang complier, do not edit
// @contact   : caibo923@gmail.com
// @time      : %s
// -------------------------------------------

`, filename, time.Now().Format(time.RFC3339)))

		// 写入包声明
		buff.WriteString(fmt.Sprintf("package %s\n", filepath.Base(script.Package().Name())))
		codes := gen.buff.String()
		// 如果两个特定的内置包中的cb文件 则不需要加包前缀
		if script.Package().Name() == "base/cluster/network" {
			codes = strings.Replace(codes, "network.", "", -1)
		}
		if script.Package().Name() == "base/cluster" {
			codes = strings.Replace(codes, "cluster.", "", -1)
		}
		// service类型会用到logger
		if hasService {
			buff.WriteString(fmt.Sprintf("import log \"%s/base/logger\"\n", moduleName))
		}

		// 如果代码中有特定packageMapping中的包名 则引入对应的包
		for key, value := range packageMapping {
			if strings.Contains(codes, key) {
				buff.WriteString(value + "\n")
			}
		}
		for _, ref := range script.Imports {
			if _, ok := packageMapping[ref.Name()+"."]; ok {
				continue
			}
			// 如果代码中有对应的包名 则引入对应的包 并取别名为 包引用的名字
			if strings.Contains(codes, ref.Name()) {
				if ref.Name() == filepath.Base(ref.Ref.Name()) {
					buff.WriteString(fmt.Sprintf("import  \"%s/%s\"\n",
						moduleName, ref.Ref))
				} else {
					buff.WriteString(fmt.Sprintf("import %s \"%s/%s\"\n",
						ref.Name(), moduleName, ref.Ref))
				}
				// buff.WriteString(fmt.Sprintf("import %s \"%s/%s\"\n",
				// 	ref.Name(), moduleName, ref.Ref))

			}
		}
		// 将代码生成器的buff附加到此buff后
		buff.Write([]byte(codes))
		// 将buff写到文件
		gen.writeFile(script, buff.Bytes())
	}

	// 格式化代码
	formatScript(script)

	return script
}

// VisitEnum 访问枚举
func (gen *Gen4Go) VisitEnum(enum *ast.Enum) ast.Node {
	if cblang.IsError(enum) {
		if err := gen.tpl.ExecuteTemplate(&gen.buff, "error", enum); err != nil {
			cberrors.Panic(err.Error())
		}
	} else {
		if err := gen.tpl.ExecuteTemplate(&gen.buff, "enum", enum); err != nil {
			cberrors.Panic(err.Error())
		}
	}
	return enum
}

// VisitTable 访问表
func (gen *Gen4Go) VisitTable(table *ast.Table) ast.Node {
	table.Sort()
	if cblang.IsStruct(table) {
		if err := gen.tpl.ExecuteTemplate(&gen.buff, "struct", table); err != nil {
			cberrors.Panic(err.Error())
		}
	} else {
		if err := gen.tpl.ExecuteTemplate(&gen.buff, "table", table); err != nil {
			cberrors.Panic(err.Error())
		}
	}
	return table
}

// VisitService 访问协议
func (gen *Gen4Go) VisitService(service *ast.Service) ast.Node {
	if err := gen.tpl.ExecuteTemplate(&gen.buff, "service", service); err != nil {
		cberrors.Panic(err.Error())
	}
	return service
}
