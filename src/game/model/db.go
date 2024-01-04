// -------------------------------------------
// @file      : db.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午11:27
// -------------------------------------------

package model

import (
	"gogs/base/config"
	"gogs/base/etcd"
	"gogs/base/gserrors"
	log "gogs/base/logger"
	"gogs/base/mongodb"
)

const (
	UserCollection = "users"
)

var (
	mongoClient = mongodb.NewMongoClient()
)

func InitMongoDB(serverID int64) {
	gameConfigNode, err := etcd.GetDepByTypeAndID(etcd.ServerTypeGameConfig, serverID)
	if err != nil {
		gserrors.Panicf("unable to find game_config in etcd. serverID:%d. err:%s", serverID, err)
	}
	if gameConfigNode == nil {
		gserrors.Panicf("unable to find game_config in etcd. serverID:%d", serverID)
	}
	mongoID := gameConfigNode.GetMongoID()
	mongoNode, err := etcd.GetDepByTypeAndID(etcd.ServerTypeMongo, mongoID)
	if err != nil {
		gserrors.Panicf("unable to find mongo in etcd. mongoID:%d. err:%s", mongoID, err)
	}
	if mongoNode == nil {
		gserrors.Panicf("unable to find mongo in etcd. mongoID:%d", mongoID)
	}
	dbName := config.GetGameConfig().DBName
	url := mongoNode.GetMongoConnectURL()
	log.Infof("start connecting to mongodb: %s", url)
	err = mongoClient.Connect(url, dbName)
	if err != nil {
		gserrors.Panicf("mongoClient connect err:%s", err)
	}
	err = mongoClient.CreateIndex(UserCollection, "id", true)
	if err != nil {
		gserrors.Panicf("mongoClient CreateIndex err:%s", err)
	}
	err = mongoClient.CreateIndex(UserCollection, "username", false)
	if err != nil {
		gserrors.Panicf("mongoClient CreateIndex err:%s", err)
	}
	log.Infof("mongoClient: %s connected", url)
}

func CloseMongoDB() {
	err := mongoClient.Disconnect()
	if err != nil {
		log.Errorf("mongo client disconnect err:%s", err)
	} else {
		log.Info("mongo client disconnected")
	}
}

func MongoClient() *mongodb.MongoClient {
	return mongoClient
}
