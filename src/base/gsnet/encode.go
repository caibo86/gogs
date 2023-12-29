// -------------------------------------------
// @file      : encode.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/19 下午6:22
// -------------------------------------------

package gsnet

import (
	"gogs/base/gserrors"
	"math"
)

// WriteFieldID 写入字段编号
func WriteFieldID(data []byte, i int, v uint16) int {
	data[i] = byte(v)
	data[i+1] = byte(v >> 8)
	return i + 2
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

// WriteInt8 写入一个有符号8位整数
func WriteInt8(data []byte, i int, v int8) int {
	data[i] = byte(v)
	return i + 1
}

// WriteUint8 写入一个无符号8位整数
func WriteUint8(data []byte, i int, v uint8) int {
	data[i] = v
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
	data[i] = byte(v)
	data[i+1] = byte(v >> 8)
	data[i+2] = byte(v >> 16)
	data[i+3] = byte(v >> 24)
	return i + 4
}

// ReadFieldID 读取字段编号
func ReadFieldID(data []byte, i int) (int, uint16) {
	return i + 2, uint16(data[i]) | uint16(data[i+1])<<8
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

// ReadInt8 读取一个有符号8位整数
func ReadInt8(data []byte, i int) (int, int8) {
	return i + 1, int8(data[i])
}

// ReadUint8 读取一个无符号8位整数
func ReadUint8(data []byte, i int) (int, uint8) {
	return i + 1, data[i]
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
	return i + 4, int32(data[i]) | int32(data[i+1])<<8 |
		int32(data[i+2])<<16 | int32(data[i+3])<<24
}

//////////////////////////////////////////////////////////////////////////////////////////////////

// MarshalBool 序列化一个布尔值
func MarshalBool(v bool) []byte {
	data := make([]byte, 1)
	if v == true {
		data[0] = 1
	} else {
		data[0] = 0
	}
	return data
}

// MarshalByte 序列化一个字节
func MarshalByte(v byte) []byte {
	data := make([]byte, 1)
	data[0] = v
	return data
}

// MarshalInt8 序列化一个有符号8位整数
func MarshalInt8(v int8) []byte {
	data := make([]byte, 1)
	data[0] = byte(v)
	return data
}

// MarshalUint8 序列化一个无符号8位整数
func MarshalUint8(v uint8) []byte {
	data := make([]byte, 1)
	data[0] = v
	return data
}

// MarshalInt16 序列化一个有符号16位整数
func MarshalInt16(v int16) []byte {
	data := make([]byte, 2)
	data[0] = byte(v)
	data[1] = byte(v >> 8)
	return data
}

// MarshalUint16 序列化一个无符号16位整数
func MarshalUint16(v uint16) []byte {
	data := make([]byte, 2)
	data[0] = byte(v)
	data[1] = byte(v >> 8)
	return data
}

// MarshalInt32 序列化一个有符号32位整数
func MarshalInt32(v int32) []byte {
	data := make([]byte, 4)
	data[0] = byte(v)
	data[1] = byte(v >> 8)
	data[2] = byte(v >> 16)
	data[3] = byte(v >> 24)
	return data
}

// MarshalUint32 序列化一个无符号32位整数
func MarshalUint32(v uint32) []byte {
	data := make([]byte, 4)
	data[0] = byte(v)
	data[1] = byte(v >> 8)
	data[2] = byte(v >> 16)
	data[3] = byte(v >> 24)
	return data
}

// MarshalInt64 序列化一个有符号64位整数
func MarshalInt64(v int64) []byte {
	data := make([]byte, 8)
	data[0] = byte(v)
	data[1] = byte(v >> 8)
	data[2] = byte(v >> 16)
	data[3] = byte(v >> 24)
	data[4] = byte(v >> 32)
	data[5] = byte(v >> 40)
	data[6] = byte(v >> 48)
	data[7] = byte(v >> 56)
	return data
}

// MarshalUint64 序列化一个无符号64位整数
func MarshalUint64(v uint64) []byte {
	data := make([]byte, 8)
	data[0] = byte(v)
	data[1] = byte(v >> 8)
	data[2] = byte(v >> 16)
	data[3] = byte(v >> 24)
	data[4] = byte(v >> 32)
	data[5] = byte(v >> 40)
	data[6] = byte(v >> 48)
	data[7] = byte(v >> 56)
	return data
}

// MarshalFloat32 序列化一个32位浮点数
func MarshalFloat32(v float32) []byte {
	return MarshalUint32(math.Float32bits(v))
}

// MarshalFloat64 序列化一个64位浮点数
func MarshalFloat64(v float64) []byte {
	return MarshalUint64(math.Float64bits(v))
}

// MarshalString 序列化一个字符串
func MarshalString(v string) []byte {
	return []byte(v)
}

// MarshalBytes 序列化一个字节流
func MarshalBytes(v []byte) []byte {
	data := make([]byte, len(v))
	copy(data, v)
	return data
}

// UnmarshalBool 反序列化一个布尔值
func UnmarshalBool(data []byte) (bool, error) {
	if len(data) != 1 {
		return false, gserrors.Newf("unmarshal bool, data length is not 1")
	}
	if data[0] == 0 {
		return false, nil
	}
	return true, nil
}

// UnmarshalByte 反序列化一个字节
func UnmarshalByte(data []byte) (byte, error) {
	if len(data) != 1 {
		return 0, gserrors.Newf("unmarshal byte, data length is not 1")
	}
	return data[0], nil
}

// UnmarshalInt8 反序列化一个有符号8位整数
func UnmarshalInt8(data []byte) (int8, error) {
	if len(data) != 1 {
		return 0, gserrors.Newf("unmarshal int8, data length is not 1")
	}
	return int8(data[0]), nil
}

// UnmarshalUint8 反序列化一个无符号8位整数
func UnmarshalUint8(data []byte) (uint8, error) {
	if len(data) != 1 {
		return 0, gserrors.Newf("unmarshal uint8, data length is not 1")
	}
	return data[0], nil
}

// UnmarshalInt16 反序列化一个有符号16位整数
func UnmarshalInt16(data []byte) (int16, error) {
	if len(data) != 2 {
		return 0, gserrors.Newf("unmarshal int16, data length is not 2")
	}
	return int16(data[0]) | int16(data[1])<<8, nil
}

// UnmarshalUint16 反序列化一个无符号16位整数
func UnmarshalUint16(data []byte) (uint16, error) {
	if len(data) != 2 {
		return 0, gserrors.Newf("unmarshal uint16, data length is not 2")
	}
	return uint16(data[0]) | uint16(data[1])<<8, nil
}

// UnmarshalInt32 反序列化一个有符号32位整数
func UnmarshalInt32(data []byte) (int32, error) {
	if len(data) != 4 {
		return 0, gserrors.Newf("unmarshal int32, data length is not 4")
	}
	return int32(data[0]) | int32(data[1])<<8 |
		int32(data[2])<<16 | int32(data[3])<<24, nil
}

// UnmarshalUint32 反序列化一个无符号32位整数
func UnmarshalUint32(data []byte) (uint32, error) {
	if len(data) != 4 {
		return 0, gserrors.Newf("unmarshal uint32, data length is not 4")
	}
	return uint32(data[0]) | uint32(data[1])<<8 |
		uint32(data[2])<<16 | uint32(data[3])<<24, nil
}

// UnmarshalInt64 反序列化一个有符号64位整数
func UnmarshalInt64(data []byte) (int64, error) {
	if len(data) != 8 {
		return 0, gserrors.Newf("unmarshal int64, data length is not 8")
	}
	return int64(data[0]) | int64(data[1])<<8 |
		int64(data[2])<<16 | int64(data[3])<<24 |
		int64(data[4])<<32 | int64(data[5])<<40 |
		int64(data[6])<<48 | int64(data[7])<<56, nil
}

// UnmarshalUint64 反序列化一个无符号64位整数
func UnmarshalUint64(data []byte) (uint64, error) {
	if len(data) != 8 {
		return 0, gserrors.Newf("unmarshal uint64, data length is not 8")
	}
	return uint64(data[0]) | uint64(data[1])<<8 |
		uint64(data[2])<<16 | uint64(data[3])<<24 |
		uint64(data[4])<<32 | uint64(data[5])<<40 |
		uint64(data[6])<<48 | uint64(data[7])<<56, nil
}

// UnmarshalString 反序列化一个字符串
func UnmarshalString(data []byte) (string, error) {
	return string(data), nil
}

// UnmarshalBytes 反序列化一个字节流
func UnmarshalBytes(data []byte) ([]byte, error) {
	bytes := make([]byte, len(data))
	copy(bytes, data)
	return bytes, nil
}
