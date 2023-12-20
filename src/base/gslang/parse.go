// -------------------------------------------
// @file      : parse.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/14 下午8:19
// -------------------------------------------

package gslang

import (
	"bytes"
	"fmt"
	"gogs/base/gslang/ast"
	log "gogs/base/logger"
	"math"
	"os"
	"path/filepath"
	"strings"
)

const (
	//  位置额外信息的key
	posExtra = "gslang_parser_pos"
	// 注释额外信息的key
	commentExtra = "gslang_parser_comment"
)

// attachPos 为某个节点添加额外的位置信息
func attachPos(node ast.Node, pos Position) {
	node.NewExtra(posExtra, pos)
}

// Pos 节点的位置信息
func Pos(node ast.Node) Position {
	if val, ok := node.Extra(posExtra); ok {
		return val.(Position)
	}
	return Position{
		Line:     0,
		Column:   0,
		Filename: "unknown",
	}
}

// attachComments 为节点增加额外的信息-注释列表
func attachComments(node ast.Node, comments []*Token) {
	node.NewExtra(commentExtra, comments)
}

// Comments 返回节点的注释列表信息
func Comments(node ast.Node) []*Token {
	if comments, ok := node.Extra(commentExtra); ok {
		return comments.([]*Token)
	}
	return nil
}

// Parser 分析器
type Parser struct {
	*Lexer               // 内嵌词法分析器
	compiler *Compiler   // 隶属的编译器
	script   *ast.Script // 指向的代码节点
	comments []*Token    // 注释列表
	attrs    []*ast.Attr // 属性列表
}

// Peek 从词法分析器 取当前Token
func (parser *Parser) Peek() *Token {
	token, err := parser.Lexer.Peek()
	if err != nil {
		log.Panic(err)
	}
	return token
}

// Next 从词法分析器取下一个Token
func (parser *Parser) Next() *Token {
	token, err := parser.Lexer.Next()
	if err != nil {
		log.Panic(err)
	}
	return token
}

// errorf 格式化报错
func (parser *Parser) errorf(position Position, template string, args ...interface{}) {
	log.Panicf(fmt.Sprintf("parse:%s error:%s", position.String(), fmt.Sprintf(template, args...)))
}

// expect 期望下一个Token的类型为目标rune expect,否则报错
func (parser *Parser) expect(expect rune) *Token {
	token := parser.Next()
	if token.Type != expect {
		parser.errorf(token.Pos, "expect '%s',but got '%s' ", TokenName(expect), TokenName(token.Type))
	}
	return token
}

// expectf 期望下一个Token的类型为目标rune expect,否则格式化报错
func (parser *Parser) expectf(expect rune, template string, args ...interface{}) *Token {
	token := parser.Next()
	if token.Type != expect {
		parser.errorf(token.Pos, fmt.Sprintf(template, args...))
	}
	return token
}

// parseTypeRef 分析类型引用 如 ast.TypeRef
func (parser *Parser) parseTypeRef() *ast.TypeRef {
	// 需求一个标识符 如ast
	start := parser.expect(TokenID)
	// 顺序 标识符 列表
	nodes := []string{start.Value.(string)}
	for {
		// 如果第偶数个Token不为. 则 认为标识符已到末尾
		token := parser.Peek()
		if token.Type != '.' {
			break
		}
		parser.Next()
		// 顺序 循环  添加到标识符列表
		token = parser.expect(TokenID)
		nodes = append(nodes, token.Value.(string))
	}
	// 用标识符列表在代码节点内新建类型引用
	ref := parser.script.NewTypeRef(nodes)
	// 给节点附加位置信息
	attachPos(ref, start.Pos)
	return ref
}

// parse 编译器进行分析流程
func (compiler *Compiler) parse(pkg *ast.Package, path string) (*ast.Script, error) {
	// 在目标代码包中新建代码节点 代码节点name为其相对文件名
	script, err := pkg.NewScript(filepath.Base(path))
	if err != nil {
		return nil, err
	}
	// 读取整个文件内容到一个字节切片 []byte
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// 新建分析器
	parser := &Parser{
		Lexer:    NewLexer(script.Name(), bytes.NewReader(content)), // 生成词法分析器
		compiler: compiler,                                          // 设置所属编译器
		script:   script,                                            // 分析器指向的代码节点
	}
	// 分析器进行分析
	err = parser.parse()
	// 返回分析后的代码节点
	log.Debug("解析完的代码", script.Name(), script.Types)
	return script, err
}

