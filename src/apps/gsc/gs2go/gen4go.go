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
	".gslang.Byte":    "gsnet.WriteByte",
	".gslang.Sbyte":   "gsnet.WriteSbyte",
	".gslang.Int16":   "gsnet.WriteInt16",
	".gslang.Uint16":  "gsnet.WriteUint16",
	".gslang.Int32":   "gsnet.WriteInt32",
	".gslang.Uint32":  "gsnet.WriteUint32",
	".gslang.Int64":   "gsnet.WriteInt64",
	".gslang.Uint64":  "gsnet.WriteUint64",
	".gslang.Float32": "gsnet.WriteFloat32",
	".gslang.Float64": "gsnet.WriteFloat64",
	".gslang.Bool":    "gsnet.WriteBool",
	".gslang.String":  "gsnet.WriteString",
}

// writeTagMapping 带标签的写入方法映射
var writeTagMapping = map[string]string{
	".gslang.Byte":    "gsnet.WriteTagByte",
	".gslang.Sbyte":   "gsnet.WriteTagSbyte",
	".gslang.Int16":   "gsnet.WriteTagInt16",
	".gslang.Uint16":  "gsnet.WriteTagUint16",
	".gslang.Int32":   "gsnet.WriteTagInt32",
	".gslang.Uint32":  "gsnet.WriteTagUint32",
	".gslang.Int64":   "gsnet.WriteTagInt64",
	".gslang.Uint64":  "gsnet.WriteTagUint64",
	".gslang.Float32": "gsnet.WriteTagFloat32",
	".gslang.Float64": "gsnet.WriteTagFloat64",
	".gslang.Bool":    "gsnet.WriteTagBool",
	".gslang.String":  "gsnet.WriteTagString",
}

// readMapping 读方法映射
var readMapping = map[string]string{
	".gslang.Byte":    "gsnet.ReadByte",
	".gslang.Sbyte":   "gsnet.ReadSbyte",
	".gslang.Int16":   "gsnet.ReadInt16",
	".gslang.Uint16":  "gsnet.ReadUint16",
	".gslang.Int32":   "gsnet.ReadInt32",
	".gslang.Uint32":  "gsnet.ReadUint32",
	".gslang.Int64":   "gsnet.ReadInt64",
	".gslang.Uint64":  "gsnet.ReadUint64",
	".gslang.Float32": "gsnet.ReadFloat32",
	".gslang.Float64": "gsnet.ReadFloat64",
	".gslang.Bool":    "gsnet.ReadBool",
	".gslang.String":  "gsnet.ReadString",
}

// tagMapping 标签映射
var tagMapping = map[string]string{
	".gslang.Byte":    "gsnet.Byte",
	".gslang.Sbyte":   "gsnet.Sbyte",
	".gslang.Int16":   "gsnet.Int16",
	".gslang.Uint16":  "gsnet.Uint16",
	".gslang.Int32":   "gsnet.Int32",
	".gslang.Uint32":  "gsnet.Uint32",
	".gslang.Int64":   "gsnet.Int64",
	".gslang.Uint64":  "gsnet.Uint64",
	".gslang.Float32": "gsnet.Float32",
	".gslang.Float64": "gsnet.Float64",
	".gslang.Bool":    "gsnet.Bool",
	".gslang.String":  "gsnet.String",
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
		"enumType":     gen.enumType,
		"symbol":       strings.Title,
		"typeName":     gen.typeName,
		"params":       gen.params,
		"returnParams": gen.returnParams,
		"returnErr":    gen.returnErr,
		"callargs":     gen.callargs,
		"returnargs":   gen.returnargs,
		"readType":     gen.readType,
		"writeType":    gen.writeType,
		"defaultVal":   gen.defaultVal,
		"pos":          gslang.Pos,
		"lowerFirst":   gen.lowerFirst,
		"sovFunc":      gen.sovFunc,
		"calTypeSize":  gen.calTypeSize,
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
				log.Panicf("not here %s", field.Type.Name())
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
					`if len(m.%s) > 0 {
						n += 6
						for _, e := range m.%s {
							n += 4 + e.Size()
						}
					}`,
					field.Name(), field.Name())
			default:
				log.Panicf("not here %s", field.Type.Name())
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
				log.Panicf("map key can only be int or string, %s not supported", hash.Key.Name())
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
				log.Panicf("map value %s not supported", hash.Value.Name())
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
	log.Panic("not here")
	return "unknown"
}

