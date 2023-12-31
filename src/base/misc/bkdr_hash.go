// -------------------------------------------
// @file      : bkdr_hash.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/1 上午1:04
// -------------------------------------------

package misc

// BKDRHashBytes 计算字节数组哈希
func BKDRHashBytes(b []byte) uint32 {
	seed := uint32(131)
	hash := uint32(0)
	for _, v := range b {
		hash = hash*seed + uint32(v)
	}
	return hash
}

// BKDRHash 计算字符串哈希
func BKDRHash(s string) uint32 {
	seed := uint32(131)
	hash := uint32(0)
	for i := 0; i < len(s); i++ {
		hash = hash*seed + uint32(s[i])
	}
	return hash
}
