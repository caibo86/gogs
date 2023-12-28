// -------------------------------------------
// @file      : eval_attr.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/22 下午2:45
// -------------------------------------------

package gslang

import (
	"gogs/base/gslang/ast"
	log "gogs/base/logger"
)

// evalAttr 访问枚举值
type evalAttr struct {
	values map[string]any
}

// VisitAttr 仅为实现访问者接口
func (visitor *evalAttr) VisitAttr(node *ast.Attr) ast.Node {
	if visitor == nil {
		log.Panicf("visitor(evalAttr) should not be nil")
	}
	if visitor.values == nil {
		visitor.values = make(map[string]any)
	}
	expr := node.Type.Ref
	table, ok := expr.(*ast.Table)
	if !ok {
		log.Panicf("attr type should be ast.Table")
	}
	args := node.Args
	if args == nil {
		for _, field := range table.Fields {
			switch field.Type.(type) {
			case *ast.Array, *ast.Slice, *ast.Map:
				visitor.values[field.Name()] = nil
			case *ast.TypeRef:
				ref := field.Type.(*ast.TypeRef).Ref
				name := ref.Name()
				switch name {
				case "Byte":
					visitor.values[field.Name()] = byte(0)
				case "SByte":
					visitor.values[field.Name()] = int8(0)
				case "Int16":
					visitor.values[field.Name()] = int16(0)
				case "UInt16":
					visitor.values[field.Name()] = uint16(0)
				case "Int32":
					visitor.values[field.Name()] = int32(0)
				case "UInt32":
					visitor.values[field.Name()] = uint32(0)
				case "Int64":
					visitor.values[field.Name()] = int64(0)
				case "UInt64":
					visitor.values[field.Name()] = uint64(0)
				case "Float32":
					visitor.values[field.Name()] = float32(0)
				case "Float64":
					visitor.values[field.Name()] = float64(0)
				case "String":
					visitor.values[field.Name()] = ""
				case "Bool":
					visitor.values[field.Name()] = false
				default:
					visitor.values[field.Name()] = nil
				}
			}
		}
	} else if uArgs, ok := args.(*ast.Args); ok {
		for _, field := range table.Fields {
			var item ast.Expr
			if field.ID < uint16(len(uArgs.Items)) {
				item = uArgs.Items[field.ID]
			}
			switch field.Type.(type) {
			case *ast.Array, *ast.Slice, *ast.Map:
				visitor.values[field.Name()] = nil
			case *ast.TypeRef:
				ref := field.Type.(*ast.TypeRef).Ref
				name := ref.Name()
				switch name {
				case "Byte", "SByte", "Int16", "UInt16",
					"Int32", "UInt32", "Int64", "UInt64":
					if item == nil {
						visitor.values[field.Name()] = int64(0)
					} else {
						i := item.(*ast.Int)
						visitor.values[field.Name()] = i.Value
					}
				case "Float32":
					if item == nil {
						visitor.values[field.Name()] = float32(0)
					} else {
						i := item.(*ast.Float)
						visitor.values[field.Name()] = float32(i.Value)
					}
				case "Float64":
					if item == nil {
						visitor.values[field.Name()] = float64(0)
					} else {
						i := item.(*ast.Float)
						visitor.values[field.Name()] = i.Value
					}
				case "String":
					if item == nil {
						visitor.values[field.Name()] = float64(0)
					} else {
						i := item.(*ast.String)
						visitor.values[field.Name()] = i.Value
					}
				case "Bool":
					if item == nil {
						visitor.values[field.Name()] = false
					} else {
						i := item.(*ast.Bool)
						visitor.values[field.Name()] = i.Value
					}
				default:
					// log.Panicf("不支持的类型:%s", name)
				}
			}
		}
	} else if nArgs, ok := args.(*ast.NamedArgs); ok {
		for _, field := range table.Fields {
			var item ast.Expr
			fieldName := field.Name()
			item = nArgs.Items[fieldName]
			switch field.Type.(type) {
			case *ast.Array, *ast.Slice, *ast.Map:
				visitor.values[field.Name()] = nil
			case *ast.TypeRef:
				ref := field.Type.(*ast.TypeRef).Ref
				name := ref.Name()
				switch name {
				case "Byte", "SByte", "Int16", "UInt16",
					"Int32", "UInt32", "Int64", "UInt64":
					if item == nil {
						visitor.values[field.Name()] = int64(0)
					} else {
						i := item.(*ast.Int)
						visitor.values[field.Name()] = i.Value
					}
				case "Float32":
					if item == nil {
						visitor.values[field.Name()] = float32(0)
					} else {
						i := item.(*ast.Float)
						visitor.values[field.Name()] = float32(i.Value)
					}
				case "Float64":
					if item == nil {
						visitor.values[field.Name()] = float64(0)
					} else {
						i := item.(*ast.Float)
						visitor.values[field.Name()] = i.Value
					}
				case "String":
					if item == nil {
						visitor.values[field.Name()] = float64(0)
					} else {
						i := item.(*ast.String)
						visitor.values[field.Name()] = i.Value
					}
				case "Bool":
					if item == nil {
						visitor.values[field.Name()] = false
					} else {
						i := item.(*ast.Bool)
						visitor.values[field.Name()] = i.Value
					}
				default:
					// log.Panicf("不支持的类型:%s", name)
				}
			}
		}
	} else {
		log.Panicf("attr args should be nil or ast.Args or ast.NamedArgs")
	}

	return nil
}

