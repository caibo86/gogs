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
	"gogs/base/gserrors"
	"gogs/base/gslang"
	"gogs/base/gslang/ast"
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
	"gsnet.":    `import "gogs/base/gsnet"`,
	"gserrors.": `import "gogs/base/gserrors"`,
	// "gss.":      `import "gogs/gss"`,
	"bytes.": `import "bytes"`,
	"fmt.":   `import "fmt"`,
	"time.":  `import "time"`,
	"bits.":  `import "math/bits"`,
	// "yfdocker.": `import "gogs/base/docker"`,
	// "yfconfig.": `import "gogs/base/config"`,
}

// gslang内置类型对应的golang表示
var keyMapping = map[string]string{
	".gslang.Byte":    "byte",
	".gslang.Sbyte":   "int8",
	".gslang.Int16":   "int16",
	".gslang.Uint16":  "uint16",
	".gslang.Int32":   "int32",
	".gslang.Uint32":  "uint32",
	".gslang.Int64":   "int64",
	".gslang.Uint64":  "uint64",
	".gslang.Float32": "float32",
	".gslang.Float64": "float64",
	".gslang.Bool":    "bool",
	".gslang.String":  "string",
	"Byte":            "byte",
	"Sbyte":           "int8",
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

// gslang内置类型的默认值对应的golang表示
var defaultVal = map[string]string{
	".gslang.Byte":    "byte(0)",
	".gslang.Sbyte":   "int8(0)",
	".gslang.Int16":   "int16(0)",
	".gslang.Uint16":  "uint16(0)",
	".gslang.Int32":   "int32(0)",
	".gslang.Uint32":  "uint32(0)",
	".gslang.Int64":   "int64(0)",
	".gslang.Uint64":  "uint64(0)",
	".gslang.Float32": "float32(0)",
	".gslang.Float64": "float64(0)",
	".gslang.Bool":    "false",
	".gslang.String":  "\"\"",
}

// gslang内置类型的零值对应的golang表示
var zeroVal = map[string]string{
	".gslang.Byte":    "0",
	".gslang.Sbyte":   "0",
	".gslang.Int16":   "0",
	".gslang.Uint16":  "0",
	".gslang.Int32":   "0",
	".gslang.Uint32":  "0",
	".gslang.Int64":   "0",
	".gslang.Uint64":  "0",
	".gslang.Float32": "0",
	".gslang.Float64": "0",
	".gslang.Bool":    "false",
	".gslang.String":  "\"\"",
}

// writeMapping 写入方法映射
var writeMapping = map[string]string{
	"Bool":    "gsnet.WriteBool",
	"Byte":    "gsnet.WriteByte",
	"Sbyte":   "gsnet.WriteSbyte",
	"Int16":   "gsnet.WriteInt16",
	"Uint16":  "gsnet.WriteUint16",
	"Int32":   "gsnet.WriteInt32",
	"Uint32":  "gsnet.WriteUint32",
	"Float32": "gsnet.WriteFloat32",
	"Int64":   "gsnet.WriteInt64",
	"Uint64":  "gsnet.WriteUint64",
	"Float64": "gsnet.WriteFloat64",
	"String":  "gsnet.WriteString",
}

// readMapping 读方法映射
var readMapping = map[string]string{
	"Bool":    "gsnet.ReadBool",
	"Byte":    "gsnet.ReadByte",
	"Sbyte":   "gsnet.ReadSbyte",
	"Int16":   "gsnet.ReadInt16",
	"Uint16":  "gsnet.ReadUint16",
	"Int32":   "gsnet.ReadInt32",
	"Uint32":  "gsnet.ReadUint32",
	"Float32": "gsnet.ReadFloat32",
	"Int64":   "gsnet.ReadInt64",
	"Uint64":  "gsnet.ReadUint64",
	"Float64": "gsnet.ReadFloat64",
	"String":  "gsnet.ReadString",
}

// compareMapping 比较方法映射
var compareMapping = map[string]string{
	"Bool":    `m.%s`,
	"Byte":    `m.%s != 0`,
	"Sbyte":   `m.%s != 0`,
	"Int16":   `m.%s != 0`,
	"Uint16":  `m.%s != 0`,
	"Int32":   `m.%s != 0`,
	"Uint32":  `m.%s != 0`,
	"Float32": `m.%s != 0`,
	"Int64":   `m.%s != 0`,
	"Uint64":  `m.%s != 0`,
	"Float64": `m.%s != 0`,
	"String":  `len(m.%s) > 0`,
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
	funcs := template.FuncMap{
		"symbol":        strings.Title,
		"typeName":      gen.typeName,
		"params":        gen.params,
		"returnParams":  gen.returnParams,
		"returnErr":     gen.returnErr,
		"callArgs":      gen.callArgs,
		"returnArgs":    gen.returnArgs,
		"readType":      gen.readType,
		"writeType":     gen.writeType,
		"defaultVal":    gen.defaultVal,
		"pos":           gslang.Pos,
		"lowerFirst":    gen.lowerFirst,
		"sovFunc":       gen.sovFunc,
		"calTypeSize":   gen.calTypeSize,
		"copyType":      gen.copyType,
		"printComments": gen.printComments,
	}
	gen.tpl, err = template.New("golang").Funcs(funcs).Parse(tpl4go)
	return
}

// lowerFirst 首字母小写
func (gen *Gen4Go) lowerFirst(name string) string {
	if len(name) == 0 {
		return name
	}
	return strings.ToLower(name[:1]) + name[1:]
}

// calTypeSize 根据类型获取计算大小的函数
func (gen *Gen4Go) calTypeSize(field *ast.Field) string {
	switch field.Type.(type) {
	case *ast.TypeRef:
		ref := field.Type.(*ast.TypeRef)
		switch ref.Ref.Name() {
		case "Bool":
			return fmt.Sprintf(
				`if m.%s {
					n += 3
				}`,
				field.Name())
		case "Byte", "Sbyte":
			return fmt.Sprintf(
				`if m.%s != 0 {
					n += 3
				}`,
				field.Name())
		case "Int16", "Uint16":
			return fmt.Sprintf(
				`if m.%s != 0 {
					n += 4
				}`,
				field.Name())
		case "Int32", "Uint32", "Float32":
			return fmt.Sprintf(
				`if m.%s != 0 {
					n += 6
				}`,
				field.Name())
		case "Int64", "Uint64", "Float64":
			return fmt.Sprintf(
				`if m.%s != 0 {
					n += 10
				}`,
				field.Name())
		case "String":
			return fmt.Sprintf(
				`l = len(m.%s)
					if l > 0 {
					n += 6 + l 
				}`,
				field.Name())
		default:
			switch ref.Ref.(type) {
			case *ast.Enum:
				return fmt.Sprintf(
					`if m.%s != 0 {
						n += 6
					}`,
					field.Name())
			case *ast.Table:
				return fmt.Sprintf(
					`if m.%s != nil {
						l = m.%s.Size()
						n += 6 + l
					}`,
					field.Name(), field.Name())
			default:
				panic(gserrors.Newf(nil, "not here %s", field.Type.Name()))
			}
		}
	case *ast.List, *ast.Array:
		// 数组 切片
		var ref *ast.TypeRef
		if list, ok := field.Type.(*ast.List); ok {
			ref = list.Element.(*ast.TypeRef)
		} else {
			ref = field.Type.(*ast.Array).Element.(*ast.TypeRef)
		}
		switch ref.Ref.Name() {
		case "Byte", "Sbyte", "Bool":
			return fmt.Sprintf(
				`l = len(m.%s)
				if l > 0 {
					n += 6 + l
				}`,
				field.Name())
		case "Uint16", "Int16":
			return fmt.Sprintf(
				`l = len(m.%s)
				if l > 0 {
					n += 6 + l * 2
				}`,
				field.Name())
		case "Uint32", "Int32", "Float32":
			return fmt.Sprintf(
				`l = len(m.%s)
				if l > 0 {
					n += 6 + l * 4
				}`,
				field.Name())
		case "Uint64", "Int64", "Float64":
			return fmt.Sprintf(
				`l = len(m.%s)
				if l > 0 {
					n += 6 + l * 8
				}`,
				field.Name())
		case "String":
			return fmt.Sprintf(
				`if len(m.%s) > 0 {
					n += 6
					for _, s := range m.%s {
						l = len(s)
						n += 4 + l
					}
				}`,
				field.Name(), field.Name())
		default:
			switch ref.Ref.(type) {
			case *ast.Enum:
				return fmt.Sprintf(
					`l = len(m.%s)
					if l > 0 {
						n += 6 + l * 4
					}`,
					field.Name())
			case *ast.Table:
				return fmt.Sprintf(
					`l = len(m.%s)
					if l > 0 {
						n += 6
						for _, e := range m.%s {
							n += 4 + e.Size()
						}
					}`,
					field.Name(), field.Name())
			default:
				panic(gserrors.Newf(nil, "not here %s", field.Type.Name()))
			}
		}
	case *ast.Map:
		// 字典
		hash := field.Type.(*ast.Map)
		var keyStr string
		var valStr string
		ref := hash.Key.(*ast.TypeRef)
		switch ref.Ref.Name() {
		case "Byte", "Sbyte", "Bool":
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
				panic(gserrors.Newf(nil, "map key can only be int or string, %s not supported", hash.Key.Name()))
			}
		}
		ref = hash.Value.(*ast.TypeRef)
		switch ref.Ref.Name() {
		case "Byte", "Sbyte", "Bool":
			valStr = "1"
		case "Uint16", "Int16":
			valStr = "2"
		case "Uint32", "Int32", "float32":
			valStr = "4"
		case "Uint64", "Int64", "float64":
			valStr = "8"
		case "String":
			valStr = "4 + len(v)"
		default:
			switch ref.Ref.(type) {
			case *ast.Enum:
				valStr = "4"
			case *ast.Table:
				valStr = "4 + v.Size()"
			default:
				panic(gserrors.Newf(nil, "map value %s not supported", hash.Value.Name()))
			}
		}
		return fmt.Sprintf(
			`if len(m.%s) > 0 {
						n += 6
						for k, v := range m.%s {
							_ = k
							_ = v
							n += %s
							n += %s
						}
					}`,
			field.Name(),
			field.Name(),
			keyStr,
			valStr)
	}
	panic(gserrors.Newf(nil, "not here"))
	return "unknown"
}