// writeType 根据字段类型生成写入函数
func (gen *Gen4Go) writeType(field *ast.Field) string {
	switch field.Type.(type) {
	case *ast.TypeRef:
		ref := field.Type.(*ast.TypeRef)
		switch ref.Ref.Name() {
		case "Bool":
			return fmt.Sprintf(
				`if m.%s {
					i = gsnet.WriteFieldID(data, i, %d)
					i = gsnet.WriteBool(data, i, m.%s)
				}`,
				field.Name(), field.ID, field.Name())
		case "Byte":
			return fmt.Sprintf(
				`if m.%s != 0 {
					i = gsnet.WriteFieldID(data, i, %d)
					i = gsnet.WriteByte(data, i, m.%s)
				}`,
				field.Name(), field.ID, field.Name())
		case "Sbyte":
			return fmt.Sprintf(
				`if m.%s != 0 {
					i = gsnet.WriteFieldID(data, i, %d)
					i = gsnet.WriteSbyte(data, i, m.%s)
				}`,
				field.Name(), field.ID, field.Name())
		case "Int16":
			return fmt.Sprintf(
				`if m.%s != 0 {
					i = gsnet.WriteFieldID(data, i, %d)
					i = gsnet.WriteInt16(data, i, m.%s)
				}`,
				field.Name(), field.ID, field.Name())
		case "Uint16":
			return fmt.Sprintf(
				`if m.%s != 0 {
					i = gsnet.WriteFieldID(data, i, %d)
					i = gsnet.WriteUint16(data, i, m.%s)
				}`,
				field.Name(), field.ID, field.Name())
		case "Int32":
			return fmt.Sprintf(
				`if m.%s != 0 {
					i = gsnet.WriteFieldID(data, i, %d)
					i = gsnet.WriteInt32(data, i, m.%s)
				}`,
				field.Name(), field.ID, field.Name())
		case "Uint32":
			return fmt.Sprintf(
				`if m.%s != 0 {
					i = gsnet.WriteFieldID(data, i, %d)
					i = gsnet.WriteUint32(data, i, m.%s)
				}`,
				field.Name(), field.ID, field.Name())
		case "Float32":
			return fmt.Sprintf(
				`if m.%s != 0 {
					i = gsnet.WriteFieldID(data, i, %d)
					i = gsnet.WriteFloat32(data, i, m.%s)
				}`,
				field.Name(), field.ID, field.Name())
		case "Int64":
			return fmt.Sprintf(
				`if m.%s != 0 {
					i = gsnet.WriteFieldID(data, i, %d)
					i = gsnet.WriteInt64(data, i, m.%s)
				}`,
				field.Name(), field.ID, field.Name())
		case "Uint64":
			return fmt.Sprintf(
				`if m.%s != 0 {
					i = gsnet.WriteFieldID(data, i, %d)
					i = gsnet.WriteUint64(data, i, m.%s)
				}`,
				field.Name(), field.ID, field.Name())
		case "Float64":
			return fmt.Sprintf(
				`if m.%s != 0 {
					i = gsnet.WriteFieldID(data, i, %d)
					i = gsnet.WriteFloat64(data, i, m.%s)
				}`,
				field.Name(), field.ID, field.Name())
		case "String":
			return fmt.Sprintf(
				`if len(m.%s) > 0 {
					i = gsnet.WriteFieldID(data, i, %d)
					i = gsnet.WriteString(data, i, m.%s)
				}`,
				field.Name(), field.ID, field.Name())
		default:
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
				log.Panicf("not here %s", field.Type.Name())
			}
		}
	case *ast.List, *ast.Array:
		// 数组
		var ref *ast.TypeRef
		if list, ok := field.Type.(*ast.List); ok {
			ref = list.Element.(*ast.TypeRef)
		} else {
			ref = field.Type.(*ast.Array).Element.(*ast.TypeRef)
		}
		var str string
		switch ref.Ref.Name() {
		case "Bool":
			str = `i = gsnet.WriteBool(data, i, e)`
		case "Byte":
			str = `i = gsnet.WriteByte(data, i, e)`
		case "Sbyte":
			str = `i = gsnet.WriteSbyte(data, i, e)`
		case "Int16":
			str = `i = gsnet.WriteInt16(data, i, e)`
		case "Uint16":
			str = `i = gsnet.WriteUint16(data, i, e)`
		case "Int32":
			str = `i = gsnet.WriteInt32(data, i, e)`
		case "Uint32":
			str = `i = gsnet.WriteUint32(data, i, e)`
		case "Float32":
			str = `i = gsnet.WriteFloat32(data, i, e)`
		case "Int64":
			str = `i = gsnet.WriteInt64(data, i, e)`
		case "Uint64":
			str = `i = gsnet.WriteUint64(data, i, e)`
		case "Float64":
			str = `i = gsnet.WriteFloat64(data, i, e)`
		case "String":
			str = `i = gsnet.WriteString(data, i, e)`
		default:
			switch ref.Ref.(type) {
			case *ast.Enum:
				str = `i = gsnet.WriteEnum(data, i, int32(e))`
			case *ast.Table:
				str = `size := e.Size()
						i = gsnet.WriteUint32(data, i, uint32(size))
						e.MarshalToSizedBuffer(data[i:])
						i += size`
			default:
				log.Panicf("not here %s", field.Type.Name())
			}
		}
		return fmt.Sprintf(
			`if len(m.%s) > 0 {
					i = gsnet.WriteFieldID(data, i , %d)
					i = gsnet.WriteUint32(data, i, uint32(len(m.%s)))
					for _, e := range m.%s {
						_ = e
						%s
					}
				}`,
			field.Name(),
			field.ID,
			field.Name(),
			field.Name(),
			str)
	case *ast.Map:
		// 字典
		hash := field.Type.(*ast.Map)
		var keyStr string
		var valStr string
		ref := hash.Key.(*ast.TypeRef)
		switch ref.Ref.Name() {
		case "Bool":
			keyStr = "i = gsnet.WriteBool(data, i, k)"
		case "Byte":
			keyStr = "i = gsnet.WriteByte(data, i, k)"
		case "Sbyte":
			keyStr = "i = gsnet.WriteSbyte(data, i, k)"
		case "Int16":
			keyStr = "i = gsnet.WriteInt16(data, i, k)"
		case "Uint16":
			keyStr = "i = gsnet.WriteUint16(data, i, k)"
		case "Int32":
			keyStr = "i = gsnet.WriteInt32(data, i, k)"
		case "Uint32":
			keyStr = "i = gsnet.WriteUint32(data, i, k)"
		case "Int64":
			keyStr = "i = gsnet.WriteInt64(data, i, k)"
		case "Uint64":
			keyStr = "i = gsnet.WriteInt64(data, i, k)"
		case "String":
			keyStr = "i = gsnet.WriteString(data, i, k)"
		default:
			switch ref.Ref.(type) {
			case *ast.Enum:
				keyStr = "i = gsnet.WriteEnum(data, i, int32(k))"
			default:
				log.Panicf("map key can only be int or string, %s not supported", hash.Key.Name())
			}
		}
		ref = hash.Value.(*ast.TypeRef)
		switch ref.Ref.Name() {
		case "Bool":
			valStr = "i = gsnet.WriteBool(data, i, v)"
		case "Byte":
			valStr = "i = gsnet.WriteByte(data, i, v)"
		case "Sbyte":
			valStr = "i = gsnet.WriteSbyte(data, i, v)"
		case "Int16":
			valStr = "i = gsnet.WriteInt16(data, i, v)"
		case "Uint16":
			valStr = "i = gsnet.WriteUint16(data, i, v)"
		case "Int32":
			valStr = "i = gsnet.WriteInt32(data, i, v)"
		case "Uint32":
			valStr = "i = gsnet.WriteUint32(data, i, v)"
		case "Float32":
			valStr = "i = gsnet.WriteFloat32(data, i, v)"
		case "Int64":
			valStr = "i = gsnet.WriteInt64(data, i, v)"
		case "Uint64":
			valStr = "i = gsnet.WriteInt64(data, i, v)"
		case "Float64":
			valStr = "i = gsnet.WriteFloat64(data, i, v)"
		case "String":
			valStr = "i = gsnet.WriteString(data, i, v)"
		default:
			switch ref.Ref.(type) {
			case *ast.Enum:
				valStr = "i = gsnet.WriteEnum(data, i, int32(v))"
			case *ast.Table:
				valStr = `size := v.Size()
						i = gsnet.WriteUint32(data, i, uint32(size))
						v.MarshalToSizedBuffer(data[i:])
						i += size`
			default:
				log.Panicf("map key can only be int or string, %s not supported", hash.Key.Name())
			}
		}

		return fmt.Sprintf(
			`if len(m.%s) > 0 {
						i = gsnet.WriteFieldID(data, i, %d)
						i = gsnet.WriteUint32(data, i, uint32(len(m.%s)))
						for k, v := range m.%s {
 							_ = k
							_ = v
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
	log.Panic("not here")
	return "unknown"
}

// readType 根据字段类型生成读取函数
func (gen *Gen4Go) readType(field *ast.Field) string {
	switch field.Type.(type) {
	case *ast.TypeRef:
		ref := field.Type.(*ast.TypeRef)
		switch ref.Ref.Name() {
		case "Bool":
			return fmt.Sprintf(
				`i, m.%s = gsnet.ReadBool(data, i)`,
				field.Name())
		case "Byte":
			return fmt.Sprintf(
				`i, m.%s = gsnet.ReadByte(data, i)`,
				field.Name())
		case "Sbyte":
			return fmt.Sprintf(
				`i, m.%s = gsnet.ReadSbyte(data, i)`,
				field.Name())
		case "Int16":
			return fmt.Sprintf(
				`i, m.%s = gsnet.ReadInt16(data, i)`,
				field.Name())
		case "Uint16":
			return fmt.Sprintf(
				`i, m.%s = gsnet.ReadUint16(data, i)`,
				field.Name())
		case "Int32":
			return fmt.Sprintf(
				`i, m.%s = gsnet.ReadInt32(data, i)`,
				field.Name())
		case "Uint32":
			return fmt.Sprintf(
				`i, m.%s = gsnet.ReadUint32(data, i)`,
				field.Name())
		case "Float32":
			return fmt.Sprintf(
				`i, m.%s = gsnet.ReadFloat32(data, i)`,
				field.Name())
		case "Int64":
			return fmt.Sprintf(
				`i, m.%s = gsnet.ReadInt64(data, i)`,
				field.Name())
		case "Uint64":
			return fmt.Sprintf(
				`i, m.%s = gsnet.ReadUint64(data, i)`,
				field.Name())
		case "Float64":
			return fmt.Sprintf(
				`i, m.%s = gsnet.ReadFloat64(data, i)`,
				field.Name())
		case "String":
			return fmt.Sprintf(
				`i, m.%s = gsnet.ReadString(data, i)`,
				field.Name())
		default:
			switch ref.Ref.(type) {
			case *ast.Enum:
				return fmt.Sprintf(
					`var v int32
						i, v = gsnet.ReadEnum(data, i)
						m.%s = %s(v)`,
					field.Name(), gen.typeName(ref))
			case *ast.Table:
				return fmt.Sprintf(
					`var size uint32
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
				log.Panicf("not here %s", field.Type.Name())
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
		switch ref.Ref.Name() {
		case "Bool":
			str = fmt.Sprintf(`i, m.%s[j] = gsnet.ReadBool(data, i)`, field.Name())
		case "Byte":
			str = fmt.Sprintf(`i, m.%s[j] = gsnet.ReadByte(data, i)`, field.Name())
		case "Sbyte":
			str = fmt.Sprintf(`i, m.%s[j] = gsnet.ReadSbyte(data, i)`, field.Name())
		case "Int16":
			str = fmt.Sprintf(`i, m.%s[j] = gsnet.ReadInt16(data, i)`, field.Name())
		case "Uint16":
			str = fmt.Sprintf(`i, m.%s[j] = gsnet.ReadUint16(data, i)`, field.Name())
		case "Int32":
			str = fmt.Sprintf(`i, m.%s[j] = gsnet.ReadInt32(data, i)`, field.Name())
		case "Uint32":
			str = fmt.Sprintf(`i, m.%s[j] = gsnet.ReadUint32(data, i)`, field.Name())
		case "Float32":
			str = fmt.Sprintf(`i, m.%s[j] = gsnet.ReadFloat32(data, i)`, field.Name())
		case "Int64":
			str = fmt.Sprintf(`i, m.%s[j] = gsnet.ReadInt64(data, i)`, field.Name())
		case "Uint64":
			str = fmt.Sprintf(`i, m.%s[j] = gsnet.ReadUint64(data, i)`, field.Name())
		case "Float64":
			str = fmt.Sprintf(`i, m.%s[j] = gsnet.ReadFloat64(data, i)`, field.Name())
		case "String":
			str = fmt.Sprintf(`i, m.%s[j] = gsnet.ReadString(data, i)`, field.Name())
		default:
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
						m.%s[j] = %s
						if err = m.%s[j].Unmarshal(data[i:i+int(size)]); err != nil {
							return
						}
						i += int(size)`, field.Name(), gen.defaultVal(ref), field.Name())
			default:
				log.Panicf("not here %s", field.Type.Name())
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
		ref := hash.Key.(*ast.TypeRef)
		switch ref.Ref.Name() {
		case "Bool":
			keyStr = `var k bool
					i, k = gsnet.ReadBool(data, i)`
		case "Byte":
			keyStr = `var k byte
					i, k = gsnet.ReadByte(data, i)`
		case "Sbyte":
			keyStr = `var k int8
					i, k = gsnet.ReadSbyte(data, i)`
		case "Int16":
			keyStr = `var k int16
					i, k = gsnet.ReadInt16(data, i)`
		case "Uint16":
			keyStr = `var k uint16
					i, k = gsnet.ReadUint16(data, i)`
		case "Int32":
			keyStr = `var k int32
					i, k = gsnet.ReadInt32(data, i)`
		case "Uint32":
			keyStr = `var k uint32
					i, k = gsnet.ReadUint32(data, i)`
		case "Int64":
			keyStr = `var k int64
					i, k = gsnet.ReadInt64(data, i)`
		case "Uint64":
			keyStr = `var k uint64
					i, k = gsnet.ReadUint64(data, i)`
		case "String":
			keyStr = `var k string
					i, k = gsnet.ReadString(data, i)`
		default:
			switch ref.Ref.(type) {
			case *ast.Enum:
				keyStr = fmt.Sprintf(`var k1 int32
					i, k1 = gsnet.ReadInt32(data, i)
					k := %s(k1)`, gen.typeName(hash.Key))
			default:
				log.Panicf("map key can only be int or string, %s not supported", hash.Key.Name())
			}
		}
		ref = hash.Value.(*ast.TypeRef)
		switch ref.Ref.Name() {
		case "Bool":
			valStr = `var v bool
					i, v = gsnet.ReadBool(data, i)`
		case "Byte":
			valStr = `var v byte
					i, v = gsnet.ReadByte(data, i)`
		case "Sbyte":
			valStr = `var v int8
					i, v = gsnet.ReadSbyte(data, i)`
		case "Int16":
			valStr = `var v int16
					i, v = gsnet.ReadInt16(data, i)`
		case "Uint16":
			valStr = `var v uint16
					i, v = gsnet.ReadUint16(data, i)`
		case "Int32":
			valStr = `var v int32
					i, v = gsnet.ReadInt32(data, i)`
		case "Uint32":
			valStr = `var v uint32
					i, v = gsnet.ReadUint32(data, i)`
		case "Float32":
			valStr = `var v float32
					i, v = gsnet.ReadFloat32(data, i)`
		case "Int64":
			valStr = `var v int64
					i, v = gsnet.ReadInt64(data, i)`
		case "Uint64":
			valStr = `var v uint64
					i, v = gsnet.ReadUint64(data, i)`
		case "Float64":
			valStr = `var v float64
					i, v = gsnet.ReadFloat64(data, i)`
		case "String":
			valStr = `var v string
					i, v = gsnet.ReadString(data, i)`
		default:
			switch ref.Ref.(type) {
			case *ast.Enum:
				valStr = fmt.Sprintf(`var v1 int32
					i, v1 = gsnet.ReadInt32(data, i)
					v := %s(v1)`, gen.typeName(hash.Value))
			case *ast.Table:
				valStr = fmt.Sprintf(`var size uint32
						v := %s
						i, size = gsnet.ReadUint32(data, i)
						if err = v.Unmarshal(data[i:i+int(size)]); err != nil {
							return
						}
						i += int(size)`, gen.defaultVal(hash.Value))
			default:
				log.Panicf("map key can only be int or string, %s not supported", hash.Key.Name())
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
	log.Panic("not here")
	return "unknown"
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
	log.Panic("not here")
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
	log.Panic("not here")
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

// callargs 根据参数生成函数的调用参数列表
func (gen *Gen4Go) callargs(params []*ast.Param) string {
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

// returnargs 根据函数生成 接收 函数的调用值 的 声明列表
func (gen *Gen4Go) returnargs(method *ast.Method) string {
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

	log.Panicf("inner error: unknown golang typeName: %s\n\t%s",
		expr,
		gslang.Pos(expr))
	return "unknown"
}

// enumType 根据枚举类型长度和符号 取对应的golang类型
func (gen *Gen4Go) enumType(enum *ast.Enum) string {
	switch {
	case enum.Length == 1 && enum.Signed == true:
		return "int8"
	case enum.Length == 1 && enum.Signed == false:
		return "byte"
	case enum.Length == 2 && enum.Signed == true:
		return "int16"
	case enum.Length == 2 && enum.Signed == false:
		return "uint16"
	case enum.Length == 4 && enum.Signed == true:
		return "int32"
	case enum.Length == 4 && enum.Signed == false:
		return "uint32"
	}
	log.Panicf("inner error: check enum ABI: %s\n\t%d",
		enum, enum.Length,
	)
	return ""
}

// sovFunc sov函数名字
func (gen *Gen4Go) sovFunc(script *ast.Script) string {
	name := script.Name()
	ss := strings.Split(name, ".")
	ret := fmt.Sprintf("sov%s", strings.Title(ss[0]))
	return ret
}

// writeFile 将代码节点对应的golang代码写入到文件
func (gen *Gen4Go) writeFile(script *ast.Script, bytes []byte) {
	fullPath, ok := gslang.FilePath(script)
	if !ok {
		log.Panic("inner error: compile must bind file path to script")
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
		log.Errorf("goimports format err:%s, check if goimports installed:\n\tgo install golang.org/x/tools/cmd/goimports@latest", err)
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
	for _, t := range script.Types {
		t.Accept(gen)
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
	log.Debugf("生成器")
	return contract
}