// parse 分析器入口函数
func (parser *Parser) parse() (err error) {
	// 捕获错误 并返回该错误
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	// 先分析 代码内导入的其他包
	parser.parseImports()
	for {
		// 分析是否有属性
		parser.parseAttrs()
		// 根据下一个Token类型决定下一步分析
		token := parser.Next()
		switch token.Type {
		case TokenEOF: // 文件末尾 跳出循环
			goto FINISH
		case KeyEnum: // enum关键字 则分析枚举
			parser.parseEnum()
		case KeyTable: // table 关键字
			parser.parseTable(false)
		case KeyStruct: // struct 关键字
			parser.parseTable(true)
		case KeyContract: // contract 关键字
			parser.parseContract()
		default: // 其余则报错
			parser.errorf(token.Pos, "expect EOF")
		}
	}
FINISH:
	// 注释列表以额外信息的形式 添加到代码节点
	// 剩余的注释列表及属性列表均附加到代码节点
	attachComments(parser.script, parser.comments)
	parser.script.AddAttrs(parser.attrs)
	return
}

// parseComments 分析注释token 并保存到分析器的注释列表
func (parser *Parser) parseComments() {
	for { // 循环判断TokenCOMMENT
		token := parser.Peek()
		if token.Type != TokenCOMMENT {
			return
		}
		parser.Next()
		parser.comments = append(parser.comments, token)
	}
}

// attachComments 将分析器保存的注释列表中符合条件的注释附加给对应节点
func (parser *Parser) attachComments(node ast.Node) {
	// 节点位置
	pos := Pos(node)
	// 节点的位置必须有效
	if !pos.Valid() {
		log.Panicf("all node have to bind with pos object by calling attachPos: %s", node)
	}

	// 用于保存选中的注释
	var selected []*Token
	// 用于保存没有被选中的注释
	var rest []*Token
	// 遍历分析器当前保存的注释列表
	for i := len(parser.comments) - 1; i >= 0; i-- {
		comment := parser.comments[i]
		// 注释节点所属文件名字必须和目标节点所属文件名相同
		if comment.Pos.Filename != pos.Filename {
			log.Panicf("comment's filename must equal with node's filename")
		}

		// 如果注释节点的行号与 目标节点行号相同 或者 小1 则认为该注释属于目标节点  递归往上找行号连续的注释
		if comment.Pos.Line == pos.Line || (comment.Pos.Line+1) == pos.Line {
			selected = append(selected, comment)
			pos = comment.Pos
		} else {
			rest = append(rest, comment)
		}
	}
	// 分析器保存未被选中的注释
	parser.comments = rest
	// 将被选中的注释列表反序 以此得到按行号递增的注释列表
	var revert []*Token
	for i := len(selected) - 1; i >= 0; i-- {
		revert = append(revert, selected[i])
	}
	// 将此结果注释列表附加给目标节点
	attachComments(node, revert)
}

