// -------------------------------------------
// @file      : encode.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/19 下午6:22
// -------------------------------------------

package gsnet

import (
	"io"
	"math"
)

// Tag 数值类型标签
type Tag byte

const (
	// None 空
	None Tag = iota
	// Byte 单字节
	Byte
	// SByte int8
	SByte
	// Uint16 16位无符号整数
	Uint16
	// Int16 16位有符号整数
	Int16
	// Uint32 32位无符号整数
	Uint32
	// Int32 32位有符号整数
	Int32
	// Uint64 64位无符号整数
	Uint64
	// Int64 64位有符号整数
	Int64
	// Float32 32位浮点数
	Float32
	// Float64 64位浮点数
	Float64
	// Enum 枚举
	Enum
	// Struct 结构
	Struct
	// Table 表
	Table
	// Array 数组
	Array
	// List 链表
	List
	// String 字符串
	String
	// Bool 布尔值
	Bool
)

// Reader 读取器
type Reader interface {
	io.ByteReader
	io.Reader
}

// Writer 写入器
type Writer interface {
	io.ByteWriter
	io.Writer
}

// WriteFieldNum 写入字段编号
func WriteFieldNum(data []byte, i int, fieldNum uint8) int {
	data[i] = fieldNum
	return i + 1
}

// WriteBool 写入一个布尔值
func WriteBool(data []byte, i int, v bool) int {
	if v == true {
		data[i] = 1
	} else {
		data[i] = 0
	}
	return i + 1
}

// WriteByte 写入一个字节
func WriteByte(data []byte, i int, v byte) int {
	data[i] = v
	return i + 1
}

// WriteSbyte 写入一个有符号字节
func WriteSbyte(data []byte, i int, v int8) int {
	data[i] = byte(v)
	return i + 1
}

// WriteUint16 写入一个无符号16位整数
func WriteUint16(data []byte, i int, v uint16) int {
	data[i] = byte(v)
	data[i+1] = byte(v >> 8)
	return i + 2
}

// WriteInt16 写入一个有符号16位整数
func WriteInt16(data []byte, i int, v int16) int {
	data[i] = byte(v)
	data[i+1] = byte(v >> 8)
	return i + 2
}

// WriteUint32 写入一个无符号32位整数
func WriteUint32(data []byte, i int, v uint32) int {
	data[i] = byte(v)
	data[i+1] = byte(v >> 8)
	data[i+2] = byte(v >> 16)
	data[i+3] = byte(v >> 24)
	return i + 4
}

// WriteInt32 写入一个有符号32位整数
func WriteInt32(data []byte, i int, v int32) int {
	data[i] = byte(v)
	data[i+1] = byte(v >> 8)
	data[i+2] = byte(v >> 16)
	data[i+3] = byte(v >> 24)
	return i + 4
}

// WriteUint64 写入一个无符号64位整数
func WriteUint64(data []byte, i int, v uint64) int {
	data[i] = byte(v)
	data[i+1] = byte(v >> 8)
	data[i+2] = byte(v >> 16)
	data[i+3] = byte(v >> 24)
	data[i+4] = byte(v >> 32)
	data[i+5] = byte(v >> 40)
	data[i+6] = byte(v >> 48)
	data[i+7] = byte(v >> 56)
	return i + 8
}

// WriteInt64 写入一个有符号64位整数
func WriteInt64(data []byte, i int, v int64) int {
	data[i] = byte(v)
	data[i+1] = byte(v >> 8)
	data[i+2] = byte(v >> 16)
	data[i+3] = byte(v >> 24)
	data[i+4] = byte(v >> 32)
	data[i+5] = byte(v >> 40)
	data[i+6] = byte(v >> 48)
	data[i+7] = byte(v >> 56)
	return i + 8
}

// WriteFloat32 写入一个32位浮点数
func WriteFloat32(data []byte, i int, v float32) int {
	return WriteUint32(data, i, math.Float32bits(v))
}

// WriteFloat64 写入一个64位浮点数
func WriteFloat64(data []byte, i int, v float64) int {
	return WriteUint64(data, i, math.Float64bits(v))
}