// writeType 根据字段类型生成写入函数
func (gen *Gen4Go) writeType(field *ast.Field) string {
	switch field.Type.(type) {
	case *ast.TypeRef:
		var str string
		var ok bool
		ref := field.Type.(*ast.TypeRef)
		if str, ok = writeMapping[ref.Ref.Name()]; ok {
			return fmt.Sprintf(
				`if `+compareMapping[ref.Ref.Name()]+` {
					i = gsnet.WriteFieldID(data, i, %d)
					i = %s(data, i, m.%s)
				}`,
				field.Name(), field.ID, str, field.Name())
		} else {
			switch ref.Ref.(type) {
			case *ast.Enum:
				return fmt.Sprintf(
					`if m.%s != 0 {
						i = gsnet.WriteFieldID(data, i, %d)
						i = gsnet.WriteEnum(data, i, int32(m.%s))
					}`,
					field.Name(), field.ID, field.Name())
			case *ast.Table:
				return fmt.Sprintf(
					`if m.%s != nil {
						i = gsnet.WriteFieldID(data, i, %d)
						size := m.%s.Size()
						i = gsnet.WriteUint32(data, i, uint32(size))
						m.%s.MarshalToSizedBuffer(data[i:])
						i += size
					}`,
					field.Name(), field.ID, field.Name(), field.Name())
			default:
				panic(gserrors.Newf(nil, "not here %s", field.Type.Name()))
			}
		}
	case *ast.List, *ast.Array:
		// 数组
		var ref *ast.TypeRef
		isList := false
		if list, ok := field.Type.(*ast.List); ok {
			ref = list.Element.(*ast.TypeRef)
			isList = true
		} else {
			ref = field.Type.(*ast.Array).Element.(*ast.TypeRef)
		}
		var str string
		var ok bool
		if str, ok = writeMapping[ref.Ref.Name()]; ok {
			str = fmt.Sprintf("i = %s(data, i, e)", str)
		} else {
			switch ref.Ref.(type) {
			case *ast.Enum:
				str = `i = gsnet.WriteEnum(data, i, int32(e))`
			case *ast.Table:
				str = `size := e.Size()
						i = gsnet.WriteUint32(data, i, uint32(size))
						if e != nil {
							e.MarshalToSizedBuffer(data[i:])
						}
						i += size`
			default:
				panic(gserrors.Newf(nil, "not here %s", field.Type.Name()))
			}
		}
		if isList {
			return fmt.Sprintf(
				`if len(m.%s) > 0 {
					i = gsnet.WriteFieldID(data, i , %d)
					i = gsnet.WriteUint32(data, i, uint32(len(m.%s)))
					for _, e := range m.%s {
						%s
					}
				}`,
				field.Name(),
				field.ID,
				field.Name(),
				field.Name(),
				str)
		} else {
			return fmt.Sprintf(
				`i = gsnet.WriteFieldID(data, i , %d)
					i = gsnet.WriteUint32(data, i, uint32(len(m.%s)))
					for _, e := range m.%s {
						%s
					}`,
				field.ID,
				field.Name(),
				field.Name(),
				str)
		}

	case *ast.Map:
		// 字典
		hash := field.Type.(*ast.Map)
		var keyStr string
		var valStr string
		ref := hash.Key.(*ast.TypeRef)
		var ok bool
		if keyStr, ok = writeMapping[ref.Ref.Name()]; ok {
			keyStr = fmt.Sprintf("i = %s(data, i, k)", keyStr)
		} else {
			switch ref.Ref.(type) {
			case *ast.Enum:
				keyStr = `i = gsnet.WriteEnum(data, i, int32(k))`
			default:
				panic(gserrors.Newf(nil, "map key can only be int or string, %s not supported", hash.Key.Name()))
			}
		}
		ref = hash.Value.(*ast.TypeRef)
		if valStr, ok = writeMapping[ref.Ref.Name()]; ok {
			valStr = fmt.Sprintf("i = %s(data, i, v)", valStr)
		} else {
			switch ref.Ref.(type) {
			case *ast.Enum:
				valStr = `i = gsnet.WriteEnum(data, i, int32(v))`
			case *ast.Table:
				valStr = `size := v.Size()
						i = gsnet.WriteUint32(data, i, uint32(size))
						if v != nil {
							v.MarshalToSizedBuffer(data[i:])
						}
						i += size`
			default:
				panic(gserrors.Newf(nil, "map value %s not supported", hash.Value.Name()))
			}
		}
		return fmt.Sprintf(
			`if len(m.%s) > 0 {
						i = gsnet.WriteFieldID(data, i, %d)
						i = gsnet.WriteUint32(data, i, uint32(len(m.%s)))
						for k, v := range m.%s {
							%s
							%s
						}
					}`,
			field.Name(),
			field.ID,
			field.Name(),
			field.Name(),
			keyStr,
			valStr)
	}
	panic(gserrors.Newf(nil, "not here"))
	return "unknown"
}

