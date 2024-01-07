// -------------------------------------------
// @file      : formatter.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/28 下午8:40
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
	"path/filepath"
	"strings"
)

// printAttrs 输出格式化后属性
func printAttrs(buff *bytes.Buffer, node ast.Node) {
	for _, attr := range node.Attrs() {
		if attr.Name() != ".Struct" && attr.Name() != ".cblang.Struct" {
			printComments(buff, attr)
			buff.WriteString(fmt.Sprintf("%s\n", attr.OriginName()))
		}
	}
}

// printCommentsToLine 输出格式化后的注释到一行
func printCommentsToLine(buff *bytes.Buffer, node ast.Node) {
	comments := cblang.Comments(node)
	if len(comments) > 0 {
		buff.WriteString("//")
		for _, comment := range comments {
			value := comment.Value.(string)
			value = strings.TrimLeft(value, " ")
			buff.WriteString(fmt.Sprintf("%s", comment.Value))
		}
	}
}

// printComments 输出格式化后的注释
func printComments(buff *bytes.Buffer, node ast.Node) bool {
	comments := cblang.Comments(node)
	if len(comments) > 0 {
		for _, comment := range comments {
			value := comment.Value.(string)
			value = strings.TrimLeft(value, " ")
			buff.WriteString(fmt.Sprintf("//%s\n", comment.Value))
		}
		return true
	}
	return false
}

// writeFormatFile 格式化后的gs文件
func writeFormatFile(script *ast.Script, bytes []byte) {
	fullPath, ok := cblang.FilePath(script)
	if !ok {
		cberrors.Panic("compile must bind file path to script")
	}
	// 写入文件名为 源文件名+.gss
	// fullPath += ".gss"
	err := os.WriteFile(fullPath, bytes, 0644)
	if err != nil {
		cberrors.Panic(err.Error())
	}
	log.Infof("Format file successfully: %s success", fullPath)
}

// formatScript 格式化代码并输出到文件
func formatScript(script *ast.Script) {
	// format gs file
	var buff bytes.Buffer

	count := 0
	for _, ref := range script.Imports {
		if ref.Name() != "cblang" {
			count++
			break
		}
	}

	// format imports
	if count > 0 {
		buff.WriteString("import (\n")
		for _, ref := range script.Imports {
			if ref.Name() != "cblang" {
				if ref.Name() == filepath.Base(ref.Ref.Name()) {
					buff.WriteString(fmt.Sprintf("\t\"%s\"\n", ref.Ref))
				} else {
					buff.WriteString(fmt.Sprintf("\t%s \"%s\"\n", ref.Name(), ref.Ref))
				}
			}
		}
		buff.WriteString(")\n\n")
	}
	// format script comments
	if printComments(&buff, script) {
		buff.WriteString("\n")
	}

	// format enum
	for _, t := range script.Types {
		if enum, ok := t.(*ast.Enum); ok {
			printComments(&buff, enum)
			printAttrs(&buff, enum)
			buff.WriteString(fmt.Sprintf("enum %s {\n", enum.Name()))
			maxLen := enum.MaxKeyLength
			maxValueLen := enum.MaxValueLength + 2
			sortedValues := enum.SortedValues()
			for _, field := range sortedValues {
				tmp := "\t%" + fmt.Sprintf("-%d", maxLen) + "s = %" + fmt.Sprintf("-%d", maxValueLen) + "s"
				buff.WriteString(fmt.Sprintf(tmp, field.Name(), fmt.Sprintf("%d; ", field.Value)))
				printCommentsToLine(&buff, field)
				buff.WriteString("\n")
			}
			buff.WriteString(fmt.Sprintf("}\n\n"))
		}
	}
	// format struct
	for _, t := range script.Types {
		if table, ok := t.(*ast.Table); ok {
			printComments(&buff, table)
			printAttrs(&buff, table)
			if cblang.IsStruct(table) {
				buff.WriteString(fmt.Sprintf("struct %s {\n", table.Name()))
			} else {
				buff.WriteString(fmt.Sprintf("table %s {\n", table.Name()))
			}
			maxNameLen := table.MaxFieldNameLength
			maxTypeLen := table.MaxFieldTypeLength
			maxIDLen := table.MaxFieldIDLength + 2
			for _, field := range table.Fields {
				tmp := "\t%" +
					fmt.Sprintf("-%d", maxNameLen) +
					"s %" +
					fmt.Sprintf("-%d", maxTypeLen) +
					"s = %" +
					fmt.Sprintf("-%d", maxIDLen) +
					"s"
				buff.WriteString(fmt.Sprintf(tmp,
					field.Name(),
					field.Type.OriginName(),
					fmt.Sprintf("%d; ", field.ID)))
				printCommentsToLine(&buff, field)
				buff.WriteString("\n")
			}
			buff.WriteString(fmt.Sprintf("}\n\n"))
		}
	}
	// format service
	for _, t := range script.Types {
		if service, ok := t.(*ast.Service); ok {
			service.CalMethodLength()
			printComments(&buff, service)
			printAttrs(&buff, service)
			buff.WriteString(fmt.Sprintf("service %s {\n", service.OriginName()))
			for _, method := range service.MethodList {
				tmp := "\t%" +
					fmt.Sprintf("-%d", service.MaxMethodFirst) +
					"s %" +
					fmt.Sprintf("-%d", service.MaxMethodSecond) +
					"s"
				buff.WriteString(fmt.Sprintf(tmp, method.OriginFirst(), method.OriginSecond()))
				printCommentsToLine(&buff, method)
				buff.WriteString("\n")
			}
			buff.WriteString(fmt.Sprintf("}\n\n"))
		}
	}
	writeFormatFile(script, buff.Bytes())
}