// VisitBinaryOp 访问二元运算
func (visitor *evalAttr) VisitBinaryOp(node *ast.BinaryOp) ast.Node {
	log.Panicf("stmt is not attr expr :%s", Pos(node))
	return nil
}

// VisitTypeRef 访问类型引用
func (visitor *evalAttr) VisitTypeRef(node *ast.TypeRef) ast.Node {
	log.Panicf("stmt is not attr expr :%s", Pos(node))
	return node
}

// VisitEnumVal 访问枚举值
func (visitor *evalAttr) VisitEnumVal(node *ast.EnumVal) ast.Node {
	log.Panicf("stmt is not attr expr :%s", Pos(node))
	return node
}

// VisitString 仅为实现访问者接口
func (visitor *evalAttr) VisitString(node *ast.String) ast.Node {
	log.Panicf("stmt is not attr expr :%s", Pos(node))
	return nil
}

// VisitFloat 仅为实现访问者接口
func (visitor *evalAttr) VisitFloat(node *ast.Float) ast.Node {
	log.Panicf("stmt is not attr expr :%s", Pos(node))
	return nil
}

// VisitInt 仅为实现访问者接口
func (visitor *evalAttr) VisitInt(node *ast.Int) ast.Node {
	log.Panicf("stmt is not attr expr :%s", Pos(node))
	return nil
}

// VisitBool 仅为实现访问者接口
func (visitor *evalAttr) VisitBool(node *ast.Bool) ast.Node {
	log.Panicf("stmt is not attr expr :%s", Pos(node))
	return nil
}

// VisitPackage 仅为实现访问者接口
func (visitor *evalAttr) VisitPackage(node *ast.Package) ast.Node {
	log.Panicf("stmt is not attr expr :%s", Pos(node))
	return nil
}

// VisitScript 仅为实现访问者接口
func (visitor *evalAttr) VisitScript(node *ast.Script) ast.Node {
	log.Panicf("stmt is not attr expr :%s", Pos(node))
	return nil
}

// VisitEnum 仅为实现访问者接口
func (visitor *evalAttr) VisitEnum(node *ast.Enum) ast.Node {
	log.Panicf("stmt is not attr expr :%s", Pos(node))
	return nil
}

// VisitTable 仅为实现访问者接口
func (visitor *evalAttr) VisitTable(node *ast.Table) ast.Node {
	log.Panicf("stmt is not attr expr :%s", Pos(node))
	return nil
}

// VisitField 仅为实现访问者接口
func (visitor *evalAttr) VisitField(node *ast.Field) ast.Node {
	log.Panicf("stmt is not attr expr :%s", Pos(node))
	return nil
}

// VisitContract 仅为实现访问者接口
func (visitor *evalAttr) VisitContract(node *ast.Contract) ast.Node {
	log.Panicf("stmt is not attr expr :%s", Pos(node))
	return nil
}

// VisitMethod 仅为实现访问者接口
func (visitor *evalAttr) VisitMethod(node *ast.Method) ast.Node {
	log.Panicf("stmt is not attr expr :%s", Pos(node))
	return nil
}

// VisitArray 仅为实现访问者接口
func (visitor *evalAttr) VisitArray(node *ast.Array) ast.Node {
	log.Panicf("stmt is not attr expr :%s", Pos(node))
	return nil
}

// VisitSlice 仅为实现访问者接口
func (visitor *evalAttr) VisitSlice(node *ast.Slice) ast.Node {
	log.Panicf("stmt is not attr expr :%s", Pos(node))
	return nil
}

// VisitMap 仅为实现访问者接口
func (visitor *evalAttr) VisitMap(node *ast.Map) ast.Node {
	log.Panicf("stmt is not attr expr :%s", Pos(node))
	return nil
}

// VisitArgs 仅为实现访问者接口
func (visitor *evalAttr) VisitArgs(node *ast.Args) ast.Node {
	log.Panicf("stmt is not attr expr :%s", Pos(node))
	return nil
}

// VisitNamedArgs 仅为实现访问者接口
func (visitor *evalAttr) VisitNamedArgs(node *ast.NamedArgs) ast.Node {
	log.Panicf("stmt is not attr expr :%s", Pos(node))
	return nil
}

// VisitParam 仅为实现访问者接口
func (visitor *evalAttr) VisitParam(node *ast.Param) ast.Node {
	log.Panicf("stmt is not attr expr :%s", Pos(node))
	return nil
}