// parseImports 分析 当前代码 需要导入的 包
func (parser *Parser) parseImports() {
	// 先分析有没有导入注释
	parser.parseComments()
	for {
		// 判断是否有导入关键字 import
		token := parser.Peek()
		if token.Type != KeyImport {
			break
		}
		parser.Next()
		// 如果接下来是字符串字面量或者标识符 则认为是单行引用 并进行解析
		token = parser.Peek()
		if token.Type == TokenSTRING || token.Type == TokenID {
			log.Debug(token.Value)
			ref := parser.parseImport()
			if ref == nil {
				log.Panicf("check parser.Import implement")
			}
			parser.parseComments()     // 行尾可能有注释
			parser.attachComments(ref) // 附加注释到目标
		} else if token.Type == '(' { // 检测到 ( 则认为是多行引用
			parser.Next()
			for { // 多行引入循环 直到token不是合法 引用包 为止
				parser.parseComments() // 目标上方行可能有注释
				if ref := parser.parseImport(); ref != nil {
					parser.parseComments()     // 行尾可能有的注释
					parser.attachComments(ref) // 附加注释到目标
					continue
				}
				break
			}
			parser.expect(')') // 期盼一个) 否则认为格式错误报错
		} else { // 非法格式
			parser.errorf(token.Pos, "expect import body: TokenString or '('")
		}
	}
	// 无论什么包都要默认引入base/gslang包 编译器自动引入 设置位置1,1
	if parser.script.Package().Name() != GSLangPackage &&
		parser.script.Imports["gslang"] == nil {
		pkg, err := parser.compiler.Compile(GSLangPackage)
		if err != nil {
			panic(err)
		}
		pos := Position{
			Filename: parser.script.Name(),
			Line:     1,
			Column:   1,
		}
		ref, ok := parser.script.NewPackageRef("gslang", pkg)
		if pkg == nil {
			log.Panicf("check Compiler and Compile implement")
		}
		if !ok {
			log.Panicf("check if the script manual import gslang package")
		}
		attachPos(ref, pos)
	}
}

// parserImport 分析导入单个包引用
func (parser *Parser) parseImport() *ast.PackageRef {
	// 取当前token
	token := parser.Peek()
	var path string
	var key string
	// 如果 token是字符串字面量 则表明是直接引入包路径 无别名
	if token.Type == TokenSTRING {
		path = token.Value.(string)
		key = filepath.Base(path) // 设置key为原始包名
		parser.Next()
	} else if token.Type == TokenID { // 是 标识符 则是为引用包提供别名 同golang语法
		parser.Next()
		key = token.Value.(string) // 设置key为 别名
		token = parser.expect(TokenSTRING)
		path = token.Value.(string)
	} else {
		return nil
	}
	// 编译目标路径的包
	pkg, err := parser.compiler.Compile(path)
	if err != nil {
		log.Panic(err)
	}
	// 将该包生成包引用节点并加入到代码节点的包引用列表中
	ref, ok := parser.script.NewPackageRef(key, pkg)
	// 检查是否已经引用了 同名的包
	if !ok {
		parser.errorf(token.Pos,
			"import same package(%s) twice: \n\tsee: %s",
			key, Pos(ref))
	}
	// 为目标包引用 添加 源文件中的位置
	attachPos(ref, token.Pos)
	return ref
}

// parseAttrs 分析属性
func (parser *Parser) parseAttrs() {
	// 先检查是否有注释
	parser.parseComments()
	for {
		token := parser.Peek()
		// 属性以符号 @ 开头
		if token.Type != '@' {
			return
		}
		parser.Next()
		// 通过类型引用分析在代码节点内新建属性
		attr := parser.script.NewAttr(parser.parseTypeRef())
		attachPos(attr, token.Pos)
		token = parser.Peek()
		if token.Type == '(' { // 如果后面跟了()则表示有参数列表 解析此参数列表附加到此属性
			parser.Next()
			// 分析参数列表附加到属性
			attr.Args = parser.parseArgs()
			parser.expect(')')
		}
		// 将属性添加到分析器的属性缓存列表
		parser.attrs = append(parser.attrs, attr)
		// 分析是否有注释
		parser.parseComments()
		// 将对应注释附加到此属性节点
		parser.attachComments(attr)
	}
}

// attachAttrs 附加属性
func (parser *Parser) attachAttrs(node ast.Node) {
	// 将分析器当前属性缓存列表中的所有属性附加到对应节点 并清空分析器属性缓存列表
	node.AddAttrs(parser.attrs)
	parser.attrs = nil
}