// WriteString 写入一个字符串
func WriteString(data []byte, i int, v string) int {
	i = WriteUint32(data, i, uint32(len(v)))
	copy(data[i:], v)
	return i + len(v)
}

// WriteBytes 写入一个字节流
func WriteBytes(data []byte, i int, bytes []byte) int {
	l := len(bytes)
	i = WriteUint32(data, i, uint32(l))
	copy(data[i:], bytes)
	return i + len(bytes)
}

// WriteEnum 写入一个枚举
func WriteEnum(data []byte, i int, v int32) int {
	return WriteInt32(data, i, v)
}

// ReadFieldNum 读取字段编号
func ReadFieldNum(data []byte, i int) (int, uint8) {
	return i + 1, data[i]
}

// ReadBool 读取一个布尔值
func ReadBool(data []byte, i int) (int, bool) {
	if data[i] == 0 {
		return i + 1, false
	}
	return i + 1, true
}

// ReadByte 读取一个字节
func ReadByte(data []byte, i int) (int, byte) {
	return i + 1, data[i]
}

// ReadSbyte 读取一个有符号字节
func ReadSbyte(data []byte, i int) (int, int8) {
	return i + 1, int8(data[i])
}

// ReadUint16 读取一个无符号16位整数
func ReadUint16(data []byte, i int) (int, uint16) {
	return i + 2, uint16(data[i]) | uint16(data[i+1])<<8
}

// ReadInt16 读取一个有符号16位整数
func ReadInt16(data []byte, i int) (int, int16) {
	return i + 2, int16(data[i]) | int16(data[i+1])<<8
}

// ReadUint32 读取一个无符号32位整数
func ReadUint32(data []byte, i int) (int, uint32) {
	return i + 4, uint32(data[i]) | uint32(data[i+1])<<8 |
		uint32(data[i+2])<<16 | uint32(data[i+3])<<24
}

// ReadInt32 读取一个有符号32位整数
func ReadInt32(data []byte, i int) (int, int32) {
	return i + 4, int32(data[i]) | int32(data[i+1])<<8 |
		int32(data[i+2])<<16 | int32(data[i+3])<<24
}

// ReadUint64 读取一个无符号64位整数
func ReadUint64(data []byte, i int) (int, uint64) {
	return i + 8, uint64(data[i]) | uint64(data[i+1])<<8 |
		uint64(data[i+2])<<16 | uint64(data[i+3])<<24 |
		uint64(data[i+4])<<32 | uint64(data[i+5])<<40 |
		uint64(data[i+6])<<48 | uint64(data[i+7])<<56
}

// ReadInt64 读取一个有符号64位整数
func ReadInt64(data []byte, i int) (int, int64) {
	return i + 8, int64(data[i]) | int64(data[i+1])<<8 |
		int64(data[i+2])<<16 | int64(data[i+3])<<24 |
		int64(data[i+4])<<32 | int64(data[i+5])<<40 |
		int64(data[i+6])<<48 | int64(data[i+7])<<56
}

// ReadFloat32 读取一个32位浮点数
func ReadFloat32(data []byte, i int) (int, float32) {
	i, v := ReadUint32(data, i)
	return i, math.Float32frombits(v)
}

// ReadFloat64 读取一个64位浮点数
func ReadFloat64(data []byte, i int) (int, float64) {
	i, v := ReadUint64(data, i)
	return i, math.Float64frombits(v)
}

// ReadString 读取一个字符串
func ReadString(data []byte, i int) (int, string) {
	i, length := ReadUint32(data, i)
	return i + int(length), string(data[i : i+int(length)])
}

// ReadBytes 读取一个字节流
func ReadBytes(data []byte, i int) (int, []byte) {
	var l uint32
	i, l = ReadUint32(data, i)
	bytes := make([]byte, int(l))
	copy(bytes, data[i:])
	return i + int(l), bytes
}

// ReadEnum 读取一个枚举
func ReadEnum(data []byte, i int) (int, int32) {
	return ReadInt32(data, i)
}
