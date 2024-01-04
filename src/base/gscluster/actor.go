// -------------------------------------------
// @file      : actor.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午3:40
// -------------------------------------------

package gscluster

import (
	"fmt"
	"gogs/base/mongodb"
	"regexp"
	"strconv"
	"sync"
)

// IActorContext 角色上下文接口
// 角色上下文是角色的数据载体
type IActorContext interface {
	Actor() IActor                   // 获取所绑定角色
	SetActor(IActor)                 // 绑定角色到上下文
	Save(*mongodb.MongoClient) error // 角色上下文存盘
}

// NilActorContext 空角色上下文
type NilActorContext struct {
	actor IActor // 角色
}

// NewNilActorContext 创建空角色上下文
func NewNilActorContext() *NilActorContext {
	return &NilActorContext{}
}

// Actor 获取所绑定角色
func (ctx *NilActorContext) Actor() IActor {
	return ctx.actor
}

// SetActor 绑定角色到上下文
func (ctx *NilActorContext) SetActor(actor IActor) {
	ctx.actor = actor
}

// Save 角色上下文存盘
func (ctx *NilActorContext) Save(*mongodb.MongoClient) error {
	return nil
}

// ActorName 角色名字
type ActorName struct {
	SystemName string // 角色系统名字
	Type       string // 角色类型
	ID         int64  // 角色ID 同类型的角色下的唯一标识符
	actor      IActor // 绑定的角色
}

// NewActorName 新建角色名
func NewActorName(url string) (*ActorName, error) {
	regex := regexp.MustCompile(`(\w*):(\w*)@(\w*)`)
	match := regex.FindStringSubmatch(url)
	if nil == match {
		return nil, fmt.Errorf("actor name: %s regex match failed", url)
	}
	if len(match) != 4 {
		return nil, fmt.Errorf("actor name: %s regex match result len must be 4", url)
	}
	id, err := strconv.ParseInt(match[3], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("actor name: %s regex match result id must be int64", url)
	}
	return &ActorName{
		SystemName: match[1],
		Type:       match[2],
		ID:         id,
	}, nil
}

// String 角色名字封装显示
func (name *ActorName) String() string {
	return fmt.Sprintf("%s:%s@%d", name.SystemName, name.Type, name.ID)
}

// Actor 获取绑定的实际角色
func (name *ActorName) Actor() IActor {
	return name.actor
}

// SetActor 将名字绑定到角色
func (name *ActorName) SetActor(actor IActor) {
	name.actor = actor
}

// IActor 角色接口
// 一个角色包含了一个服务
// 角色可以看成是在服务的基础上加载了更多应用层数据的数据单元
type IActor interface {
	Name() string           // 角色名字,格式为:ActorSystemName:ActorType@ActorID
	ID() int64              // 角色ID
	Type() string           // 角色类型
	System() *ActorSystem   // 角色系统
	Service() IService      // 角色服务
	Context() IActorContext // 角色上下文
	Locker() *sync.Mutex    // 在分组锁中根据ID获取锁
	Lock()                  // 加锁
	Unlock()                // 解锁
}

// baseActor 基础Actor
type baseActor struct {
	locker  *sync.Mutex   // 分组锁
	name    *ActorName    // 角色名字
	system  *ActorSystem  // 角色所属系统
	service IService      // 角色挂载的服务
	context IActorContext // 角色上下文
}

// newBaseActor 创建基本角色
func newBaseActor(system *ActorSystem, name *ActorName, locker *sync.Mutex, context IActorContext) *baseActor {
	return &baseActor{
		system:  system,
		locker:  locker,
		context: context,
		name:    name,
	}
}

// Name 获取角色名字
func (actor *baseActor) Name() string {
	return actor.name.String()
}

// ID 获取角色唯一ID
func (actor *baseActor) ID() int64 {
	return actor.name.ID
}

// Type 获取角色类型
func (actor *baseActor) Type() string {
	return actor.name.Type
}

// System 获取角色系统
func (actor *baseActor) System() *ActorSystem {
	return actor.system
}

// Service 获取角色服务
func (actor *baseActor) Service() IService {
	return actor.service
}

// Context 获取角色上下文
func (actor *baseActor) Context() IActorContext {
	return actor.context
}

// Locker 获取分组锁
func (actor *baseActor) Locker() *sync.Mutex {
	return actor.locker
}

// Lock 加锁
func (actor *baseActor) Lock() {
	actor.locker.Lock()
}

// Unlock 解锁
func (actor *baseActor) Unlock() {
	actor.locker.Unlock()
}