// parseArgs 分析参数列表
func (parser *Parser) parseArgs() ast.Expr {
	// 参数列表的分析 总是从跳过 ( 的下一个Token开始
	token := parser.Peek()
	// 开始的第一个Token 为 ) 则表示一个空的参数列表 直接返回 nil
	if token.Type == ')' {
		return nil
	}
	// TokenLABEL 是 源文件中 如 lang: 格式的Token  (lang:"en",age:1)
	if token.Type == TokenLABEL {
		// 代码内新建一个命名参数列表
		args := parser.script.NewNamedArgs()
		parser.Next()
		name := token
		for {
			if arg, ok := args.NewArg(token.Value.(string), parser.parseArg()); !ok {
				// 命令参数列表内已存在同名的参数
				parser.errorf(name.Pos, "duplicate param assign: \n\tsee: %s", Pos(arg))
			} else {
				// 分析注释并添加到对应参数
				parser.parseComments()
				parser.attachComments(arg)
			}
			// 检查是否是参数分隔符逗号,不是则弹出循环 返回命名参数列表
			token = parser.Peek()
			if token.Type != ',' {
				break
			}
			parser.Next()
			// 期待一个TokenLABEL
			name = parser.expect(TokenLABEL)
		}
		return args
	}
	// 新建一个参数列表
	args := parser.script.NewArgs()
	for {
		// 分析参数
		arg := args.NewArg(parser.parseArg())
		// 分析并附加注释到参数
		parser.parseComments()
		parser.attachComments(arg)
		token = parser.Peek()
		// 参数分隔符
		if token.Type != ',' {
			break
		}
	}
	return args
}

// parseArg 分析参数
func (parser *Parser) parseArg() ast.Expr {
	// 二元运算表达式
	var lhs *ast.BinaryOp
	for {
		token := parser.Peek()
		var rhs ast.Expr // 表达式对象
		switch token.Type {
		case TokenINT: // 字面量整数值  100
			parser.Next()
			rhs = parser.script.NewInt(token.Value.(int64))
		case TokenFLOAT: // 字面量浮点值 3.14
			parser.Next()
			rhs = parser.script.NewFloat(token.Value.(float64))
		case TokenSTRING: // 字面量字符串  "caibo"
			parser.Next()
			rhs = parser.script.NewString(token.Value.(string))
		case TokenTrue: // 字面量布尔值真 true
			parser.Next()
			rhs = parser.script.NewBool(true)
		case TokenFalse: // 字面量布尔值假 false
			parser.Next()
			rhs = parser.script.NewBool(false)
		case '-': // 字面量 负整数 负浮点数
			parser.Next()
			next := parser.Next()
			if next.Type == TokenINT {
				rhs = parser.script.NewInt(-next.Value.(int64))
			} else if next.Type == TokenFLOAT {
				rhs = parser.script.NewFloat(-next.Value.(float64))
			}
			parser.errorf(token.Pos, "unexpect token '-'")
		case '+': // 字面量 正整数  正浮点数
			parser.Next()
			next := parser.Next()
			if next.Type == TokenINT {
				rhs = parser.script.NewInt(next.Value.(int64))
			} else if next.Type == TokenFLOAT {
				rhs = parser.script.NewFloat(next.Value.(float64))
			}
			parser.errorf(token.Pos, "unexpect token '+'")
		case TokenID: // 标识符 节点对象
			rhs = parser.parseTypeRef()
		default:
			parser.errorf(token.Pos, "unexpect token '%s', expect argument stmt	", TokenName(token.Type))
		}
		attachPos(rhs, token.Pos)
		if lhs != nil { // 已经是第二个操作数的时候 保存为右操作数
			lhs.Right = rhs
		}
		token = parser.Peek()
		if token.Type == '|' {
			parser.Next()
			if lhs != nil { // 递归的保存方式
				lhs = parser.script.NewBinaryOp("|", lhs, nil)
			} else {
				lhs = parser.script.NewBinaryOp("|", rhs, nil)
			}
			attachPos(lhs, token.Pos)
			continue
		}
		if lhs != nil { // 有2个以上操作数 用|连接的
			return lhs
		}
		// 只有一个操作数
		return rhs
	}
}

