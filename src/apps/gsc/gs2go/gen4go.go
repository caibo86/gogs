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

// 包名映射的引入包的go代码
var packageMapping = map[string]string{
	"gsnet.": `import "gogs/base/gsnet"`,
	// "yfdocker.": `import "gsgo/base/docker"`,
	// "yfconfig.": `import "gsgo/base/config"`,
	// "yferrors.": `import "gsgo/base/errors"`,
	"bytes.": `import "bytes"`,
	"fmt.":   `import "fmt"`,
	"time.":  `import "time"`,
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
		"enumWrite":    gen.enumWrite,
		"enumRead":     gen.enumRead,
		"writeTagType": gen.writeTagType,
		"pos":          gslang.Pos,
		"tag":          gen.tag,
	}
	gen.tpl, err = template.New("golang").Funcs(funcs).Parse(tpl4go)
	return
}

// tag 取表达式类型的标签
func (gen *Gen4Go) tag(expr ast.Expr) string {
	switch expr.(type) {
	case *ast.TypeRef:
		if tag, ok := tagMapping[expr.Name()]; ok {
			return tag
		}
		ref := expr.(*ast.TypeRef)
		if _, ok := ref.Ref.(*ast.Enum); ok {
			return "gsnet.Enum"
		}
		if gslang.IsStruct(ref.Ref.(*ast.Struct)) {
			return "gsnet.Struct"
		}
		return "gsnet.Table"
	case *ast.Array:
		return "gsnet.Array"
	case *ast.List:
		return "gsnet.List"
	}
	log.Panic("not here")
	return "unknown"
}

// enumWrite 根据枚举的值类型取写入语句
func (gen *Gen4Go) enumWrite(enum *ast.Enum) string {
	switch {
	case enum.Length == 1 && enum.Signed == true:
		return "gsnet.WriteSByte(writer,int8(val))"
	case enum.Length == 1 && enum.Signed == false:
		return "gsnet.WriteByte(writer,byte(val))"
	case enum.Length == 2 && enum.Signed == true:
		return "gsnet.WriteInt16(writer,int16(val))"
	case enum.Length == 2 && enum.Signed == false:
		return "gsnet.WriteUint16(writer,uint16(val))"
	case enum.Length == 4 && enum.Signed == true:
		return "gsnet.WriteInt32(writer,int32(val))"
	case enum.Length == 4 && enum.Signed == false:
		return "gsnet.WriteUint32(writer,uint32(val))"
	}
	log.Panic("not here")
	return "unknown"
}

// enumRead 根据枚举的值类型取读入语句
func (gen *Gen4Go) enumRead(enum *ast.Enum) string {
	switch {
	case enum.Length == 1 && enum.Signed == true:
		return "gsnet.ReadSByte(reader)"
	case enum.Length == 1 && enum.Signed == false:
		return "gsnet.ReadByte(reader)"
	case enum.Length == 2 && enum.Signed == true:
		return "gsnet.ReadInt16(reader)"
	case enum.Length == 2 && enum.Signed == false:
		return "gsnet.ReadUint16(reader)"
	case enum.Length == 4 && enum.Signed == true:
		return "gsnet.ReadInt32(reader)"
	case enum.Length == 4 && enum.Signed == false:
		return "gsnet.ReadUint32(reader)"
	}
	log.Panic("not here")
	return "unknown"
}

// writeTagType 根据类型取 带标签的写入函数
func (gen *Gen4Go) writeTagType(expr ast.Expr) string {
	switch expr.(type) {
	case *ast.TypeRef:
		// 内置类型
		if f, ok := writeTagMapping[expr.Name()]; ok {
			return f
		}
		// 自定义类型
		ref := expr.(*ast.TypeRef)
		if _, ok := expr.Script().Imports[ref.NamePath[0]]; ok {
			return fmt.Sprintf(
				"%s.WriteTag%s",
				ref.NamePath[0],
				strings.Title(ref.NamePath[1]),
			)
		}
		return fmt.Sprintf("WriteTag%s", strings.Title(ref.NamePath[0]))
	case *ast.Array:
		// 数组
		array := expr.(*ast.Array)
		var buff bytes.Buffer
		if array.Element.Name() == string(".yflang.Byte") {
			if err := gen.tpl.ExecuteTemplate(&buff, "writeTagByteArray", array); err != nil {
				panic(err)
			}
		} else {
			if err := gen.tpl.ExecuteTemplate(&buff, "writeTagArray", array); err != nil {
				panic(err)
			}
		}
		return buff.String()
	case *ast.List:
		// 切片
		list := expr.(*ast.List)
		var buff bytes.Buffer
		if list.Element.Name() == string(".yflang.Byte") {
			if err := gen.tpl.ExecuteTemplate(&buff, "writeTagByteList", list); err != nil {
				panic(err)
			}
		} else {
			if err := gen.tpl.ExecuteTemplate(&buff, "writeTagList", list); err != nil {
				panic(err)
			}
		}
		return buff.String()
	}
	log.Panic("not here")
	return "unknown"
}