// readType 根据字段类型生成读取函数
func (gen *Gen4Go) readType(field *ast.Field) string {
	switch field.Type.(type) {
	case *ast.TypeRef:
		ref := field.Type.(*ast.TypeRef)
		var ok bool
		var str string
		if str, ok = readMapping[ref.Ref.Name()]; ok {
			return fmt.Sprintf("i, m.%s = %s(data, i)", field.Name(), str)
		} else {
			switch ref.Ref.(type) {
			case *ast.Enum:
				return fmt.Sprintf(`var v int32
						i, v = gsnet.ReadEnum(data, i)
						m.%s = %s(v)`,
					field.Name(), gen.typeName(ref))
			case *ast.Table:
				return fmt.Sprintf(`var size uint32
						i, size = gsnet.ReadUint32(data, i)
						if m.%s == nil {
							m.%s = %s
						}
						if err = m.%s.Unmarshal(data[i:i+int(size)]); err != nil {
							return
						} 
						i += int(size)`,
					field.Name(), field.Name(), gen.defaultVal(ref), field.Name())
			default:
				panic(gserrors.Newf(nil, "not here %s", field.Type.Name()))
			}
		}
	case *ast.List, *ast.Array:
		// 数组
		var ref *ast.TypeRef
		var isList bool
		if list, ok := field.Type.(*ast.List); ok {
			isList = true
			ref = list.Element.(*ast.TypeRef)
		} else {
			ref = field.Type.(*ast.Array).Element.(*ast.TypeRef)
		}
		var str string
		var ok bool
		if str, ok = readMapping[ref.Ref.Name()]; ok {
			str = fmt.Sprintf("i, m.%s[j] = %s(data, i)", field.Name(), str)
		} else {
			switch ref.Ref.(type) {
			case *ast.Enum:
				str = fmt.Sprintf(`var v int32
				i, v = gsnet.ReadEnum(data, i)
				m.%s[j] = %s(v)`,
					field.Name(), gen.typeName(ref))
			case *ast.Table:
				str = fmt.Sprintf(
					`var size uint32
						i, size = gsnet.ReadUint32(data, i)
						if size > 0 {
							m.%s[j] = %s
							if err = m.%s[j].Unmarshal(data[i:i+int(size)]); err != nil {
								return
							}
						}
						i += int(size)`, field.Name(), gen.defaultVal(ref), field.Name())
			default:
				panic(gserrors.Newf(nil, "not here %s", field.Type.Name()))
			}
		}
		if isList {
			return fmt.Sprintf(
				`var length uint32
				i, length = gsnet.ReadUint32(data, i)
				m.%s = make([]%s, length)
				for j := uint32(0); j < length; j++ {
					%s
				}`,
				field.Name(), gen.typeName(ref), str)
		} else {
			return fmt.Sprintf(
				`var length uint32
				i, length = gsnet.ReadUint32(data, i)
				for j := uint32(0); j < length; j++ {
					%s
				}`,
				str)
		}
	case *ast.Map:
		// 字典
		hash := field.Type.(*ast.Map)
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
					i, k1 = gsnet.ReadEnum(data, i)
					k := %s(k1)`, gen.typeName(hash.Key))
			default:
				panic(gserrors.Newf(nil, "map key can only be int or string, %s not supported", hash.Key.Name()))
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
					i, v1 = gsnet.ReadEnum(data, i)
					v := %s(v1)`, gen.typeName(hash.Value))
			case *ast.Table:
				valStr = fmt.Sprintf(`var size uint32
						var v %s
						i, size = gsnet.ReadUint32(data, i)
						if size > 0 {
							v = %s
							if err = v.Unmarshal(data[i:i+int(size)]); err != nil {
								return
							}
						}
						i += int(size)`, gen.typeName(hash.Value), gen.defaultVal(hash.Value))
			default:
				panic(gserrors.Newf(nil, "map value %s not supported", hash.Value.Name()))
			}
		}
		return fmt.Sprintf(
			`var length uint32
					i, length = gsnet.ReadUint32(data, i)
					if m.%s == nil{
						m.%s = make(map[%s]%s)
					}
					for j := uint32(0); j < length; j++ {
						%s
						%s
						m.%s[k] = v
					}`,
			field.Name(), field.Name(),
			gen.typeName(hash.Key),
			gen.typeName(hash.Value),
			keyStr, valStr, field.Name())
	}
	panic(gserrors.Newf(nil, "not here"))
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
	case *ast.List:
		list := field.Type.(*ast.List)
		ref := list.Element.(*ast.TypeRef)
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
			panic(err)
		}
		return buff.String()
	case *ast.List:
		return "nil"
	case *ast.Map:
		// 字典
		return fmt.Sprintf("make(%s)", gen.typeName(expr))
	}
	panic(gserrors.Newf(nil, "not here"))
	return "unknown"
}