// parseContract	分析协议(一组函数)
func (parser *Parser) parseContract() {
	// contract后第一个标识符为协议名字
	name := parser.expect(TokenID)
	contract := parser.script.NewService(name.Value.(string))
	// 协议也认为是类型 代码包内不能有同名协议
	if old, ok := parser.script.NewType(contract); !ok {
		parser.errorf(name.Pos, "duplicate type name:\n\tsee: %s", Pos(old))
	}
	// 附加位置信息到协议节点
	attachPos(contract, name.Pos)
	// 附加注释和属性
	parser.attachComments(contract)
	parser.attachAttrs(contract)
	token := parser.Peek()
	// 协议名后如果跟小括号则代表继承自某个协议类型
	if token.Type == '(' {
		parser.Next()
		for {
			// 分析注释
			parser.parseComments()
			// 分析类型引用 此处目标类型引用为某个协议类型
			base := parser.parseTypeRef()
			// 添加引用到协议的父类型列表
			if old, ok := contract.NewBase(base); ok {
				// 分析注释并附加注释到 父类型引用节点
				parser.parseComments()
				parser.attachComments(base)
			} else { // 不能重复继承相同协议
				parser.errorf(Pos(base), "duplicate inherit from same contract:\n\tsee: %s", Pos(old))
			}
			next := parser.Peek()
			// ,分隔多个父协议
			if next.Type == ',' {
				parser.Next()
				continue
			}
			break
		}
		parser.expect(')')
	}
	// 协议体以大括号开始
	parser.expect('{')
	for {
		// 分析属性
		parser.parseAttrs()
		token := parser.Peek()
		if token.Type != TokenID {
			break
		}
		// 取函数名字并在协议内新建函数节点
		methodName := parser.Next()
		method, ok := contract.NewMethod(methodName.Value.(string))
		if !ok {
			// 单个协议内不能有同名函数
			parser.errorf(methodName.Pos, "duplicate method name:\n\tsee: %s", Pos(method))
		}
		// 附加位置
		attachPos(method, methodName.Pos)
		// 取函数参数列表
		parser.expect('(')
		next := parser.Peek()
		// 非空参数列表
		if next.Type != ')' {
			for {
				// 分析属性
				parser.parseAttrs()
				// 分析类型
				paramType := parser.parseType()
				next = parser.Peek()
				if next.Type != ',' &&
					next.Type != ')' &&
					next.Type != TokenCOMMENT {
					paramType = parser.parseType()
				}
				// 添加类型到函数的输入参数列表
				param := method.NewParam(paramType)
				// 附加位置信息 注释及属性
				attachPos(param, Pos(param.Type))
				parser.parseComments()
				parser.attachComments(param)
				parser.attachAttrs(param)
				next = parser.Peek()
				// 多个输入参数以逗号间隔
				if next.Type == ',' {
					parser.Next()
					continue
				}
				break
			}
		}
		parser.expect(')')
		next = parser.Peek()
		// 函数输入参数后如果有->符号则表示有返回参数列表 分析基本同输入参数
		if next.Type == TokenArrowRight {
			parser.Next()
			parser.expect('(')
			for {
				parser.parseAttrs()
				paramType := parser.parseType()
				next := parser.Peek()
				if next.Type != ',' &&
					next.Type != ')' &&
					next.Type != TokenCOMMENT {
					paramType = parser.parseType()
				}
				// 将类型添加到返回参数列表
				param := method.NewReturn(paramType)
				attachPos(param, Pos(param.Type))
				parser.parseComments()
				parser.attachComments(param)
				parser.attachAttrs(param)
				next = parser.Peek()
				if next.Type == ',' {
					parser.Next()
					continue
				}
				break
			}
			parser.expect(')')
		}
		// 多个函数声明以分好分隔
		parser.expect(';')
		// 给函数附加注释和属性
		parser.parseComments()
		parser.attachComments(method)
		parser.attachAttrs(method)
	}
	parser.expect('}')
}

// newGSLangAttr 在代码节点内生成指定名字的类型引用 如果不是gslang包下的 还需加入gslang.前缀,并用此类型引用节点生成一个属性
func (parser *Parser) newGSLangAttr(name string) *ast.Attr {
	if parser.script.Package().Name() != GSLangPackage {
		return parser.script.NewAttr(parser.script.NewTypeRef([]string{
			"gslang", name,
		}))
	}
	return parser.script.NewAttr(parser.script.NewTypeRef([]string{name}))
}

// newGSLangTypeRef 在代码节点内生成指定名字的类型引用 如果不是gslang包下的 还需加入gslang.前缀
func (parser *Parser) newGSLangTypeRef(name string) *ast.TypeRef {
	if parser.script.Package().Name() != GSLangPackage {
		return parser.script.NewTypeRef([]string{"gslang", name})
	}
	return parser.script.NewTypeRef([]string{name})
}

