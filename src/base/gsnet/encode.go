// -------------------------------------------
// @file      : encode.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/19 下午6:22
// -------------------------------------------

package gsnet

import (
	"errors"
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

var (
	// ErrEncode 编码错误
	ErrEncode = errors.New("encode error")
	// ErrDecode 解码错误
	ErrDecode = errors.New("decode error")
	// ErrWriteNone 写入内容空错误
	ErrWriteNone = errors.New("write nothing")
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

// ReadTag 读取一个数据类型标签
func ReadTag(reader Reader) (Tag, error) {
	v, err := reader.ReadByte()
	return Tag(v), err
}

// ReadByte 读取一个字节
func ReadByte(reader Reader) (byte, error) {
	return reader.ReadByte()
}

// ReadSbyte 读取一个字节 并转换为int8
func ReadSbyte(reader Reader) (int8, error) {
	v, err := reader.ReadByte()
	return int8(v), err
}

// ReadUint16 读取一个无符号16位整数 每8位转成一个byte传送 小端先行
func ReadUint16(reader Reader) (uint16, error) {
	buf := make([]byte, 2)
	_, err := reader.Read(buf)
	if err != nil {
		return 0, err
	}
	return uint16(buf[0]) | uint16(buf[1])<<8, nil
}

// ReadInt16 读取一个有符号16位整数
func ReadInt16(reader Reader) (int16, error) {
	v, err := ReadUint16(reader)
	return int16(v), err
}

// ReadUint32 读取一个无符号32位整数
func ReadUint32(reader Reader) (uint32, error) {
	buf := make([]byte, 4)
	_, err := reader.Read(buf)
	if err != nil {
		return 0, err
	}
	return uint32(buf[3])<<24 | uint32(buf[2])<<16 | uint32(buf[1])<<8 | uint32(buf[0]), nil
}

// ReadInt32 读取一个有符号32位整数
func ReadInt32(reader Reader) (int32, error) {
	v, err := ReadUint32(reader)
	return int32(v), err
}

// ReadUint64 读取一个无符号64位整数
func ReadUint64(reader Reader) (uint64, error) {
	buf := make([]byte, 8)
	_, err := reader.Read(buf)
	if err != nil {
		return 0, err
	}
	var ret uint64
	for i, v := range buf {
		ret |= uint64(v) << uint(i*8)
	}
	return ret, nil
}

// ReadInt64 读取一个有符号64位整数
func ReadInt64(reader Reader) (int64, error) {
	v, err := ReadUint64(reader)
	return int64(v), err
}

// ReadFloat32 读取一个32位浮点数
func ReadFloat32(reader Reader) (float32, error) {
	v, err := ReadUint32(reader)
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(v), nil
}

// ReadFloat64 读取一个64位浮点数
func ReadFloat64(reader Reader) (float64, error) {
	v, err := ReadUint64(reader)
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(v), nil
}

// ReadString 读取一个字符串 前2个字节uint16表示字符串长度
func ReadString(reader Reader) (string, error) {
	length, err := ReadUint16(reader)
	if err != nil {
		return "", err
	}
	buf := make([]byte, length)
	_, err = reader.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// ReadBytes 读取一个字节流到指定的buf
func ReadBytes(reader Reader, buf []byte) error {
	_, err := reader.Read(buf)
	return err
}

// ReadBool 读取一个布尔值
func ReadBool(reader Reader) (bool, error) {
	v, err := ReadByte(reader)
	if v == 1 {
		return true, err
	}
	return false, err
}

// WriteByte 写入一个字节
func WriteByte(writer Writer, v byte) error {
	return writer.WriteByte(v)
}

// WriteTag 写入一个数据类型标签
func WriteTag(writer Writer, v Tag) error {
	return writer.WriteByte(byte(v))
}

// WriteBool 写入一个布尔值
func WriteBool(writer Writer, v bool) error {
	if v == true {
		return writer.WriteByte(1)
	}
	return writer.WriteByte(0)
}

// WriteSbyte 写入一个有符号单字节
func WriteSbyte(writer Writer, v int8) error {
	return writer.WriteByte(byte(v))
}

// WriteUint16 写入一个无符号16位整数
func WriteUint16(writer Writer, v uint16) error {
	if err := WriteByte(writer, byte(v)); err != nil {
		return err
	}
	if err := WriteByte(writer, byte(v>>8)); err != nil {
		return err
	}
	return nil
}

// WriteInt16 写入一个有符号16位整数
func WriteInt16(writer Writer, v int16) error {
	for i := uint(0); i < 2; i++ {
		if err := WriteByte(writer, byte(v>>(i*8))); err != nil {
			return err
		}
	}
	return nil
}

// WriteUint32 写入一个无符号32位整数
func WriteUint32(writer Writer, v uint32) error {
	for i := uint(0); i < 4; i++ {
		if err := WriteByte(writer, byte(v>>(i*8))); err != nil {
			return err
		}
	}
	return nil
}

// WriteInt32 写入一个有符号32位整数
func WriteInt32(writer Writer, v int32) error {
	for i := uint(0); i < 4; i++ {
		if err := WriteByte(writer, byte(v>>(i*8))); err != nil {
			return err
		}
	}
	return nil
}

// WriteUint64 写入一个无符号64位整数
func WriteUint64(writer Writer, v uint64) error {
	for i := uint(0); i < 8; i++ {
		if err := WriteByte(writer, byte(v>>(i*8))); err != nil {
			return err
		}
	}
	return nil
}

// WriteInt64 写入一个有符号64位整数
func WriteInt64(writer Writer, v int64) error {
	for i := uint(0); i < 8; i++ {
		if err := WriteByte(writer, byte(v>>(i*8))); err != nil {
			return err
		}
	}
	return nil
}

// WriteFloat32 写入一个32位浮点数
func WriteFloat32(writer Writer, v float32) error {
	return WriteUint32(writer, math.Float32bits(v))
}

// WriteFloat64 写入一个64位浮点数
func WriteFloat64(writer Writer, v float64) error {
	return WriteUint64(writer, math.Float64bits(v))
}

// WriteString 写入一个字符串
func WriteString(writer Writer, v string) error {
	if err := WriteUint16(writer, uint16(len(v))); err != nil {
		return err
	}
	_, err := writer.Write([]byte(v))
	return err
}

// WriteBytes 写入一个字节流
func WriteBytes(writer Writer, bytes []byte) error {
	_, err := writer.Write(bytes)
	return err
}

// WriteTagByte 写入一个字节+前置标签
func WriteTagByte(writer Writer, v byte) error {
	err := WriteTag(writer, Byte)
	if err != nil {
		return err
	}
	return writer.WriteByte(v)
}

// WriteTagBool 写入一个布尔值+前置标签
func WriteTagBool(writer Writer, v bool) error {
	err := WriteTag(writer, Bool)
	if err != nil {
		return err
	}
	if v == true {
		return writer.WriteByte(1)
	}
	return writer.WriteByte(0)
}

// WriteTagSByte 写入一个有符号单字节+前置标签
func WriteTagSByte(writer Writer, v int8) error {
	err := WriteTag(writer, SByte)
	if err != nil {
		return err
	}
	return writer.WriteByte(byte(v))
}

// WriteTagUint16 写入一个无符号16位整数+前置标签
func WriteTagUint16(writer Writer, v uint16) error {
	err := WriteTag(writer, Uint16)
	if err != nil {
		return err
	}
	if err = WriteByte(writer, byte(v)); err != nil {
		return err
	}
	if err = WriteByte(writer, byte(v>>8)); err != nil {
		return err
	}
	return nil
}

// WriteTagInt16 写入一个有符号16位整数+前置标签
func WriteTagInt16(writer Writer, v int16) error {
	err := WriteTag(writer, Int16)
	if err != nil {
		return err
	}
	for i := uint(0); i < 2; i++ {
		if err = WriteByte(writer, byte(v>>i*8)); err != nil {
			return err
		}
	}
	return nil
}

// WriteTagUint32 写入一个无符号32位整数+前置标签
func WriteTagUint32(writer Writer, v uint32) error {
	err := WriteTag(writer, Uint32)
	if err != nil {
		return err
	}
	for i := uint(0); i < 4; i++ {
		if err = WriteByte(writer, byte(v>>i*8)); err != nil {
			return err
		}
	}
	return nil
}

// WriteTagInt32 写入一个有符号32位整数+前置标签
func WriteTagInt32(writer Writer, v int32) error {
	err := WriteTag(writer, Int32)
	if err != nil {
		return err
	}
	for i := uint(0); i < 4; i++ {
		if err = WriteByte(writer, byte(v>>i*8)); err != nil {
			return err
		}
	}
	return nil
}

// WriteTagUint64 写入一个无符号64位整数+前置标签
func WriteTagUint64(writer Writer, v uint64) error {
	err := WriteTag(writer, Uint64)
	if err != nil {
		return err
	}
	for i := uint(0); i < 4; i++ {
		if err = WriteByte(writer, byte(v>>i*8)); err != nil {
			return err
		}
	}
	return nil
}

// WriteTagInt64 写入一个有符号64位整数+前置标签
func WriteTagInt64(writer Writer, v int64) error {
	err := WriteTag(writer, Int64)
	if err != nil {
		return err
	}
	for i := uint(0); i < 8; i++ {
		if err := WriteByte(writer, byte(v>>i*8)); err != nil {
			return err
		}
	}
	return nil
}

// WriteTagFloat32 写入一个32位浮点数+前置标签
func WriteTagFloat32(writer Writer, v float32) error {
	err := WriteTag(writer, Float32)
	if err != nil {
		return err
	}
	return WriteUint32(writer, math.Float32bits(v))
}

// WriteTagFloat64 写入一个64位浮点数+前置标签
func WriteTagFloat64(writer Writer, v float64) error {
	err := WriteTag(writer, Float64)
	if err != nil {
		return err
	}
	return WriteUint64(writer, math.Float64bits(v))
}

// WriteTagString 写入一个字符串+前置标签
func WriteTagString(writer Writer, v string) error {
	err := WriteTag(writer, String)
	if err != nil {
		return err
	}
	if err = WriteUint16(writer, uint16(len(v))); err != nil {
		return err
	}
	_, err = writer.Write([]byte(v))
	return err
}