// zeroVal 根据类型取其零值
func (gen *Gen4Go) zeroVal(expr ast.Expr) string {
	switch expr.(type) {
	case *ast.TypeRef:
		// 内置类型
		if val, ok := zeroVal[expr.Name()]; ok {
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
			panic(err)
		}
		return buff.String()
	case *ast.List:
		return "nil"
	case *ast.Map:
		// 字典
		return fmt.Sprintf("make(%s)", gen.typeName(expr))
	}
	panic(gserrors.Newf(nil, "not here"))
	return "unknown"
}

// params 根据参数生成函数声明的入参列表
func (gen *Gen4Go) params(params []*ast.Param) string {
	if len(params) == 0 {
		return "()"
	}
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("(arg0 %s", gen.typeName(params[0].Type)))
	for i := 1; i < len(params); i++ {
		buff.WriteString(fmt.Sprintf(",arg%d %s", i, gen.typeName(params[i].Type)))
	}
	buff.WriteString(")")
	return buff.String()
}

// returnParams 根据参数生成函数声明的返回参数列表
func (gen *Gen4Go) returnParams(params []*ast.Param) string {
	if len(params) == 0 {
		return "(err error)"
	}
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("(ret0 %s", gen.typeName(params[0].Type)))
	for i := 1; i < len(params); i++ {
		buff.WriteString(fmt.Sprintf(",ret%d %s", i, gen.typeName(params[i].Type)))
	}
	buff.WriteString(",err error)")
	return buff.String()
}