// parseType 分析类型 返回一个表达式接口
func (parser *Parser) parseType() ast.Expr {
	token := parser.Peek()
	switch token.Type {
	// [] 代表数组或切片
	case '[':
		parser.Next()
		next := parser.Peek()
		length := uint32(0)
		// 有长度的数组 无长度的切片
		if next.Type == TokenINT {
			parser.Next()
			val := next.Value.(int64)
			if val < 1 || val > math.MaxUint32 {
				parser.errorf(next.Pos, "array length out of range: %d", val)
			}
			length = uint32(val)
		}
		parser.expect(']')
		// 递归分析类型
		element := parser.parseType()
		switch element.(type) {
		case *ast.List, *ast.Array:
			// 不支持递归数组或切片
			parser.errorf(token.Pos, "gslang didn't support Recursively define array or list")
		}
		// 包装并返回 对应类型的数组或切片
		var expr ast.Expr
		if length > 0 {
			expr = parser.script.NewArray(length, element)
		} else {
			expr = parser.script.NewList(element)
		}
		attachPos(expr, token.Pos)
		return expr
	case KeyMap:
		parser.Next()
		parser.expect('[')
		// 分析key
		key := parser.parseType()
		switch key.(type) {
		case *ast.List, *ast.Array, *ast.Map:
			parser.errorf(token.Pos, "gslang didn't support key(map array list) for map")
		}
		parser.expect(']')
		// 分析value
		value := parser.parseType()
		switch value.(type) {
		case *ast.List, *ast.Array, *ast.Map:
			parser.errorf(token.Pos, "gslang didn't support value(map array list) for map")
		}
		// 包装map并返回
		var expr ast.Expr
		expr = parser.script.NewMap(key, value)
		attachPos(expr, token.Pos)
		return expr
	case KeyByte, KeySByte, KeyInt16, KeyUInt16, KeyInt32, KeyUInt32,
		KeyInt64, KeyUInt64, KeyBool, KeyFloat32, KeyFloat64, KeyString:
		// gslang内置数据类型
		parser.Next()
		// 生成类型引用并返回
		expr := parser.newGSLangTypeRef(strings.Title(TokenName(token.Type)))
		attachPos(expr, token.Pos)
		return expr
	case TokenID:
		// 如果是非特殊标识符 则按路径生成类型引用 并返回
		expr := parser.parseTypeRef()
		attachPos(expr, token.Pos)
		return expr
	default:
		parser.errorf(token.Pos, "expect type declare")
	}
	return nil
}

// parseTable 分析表(isStruct=false) 结构体(isStruct=true)
func (parser *Parser) parseTable(isStruct bool) {
	name := parser.expect(TokenID)
	table := parser.script.NewTable(name.Value.(string))
	// 不能有重名类型
	if old, ok := parser.script.NewType(table); !ok {
		parser.errorf(name.Pos, "duplicate type name:\n\tsee: %s", Pos(old))
	}
	// 附加位置 注释 属性
	attachPos(table, name.Pos)
	parser.attachComments(table)
	parser.attachAttrs(table)
	if isStruct {
		// 如果是结构体 生成一个属性
		attr := parser.newGSLangAttr("Struct")
		attachPos(attr, name.Pos)
		// 将属性附加到这个表
		table.AddAttr(attr)
	}
	parser.expect('{')
	for { // 分析表或者结构体的域
		parser.parseAttrs()
		token := parser.Peek()
		if token.Type != TokenID {
			break
		}
		fieldName := parser.expect(TokenID)
		// 表或结构体中新建一个域
		field, ok := table.NewField(fieldName.Value.(string))
		if !ok { // 不能有重名域
			parser.errorf(fieldName.Pos, "duplicate field name:\n\tsee: %s", Pos(field))
		}
		// 附加位置
		attachPos(field, fieldName.Pos)
		// 分析域的类型
		field.Type = parser.parseType()
		// 域间用分号分隔
		parser.expect(';')
		// 分析注释 附加注释 附加属性
		parser.parseComments()
		parser.attachComments(field)
		parser.attachAttrs(field)
	}
	parser.expect('}')
}

