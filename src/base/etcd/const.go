package etcd

import "time"

type Event int32

const (
	EventAdd Event = iota
	EventDelete
	EventUpdate
)

func (e Event) String() string {
	switch e {
	case EventAdd:
		return "EventAdd"
	case EventDelete:
		return "EventDelete"
	case EventUpdate:
		return "EventUpdate"
	default:
		return "Unknown"
	}
}

const ServerIDDefault = -1

const (
	ServerTypeDefault     = "DEFAULT"
	ServerTypeGame        = "GAME"
	ServerTypeGameConfig  = "GAME_CONFIG"
	ServerTypeLogin       = "LOGIN"
	ServerTypeLoginConfig = "LOGIN_CONFIG"
	ServerTypeGate        = "GATE"
	ServerTypeMongo       = "MONGO"
)

const (
	NodeInfoKeyID         = "ID"
	NodeInfoKeyType       = "Type"
	NodeInfoKeyPID        = "PID"
	NodeInfoKeyHostname   = "Hostname"
	NodeInfoKeyAddrList   = "AddrList"
	NodeInfoKeyUser       = "User"
	NodeInfoKeyPassword   = "Password"
	NodeInfoKeyURL        = "URL"
	NodeInfoKeyMongoID    = "MongoID"
	NodeInfoKeyReplicaSet = "ReplicaSet"
)

const (
	NodeInfoKeyOpenTime     = "OpenTime" // GS
	NodeInfoKeyStatus       = "Status"
	NodeInfoKeyCountry      = "Country"
	NodeInfoKeyMaxRegister  = "MaxRegister"
	NodeInfoKeyCurRegister  = "CurRegister"
	NodeInfoKeyOpenRegister = "OpenRegister"
	NodeInfoKeyHide         = "Hide"
	NodeInfoKeyMaintain     = "Maintain"
	NodeInfoKeyCurOnline    = "CurOnline"
	NodeInfoKeyIsBlockLogin = "IsBlockLogin" // Login
)

const (
	StateDisconnect = iota
	StateConnected
	StateClosed
)

const (
	Timeout           = 5 * time.Second
	ReconnectDuration = 5 * time.Second
	DefaultTTL        = 10
)

type NodeEvent struct {
	Node  NodeInfo
	Event Event
}