// readType 根据类型取读入函数
func (gen *Gen4Go) readType(expr ast.Expr) string {
	switch expr.(type) {
	case *ast.TypeRef:
		// 内置类型
		if f, ok := readMapping[expr.Name()]; ok {
			return f
		}
		// 自定义类型
		ref := expr.(*ast.TypeRef)
		if _, ok := expr.Script().Imports[ref.NamePath[0]]; ok {
			return fmt.Sprintf(
				"%s.Read%s",
				ref.NamePath[0],
				strings.Title(ref.NamePath[1]),
			)
		}
		return fmt.Sprintf(
			"Read%s",
			strings.Title(ref.NamePath[0]),
		)
	case *ast.Array:
		// 数组
		array := expr.(*ast.Array)
		var buff bytes.Buffer
		if array.Element.Name() == string(".yflang.Byte") {
			if err := gen.tpl.ExecuteTemplate(&buff, "readByteArray", array); err != nil {
				panic(err)
			}
		} else {
			if err := gen.tpl.ExecuteTemplate(&buff, "readArray", array); err != nil {
				panic(err)
			}
		}
		return buff.String()
	case *ast.List:
		// 切片
		list := expr.(*ast.List)
		var buff bytes.Buffer
		if list.Element.Name() == string(".yflang.Byte") {
			if err := gen.tpl.ExecuteTemplate(&buff, "readByteList", list); err != nil {
				panic(err)
			}
		} else {
			if err := gen.tpl.ExecuteTemplate(&buff, "readList", list); err != nil {
				panic(err)
			}
		}
		return buff.String()
	case *ast.Map:
		// 字典
		hash := expr.(*ast.Map)
		var buff bytes.Buffer
		if err := gen.tpl.ExecuteTemplate(&buff, "readMap", hash); err != nil {
			panic(err)
		}
		return buff.String()
	}
	log.Panic("not here")
	return "unknown"
}

// writeType 根据类型取写入函数
func (gen *Gen4Go) writeType(expr ast.Expr) string {
	switch expr.(type) {
	case *ast.TypeRef:
		// 内置类型
		if f, ok := writeMapping[expr.Name()]; ok {
			return f
		}
		// 自定义类型
		ref := expr.(*ast.TypeRef)
		if _, ok := expr.Script().Imports[ref.NamePath[0]]; ok {
			return fmt.Sprintf(
				"%s.Write%s",
				ref.NamePath[0],
				strings.Title(ref.NamePath[1]),
			)
		}
		return fmt.Sprintf("Write%s", strings.Title(ref.NamePath[0]))
	case *ast.Array:
		// 数组
		array := expr.(*ast.Array)
		var buff bytes.Buffer
		if array.Element.Name() == string(".yflang.Byte") {
			if err := gen.tpl.ExecuteTemplate(&buff, "writeByteArray", array); err != nil {
				panic(err)
			}
		} else {
			if err := gen.tpl.ExecuteTemplate(&buff, "writeArray", array); err != nil {
				panic(err)
			}
		}
		return buff.String()
	case *ast.List:
		// 切片
		list := expr.(*ast.List)
		var buff bytes.Buffer
		if list.Element.Name() == string(".yflang.Byte") {
			if err := gen.tpl.ExecuteTemplate(&buff, "writeByteList", list); err != nil {
				panic(err)
			}
		} else {
			if err := gen.tpl.ExecuteTemplate(&buff, "writeList", list); err != nil {
				panic(err)
			}
		}
		return buff.String()
	case *ast.Map:
		// 字典
		hash := expr.(*ast.Map)
		var buff bytes.Buffer
		if err := gen.tpl.ExecuteTemplate(&buff, "writeMap", hash); err != nil {
			panic(err)
		}
		return buff.String()
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
		return "nil"
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
	log.Debugf("write to file:%s success", fullPath)
	// cmd := exec.Command("go", "fmt", fullPath)
	// _, err = cmd.Output()
	// if err != nil {
	// 	log.Debugf("format err:%s, check go fmt", err)
	// }
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
		log.Debug("生成器")
		script.Accept(gen)
	}
	return pkg
}

// VisitScript 访问代码
func (gen *Gen4Go) VisitScript(script *ast.Script) ast.Node {
	gen.buff.Reset()
	// 轮询访问代码中的所有类型 Enum Table Struct Contract
	for _, ctype := range script.Types {
		log.Debugf("生成器 %v", ctype.Name())
		ctype.Accept(gen)
	}
	log.Debug("生成器")
	// 代码中有类型
	if gen.buff.Len() > 0 {
		var buff bytes.Buffer
		// 写入额外信息
		buff.WriteString(fmt.Sprintf("// -------------------------------------------\n// @file      : gsc.go\n// @author    : generated by gsc, do not edit\n// @contact   : caibo923@gmail.com\n// @time      : %s\n// -------------------------------------------\n\n", time.Now().Format(time.RFC3339)))

		// 		buff.WriteString(
		// 			`/**************************************************/
		// /*     @author:caibo                              */
		// /*     @mail:48904088@qq.com                      */
		// /*     @date:2017-06-01                           */
		// /*     @module:GSC                                */
		// /*     @desc:自动生成Golang代码                     */
		// /**************************************************/
		//
		// `)

		// 写入包声明
		buff.WriteString(fmt.Sprintf("package %s\n", filepath.Base(script.Package().Name())))
		codes := gen.buff.String()
		// 如果两个特定的内置包中的gs文件 则不需要加包前缀
		if script.Package().Name() == "yf/platform/yfnet" {
			codes = strings.Replace(codes, "yfnet.", "", -1)
		}
		if script.Package().Name() == "yf/platform/yfdocker" {
			codes = strings.Replace(codes, "yfdocker.", "", -1)
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
				buff.WriteString(fmt.Sprintf("import %s \"%s\"\n", ref.Name(), ref.Ref))
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
func (gen *Gen4Go) VisitTable(table *ast.Struct) ast.Node {
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
func (gen *Gen4Go) VisitContract(contract *ast.Service) ast.Node {
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