// callArgs 根据参数生成函数的调用参数列表
func (gen *Gen4Go) callArgs(params []*ast.Param) string {
	if len(params) == 0 {
		return "()"
	}
	var buff bytes.Buffer
	buff.WriteString("(arg0")
	for i := 1; i < len(params); i++ {
		buff.WriteString(fmt.Sprintf(",arg%d", i))
	}
	buff.WriteString(")")
	return buff.String()
}

// returnArgs 根据函数生成 接收 函数的调用值 的 声明列表
func (gen *Gen4Go) returnArgs(method *ast.Method) string {
	params := method.Return
	if len(params) == 0 {
		return "err = "
	}
	var buff bytes.Buffer
	buff.WriteString("ret0")
	for i := 1; i < len(params); i++ {
		buff.WriteString(fmt.Sprintf(", ret%d", i))
	}
	buff.WriteString(", err = ")
	return buff.String()
}

// returnErr .
func (gen *Gen4Go) returnErr(params []*ast.Param) string {
	if len(params) == 0 {
		return "err"
	}
	var buff bytes.Buffer
	buff.WriteString(gen.defaultVal(params[0].Type))
	for i := 1; i < len(params); i++ {
		buff.WriteString(fmt.Sprintf(",%s", gen.defaultVal(params[0].Type)))
	}
	buff.WriteString(",err)")
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
	case *ast.List:
		// 切片
		list := expr.(*ast.List)
		return fmt.Sprintf("[]%s", gen.typeName(list.Element))
	case *ast.Map:
		// 字典
		hash := expr.(*ast.Map)
		return fmt.Sprintf("map[%s]%s", gen.typeName(hash.Key), gen.typeName(hash.Value))
	}
	panic(gserrors.Newf(nil, "unknown golang typeName: %s\n\t%s",
		expr,
		gslang.Pos(expr)))
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
	comments := gslang.Comments(node)
	if len(comments) > 0 {
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

// writeFile 将代码节点对应的golang代码写入到文件
func (gen *Gen4Go) writeFile(script *ast.Script, bytes []byte) {
	fullPath, ok := gslang.FilePath(script)
	if !ok {
		panic(gserrors.Newf(nil, "compile must bind file path to script"))
	}
	// 写入文件名为 源文件名+.go
	fullPath += ".go"
	err := os.WriteFile(fullPath, bytes, 0644)
	if err != nil {
		panic(err)
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
	// 内置gslang包则直接返回
	if pkg.Name() == "gogs/base/gslang" {
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
		panic(err)
	}

	// 轮询访问代码中的所有类型 Enum Struct Table Contract
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
	for _, t := range script.Types {
		if _, ok := t.(*ast.Contract); ok {
			t.Accept(gen)
		}
	}
	// 代码中有类型
	if gen.buff.Len() > 0 {
		var buff bytes.Buffer
		filename, _ := gslang.FilePath(script)
		filename += ".go"
		filename = filepath.Base(filename)
		// 写入额外信息
		buff.WriteString(fmt.Sprintf(
			`// -------------------------------------------
// @file      : %s
// @author    : generated by gsc, do not edit
// @contact   : caibo923@gmail.com
// @time      : %s
// -------------------------------------------

`, filename, time.Now().Format(time.RFC3339)))

		// 写入包声明
		buff.WriteString(fmt.Sprintf("package %s\n", filepath.Base(script.Package().Name())))
		codes := gen.buff.String()
		// 如果两个特定的内置包中的gs文件 则不需要加包前缀
		if script.Package().Name() == "gogs/base/gsnet" {
			codes = strings.Replace(codes, "gsnet.", "", -1)
		}
		if script.Package().Name() == "gogs/base/gsdocker" {
			codes = strings.Replace(codes, "gsdocker.", "", -1)
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
				// buff.WriteString(fmt.Sprintf("import %s \"%s/%s\"\n",
				// 	ref.Name(), moduleName, ref.Ref))
				buff.WriteString(fmt.Sprintf("import \"%s/%s\"\n",
					moduleName, ref.Ref))
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
	if gslang.IsError(enum) {
		if err := gen.tpl.ExecuteTemplate(&gen.buff, "error", enum); err != nil {
			panic(err)
		}
	} else {
		if err := gen.tpl.ExecuteTemplate(&gen.buff, "enum", enum); err != nil {
			panic(err)
		}
	}
	return enum
}

// VisitTable 访问表
func (gen *Gen4Go) VisitTable(table *ast.Table) ast.Node {
	table.Sort()

	if gslang.IsStruct(table) {
		if err := gen.tpl.ExecuteTemplate(&gen.buff, "struct", table); err != nil {
			panic(err)
		}
	} else {
		if err := gen.tpl.ExecuteTemplate(&gen.buff, "table", table); err != nil {
			panic(err)
		}
	}
	return table
}

// VisitContract 访问协议
func (gen *Gen4Go) VisitContract(contract *ast.Contract) ast.Node {
	log.Debugf("%v", contract.Name())
	log.Debugf("%v", contract.Path())
	log.Debugf("%v", contract.Methods)
	log.Debugf("%v", contract.Bases)
	for _, method := range contract.Methods {
		log.Debugf("%v", method.Return)
		log.Debugf("%v", method.Params)
	}
	if err := gen.tpl.ExecuteTemplate(&gen.buff, "contract", contract); err != nil {
		panic(err)
	}
	return contract
}