// parseEnumBase 根据枚举标识符后面括号内的数值类型标识符确定枚举值的长度和有无符号
// 括号内枚举类型仅支持内置几种类型
func (parser *Parser) parseEnumBase() (length uint, signed bool) {
	parser.expect('(')
	token := parser.Next()
	switch token.Type {
	case KeyByte:
		length = 1
	case KeySByte:
		length = 1
		signed = true
	case KeyInt16:
		length = 2
		signed = true
	case KeyUInt16:
		length = 2
		signed = false
	case KeyInt32:
		length = 4
		signed = true
	case KeyUInt32:
		length = 4
		signed = false
	default:
		parser.errorf(token.Pos, "enum must inherit from integer types, got: %s", TokenName(token.Type))
	}
	parser.expect(')')
	return
}

// parseEnum 分析枚举
func (parser *Parser) parseEnum() {
	// 枚举名字
	name := parser.expect(TokenID)
	token := parser.Peek()
	length := uint(1)
	signed := false
	// 如果带括号 则根据括号内的类型标识符关键字 确定 枚举值的长度和符号
	if token.Type == '(' {
		length, signed = parser.parseEnumBase()
	}
	// 新建枚举
	enum := parser.script.NewEnum(name.Value.(string), length, signed)
	// 枚举作为一中类型添加包及代码节点 且不能有重名类型
	if old, ok := parser.script.NewType(enum); !ok {
		parser.errorf(name.Pos, "duplicate type name:\n\tsee: %s", Pos(old))
	}
	// 附加位置 注释 属性
	attachPos(enum, name.Pos)
	parser.attachComments(enum)
	parser.attachAttrs(enum)
	parser.expect('{')
	for { // 循环分析枚举的每个枚举值,其中枚举值可以为负值
		// 分析属性
		parser.parseAttrs()
		token := parser.expectf(TokenID, "expect enum value field")
		parser.expect('(')
		next := parser.Peek()
		negative := false
		if next.Type == '-' {
			parser.Next()
			negative = true
		}
		valToken := parser.expect(TokenINT)
		val := valToken.Value.(int64)
		if negative {
			val = -val
		}
		switch { // 判断值是否越界
		case enum.Length == 1 && enum.Signed:
			if val > math.MaxInt8 || val < math.MinInt8 {
				parser.errorf(valToken.Pos, "out of enum[%s] type's range", enum)
			}
		case enum.Length == 1 && !enum.Signed:
			if val > math.MaxUint8 || val < 0 {
				parser.errorf(valToken.Pos, "out of enum[%s] type's range", enum)
			}
		case enum.Length == 2 && enum.Signed:
			if val > math.MaxInt16 || val < math.MinInt16 {
				parser.errorf(valToken.Pos, "out of enum[%s] type's range", enum)
			}
		case enum.Length == 2 && !enum.Signed:
			if val > math.MaxUint16 || val < 0 {
				parser.errorf(valToken.Pos, "out of enum[%s] type's range", enum)
			}
		case enum.Length == 4 && enum.Signed:
			if val > math.MaxInt32 || val < math.MinInt32 {
				parser.errorf(valToken.Pos, "out of enum[%s] type's range", enum)
			}
		case enum.Length == 4 && !enum.Signed:
			if val > math.MaxUint32 || val < 0 {
				parser.errorf(valToken.Pos, "out of enum[%s] type's range", enum)
			}
		}
		parser.expect(')')
		// 在枚举内新建单挑枚举值
		enumVal, ok := enum.NewEnumVal(token.Value.(string), val)
		if !ok { // 不能有重名枚举值
			parser.errorf(token.Pos,
				"duplicate enum val name(%s):\n\tsee: %s",
				enumVal.Name(), Pos(enumVal))
		}
		attachPos(enumVal, token.Pos)
		parser.attachAttrs(enumVal)
		next = parser.Peek()
		if next.Type != ',' { // 枚举值之间用逗号分隔
			parser.parseComments()
			parser.attachComments(enumVal)
			break
		}
		parser.Next()
		parser.parseComments()
		parser.attachComments(enumVal)
	}
	parser.expect('}')
}
