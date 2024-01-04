package etcd

import (
	"encoding/json"
	"fmt"
	log "gogs/base/logger"
	"reflect"
	"strconv"
	"strings"
)

// NodeInfo 节点信息
type NodeInfo map[string]interface{}

func (n NodeInfo) GetID() int64 {
	ret, _ := n.GetInt64(NodeInfoKeyID)
	return ret
}

func (n NodeInfo) GetType() string {
	return n.GetString(NodeInfoKeyType)
}

func (n NodeInfo) GetUser() string {
	return n.GetString(NodeInfoKeyUser)
}

func (n NodeInfo) GetPassword() string {
	return n.GetString(NodeInfoKeyPassword)
}

func (n NodeInfo) GetReplicaSet() string {
	return n.GetString(NodeInfoKeyReplicaSet)
}

func (n NodeInfo) GetAddrList() []string {
	return n.GetStringSlice(NodeInfoKeyAddrList)
}

func (n NodeInfo) GetMongoID() int64 {
	ret, _ := n.GetInt64(NodeInfoKeyMongoID)
	return ret
}

func (n NodeInfo) GetMongoConnectURL() string {
	url := n.GetString(NodeInfoKeyURL)
	if url != "" {
		return url
	}
	addrList := n.GetAddrList()
	if len(addrList) <= 0 {
		return ""
	}
	hosts := strings.Join(addrList, ",")
	user := n.GetUser()
	password := n.GetPassword()
	replicaSet := n.GetReplicaSet()
	if user != "" {
		hosts = user + ":" + password + "@" + hosts
	}
	url = fmt.Sprintf("mongodb://%s/", hosts)
	if replicaSet != "" {
		url += "?replicaSet=" + replicaSet
	} else {
		url += "?connect=direct"
	}
	return url
}

func (n NodeInfo) GetCommonConnectURL() string {
	addrList := n.GetAddrList()
	if len(addrList) <= 0 {
		return ""
	}
	return addrList[0]
}

func (n NodeInfo) GetConnectURL() string {
	switch n.GetType() {
	case ServerTypeMongo:
		return n.GetMongoConnectURL()
	}
	return n.GetCommonConnectURL()
}

func (n NodeInfo) GetInt32(key string, defaultValue int32) int32 {
	value, ok := n.GetInt64(key)
	if !ok {
		return defaultValue
	}
	return int32(value)
}

func (n NodeInfo) GetInt64(key string) (int64, bool) {
	var ret int64
	value, ok := n[key]
	if !ok {
		return 0, false
	}
	switch v := value.(type) {
	case int:
		ret = int64(v)
	case int8:
		ret = int64(v)
	case int16:
		ret = int64(v)
	case int32:
		ret = int64(v)
	case int64:
		ret = v
	case uint:
		ret = int64(v)
	case uint8:
		ret = int64(v)
	case uint16:
		ret = int64(v)
	case uint32:
		ret = int64(v)
	case uint64:
		ret = int64(v)
	case float32:
		ret = int64(v)
	case float64:
		ret = int64(v)
	case string:
		ret, _ = strconv.ParseInt(v, 10, 64)
	default:
		return 0, false
	}
	return ret, true
}

func (n NodeInfo) GetString(key string) string {
	if _, ok := n[key]; !ok {
		return ""
	}
	ret, ok := n[key].(string)
	if !ok {
		return ""
	}
	return ret
}

func (n NodeInfo) GetStringSlice(key string) []string {
	if _, ok := n[key]; !ok {
		return nil
	}
	list, ok := n[key].([]interface{})
	if !ok {
		return nil
	}
	ret := make([]string, len(list))
	for index, info := range list {
		if str, ok := info.(string); ok {
			ret[index] = str
		}
	}
	return ret
}

func (n NodeInfo) GetBool(key string, defaultValue bool) bool {
	value, exists := n[key]
	if !exists {
		return defaultValue
	}
	switch v := value.(type) {
	case int8, int16, int32, int, int64, uint8, uint16, uint32, uint, uint64, float32, float64:
		num, _ := n.GetInt64(key)
		return num != 0
	case string:
		ret, _ := strconv.ParseBool(v)
		return ret
	case bool:
		return v
	default:
		return defaultValue
	}
}

func (n NodeInfo) IsInvalid() bool {
	return n.GetID() <= 0 || n.GetType() == ""
}

func (n NodeInfo) Clone() NodeInfo {
	ret := make(map[string]interface{})
	for k, v := range n {
		ret[k] = v
	}
	return ret
}

func (n NodeInfo) Copy() NodeInfo {
	ret := make(map[string]interface{})
	for key, val := range n {
		if list, ok := val.([]interface{}); ok {
			newList := make([]interface{}, len(list))
			copy(newList, list)
			ret[key] = newList
		} else {
			ret[key] = val
		}
	}
	return ret
}

func NodeDecode(value string) NodeInfo {
	node := NodeInfo{}
	err := json.Unmarshal([]byte(value), &node)
	if err != nil {
		log.Errorf("[ETCD] NodeDecode Unmarshal err: %v", err)
		return nil
	}
	return node
}

func IsNodeEqual(n1, n2 NodeInfo) bool {
	return reflect.DeepEqual(n1, n2)
}
