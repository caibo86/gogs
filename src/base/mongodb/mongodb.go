package mongodb

import (
	"bytes"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	log "gogs/base/logger"
	"time"
)

var (
	Timeout                = 15 * time.Second
	ErrMongoDBNotConnected = fmt.Errorf("mongodb not connected")
	ErrMongoDBNotAsync     = fmt.Errorf("cur mongodb doesn't support async op")
)

type MongoClient struct {
	client       *mongo.Client
	dbName       string
	ctx          context.Context
	reqChan      chan *OpReq
	closeChan    chan struct{}
	async        bool           // 是否支持异步
	AsyncRetChan chan *AsyncRet // 回传异步操作的结果
}

func NewMongoClient() *MongoClient {
	return &MongoClient{}
}

// Connect 连接到mongoDB
func (mc *MongoClient) Connect(url, dbName string) error {
	if mc.ctx == nil {
		mc.ctx = context.Background()
	}
	ctx, cancel := context.WithTimeout(mc.ctx, Timeout)
	defer cancel()
	clientOpts := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return err
	}
	if err = client.Ping(ctx, nil); err != nil {
		return err
	}
	mc.client = client
	mc.dbName = dbName
	return nil
}

// Disconnect 断开
func (mc *MongoClient) Disconnect() error {
	if mc.client == nil {
		return nil
	}
	return mc.client.Disconnect(mc.ctx)
}

// Collection 集合
func (mc *MongoClient) Collection(collection string) *mongo.Collection {
	return mc.client.Database(mc.dbName).Collection(collection)
}

// DropDatabase 删除数据库
func (mc *MongoClient) DropDatabase() error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	return mc.client.Database(mc.dbName).Drop(mc.ctx)
}

// Tx 事务
func (mc *MongoClient) Tx(ctx context.Context, f func(sessCtx mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	if mc.client == nil {
		return nil, ErrMongoDBNotConnected
	}
	session, err := mc.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)
	return session.WithTransaction(ctx, f)
}

// CreateIndexByCtx 创建索引
func (mc *MongoClient) CreateIndexByCtx(ctx context.Context, collection, key string, unique bool) error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	model := mongo.IndexModel{
		Keys: bson.D{{Key: key, Value: 1}},
	}
	if unique {
		model.Options = options.Index().SetUnique(true)
	}
	_, err := mc.Collection(collection).Indexes().CreateOne(ctx, model)
	return err
}

// CreateIndex 创建索引
func (mc *MongoClient) CreateIndex(collection, key string, unique bool) error {
	ctx, cancel := context.WithTimeout(mc.ctx, 60*time.Second)
	defer cancel()
	return mc.CreateIndexByCtx(ctx, collection, key, unique)
}

// CreateIndexWithExpireByCtx 创建带TTL机制的索引
func (mc *MongoClient) CreateIndexWithExpireByCtx(ctx context.Context, collection, key string, seconds int32) error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	model := mongo.IndexModel{
		Keys: bson.D{{Key: key, Value: 1}},
	}
	if seconds > 0 {
		model.Options = options.Index().SetExpireAfterSeconds(seconds)
	}
	_, err := mc.Collection(collection).Indexes().CreateOne(ctx, model)
	return err
}

// CreateIndexWithExpire 创建带TTL机制的索引
func (mc *MongoClient) CreateIndexWithExpire(ctx context.Context, collection, key string, seconds int32) error {
	ctx, cancel := context.WithTimeout(mc.ctx, 60*time.Second)
	defer cancel()
	return mc.CreateIndexWithExpireByCtx(ctx, collection, key, seconds)
}

// CreateCompoundIndex 创建组合索引
func (mc *MongoClient) CreateCompoundIndex(collection string, keys ...string) error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	if len(keys) == 0 {
		return nil
	}
	d := make(bson.D, len(keys))
	for idx, key := range keys {
		d[idx] = bson.E{Key: key, Value: 1}
	}
	model := mongo.IndexModel{
		Keys: d,
	}
	_, err := mc.Collection(collection).Indexes().CreateOne(mc.ctx, model)
	return err
}

// doFindOne 查找
func (mc *MongoClient) doFindOne(collection string, filter interface{}) (bson.M, error) {
	ret := bson.M{}
	err := mc.Collection(collection).FindOne(mc.ctx, filter).Decode(&ret)
	return ret, err
}

// FindOne 同步查找一条
func (mc *MongoClient) FindOne(collection string, filter interface{}) (bson.M, error) {
	if mc.client == nil {
		return nil, ErrMongoDBNotConnected
	}
	return mc.doFindOne(collection, filter)
}

// FindOneDecode 同步查找一条并反序列到结构体
func (mc *MongoClient) FindOneDecode(collection string, filter interface{}, ret interface{}) error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	return mc.Collection(collection).FindOne(mc.ctx, filter).Decode(ret)
}

// FindOneToDo 同步查找一条并回调
func (mc *MongoClient) FindOneToDo(collection string, filter interface{}, handler func(result *mongo.SingleResult) error) error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	ret := mc.Collection(collection).FindOne(mc.ctx, filter)
	if ret.Err() != nil {
		return ret.Err()
	}
	return handler(mc.Collection(collection).FindOne(mc.ctx, filter))
}

// Find 同步查找返回游标
func (mc *MongoClient) Find(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	if mc.client == nil {
		return nil, ErrMongoDBNotConnected
	}
	return mc.Collection(collection).Find(ctx, filter, opts...)
}

// FindOneAsync 异步查找一条 异步查找一定是需要返回结果的
func (mc *MongoClient) FindOneAsync(collection string, filter interface{}, event Event, context interface{}) error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	if !mc.async {
		return ErrMongoDBNotAsync
	}
	req := NewOpReq(opTypeInsert, collection, "", filter, nil, true, event, context)
	mc.reqChan <- req
	return nil
}

// doInsertOne 插入
func (mc *MongoClient) doInsertOne(collection string, doc interface{}) (*mongo.InsertOneResult, error) {
	return mc.Collection(collection).InsertOne(mc.ctx, doc)
}

// InsertOne 同步插入一条
func (mc *MongoClient) InsertOne(collection string, doc interface{}) error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	_, err := mc.doInsertOne(collection, doc)
	return err
}

// InsertMany 同步插入多条
func (mc *MongoClient) InsertMany(collection string, docs []interface{}) error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	_, err := mc.Collection(collection).InsertMany(mc.ctx, docs)
	return err
}

// InsertOneAsync 异步插入一条
func (mc *MongoClient) InsertOneAsync(collection string, doc interface{}, needRet bool, event Event, context interface{}) error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	if !mc.async {
		return ErrMongoDBNotAsync
	}
	req := NewOpReq(opTypeInsert, collection, "", nil, doc, needRet, event, context)
	mc.reqChan <- req
	return nil
}

// doReplaceOne 全量替换
func (mc *MongoClient) doReplaceOne(collection string, filter interface{}, doc interface{}) (*mongo.UpdateResult, error) {
	return mc.Collection(collection).ReplaceOne(mc.ctx, filter, doc)
}

// ReplaceOne 同步全量替换一条
func (mc *MongoClient) ReplaceOne(collection string, filter interface{}, doc interface{}) error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	_, err := mc.doReplaceOne(collection, filter, doc)
	return err
}

// ReplaceOneAsync 异步全量替换一条
func (mc *MongoClient) ReplaceOneAsync(collection string, filter interface{}, doc interface{}, needRet bool, event Event, context interface{}) error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	if !mc.async {
		return ErrMongoDBNotAsync
	}
	req := NewOpReq(opTypeReplace, collection, "", filter, doc, needRet, event, context)
	mc.reqChan <- req
	return nil
}

// doUpdateOne 更新
func (mc *MongoClient) doUpdateOne(collection string, filter interface{}, doc interface{}) (*mongo.UpdateResult, error) {
	update := bson.M{"$set": doc}
	return mc.Collection(collection).UpdateOne(mc.ctx, filter, update)
}

// UpdateOne 同步更新一条
func (mc *MongoClient) UpdateOne(collection string, filter interface{}, doc interface{}) error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	_, err := mc.doUpdateOne(collection, filter, doc)
	return err
}

// UpdateMany 同步更新多条
func (mc *MongoClient) UpdateMany(collection string, filter interface{}, docs interface{}) error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	_, err := mc.Collection(collection).UpdateMany(mc.ctx, filter, docs)
	return err
}

// UpdateOneAsync 异步更新一条
func (mc *MongoClient) UpdateOneAsync(collection string, filter interface{}, doc interface{}, needRet bool, event Event, context interface{}) error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	if !mc.async {
		return ErrMongoDBNotAsync
	}
	req := NewOpReq(opTypeUpdate, collection, "", filter, doc, needRet, event, context)
	mc.reqChan <- req
	return nil
}

// doUpsertOne 插入或者更新
func (mc *MongoClient) doUpsertOne(collection string, filter interface{}, doc interface{}) (*mongo.UpdateResult, error) {
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)
	return mc.Collection(collection).UpdateOne(mc.ctx, filter, update, opts)
}

// UpsertOne 同步插入一条
func (mc *MongoClient) UpsertOne(collection string, filter interface{}, doc interface{}) error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	_, err := mc.doUpsertOne(collection, filter, doc)
	return err
}

// UpsertMany 同步更新或插入多条
func (mc *MongoClient) UpsertMany(collection string, filter interface{}, docs interface{}) error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	opts := options.Update().SetUpsert(true)
	_, err := mc.Collection(collection).UpdateMany(mc.ctx, filter, docs, opts)
	return err
}

// UpsertOneAsync 异步更新或插入一条
func (mc *MongoClient) UpsertOneAsync(collection string, filter interface{}, doc interface{}, needRet bool, event Event, context interface{}) error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	if !mc.async {
		return ErrMongoDBNotAsync
	}
	req := NewOpReq(opTypeUpsert, collection, "", filter, doc, needRet, event, context)
	mc.reqChan <- req
	return nil
}

// doDeleteOne 删除
func (mc *MongoClient) doDeleteOne(collection string, filter interface{}) (bool, error) {
	ret, err := mc.Collection(collection).DeleteOne(mc.ctx, filter)
	return ret.DeletedCount > 0, err
}

// DeleteOne 同步删除一条
func (mc *MongoClient) DeleteOne(collection string, filter interface{}) (bool, error) {
	if mc.client == nil {
		return false, ErrMongoDBNotConnected
	}
	return mc.doDeleteOne(collection, filter)
}

// DeleteMany 同步删除多条
func (mc *MongoClient) DeleteMany(collection string, filter interface{}) (int64, error) {
	if mc.client == nil {
		return 0, ErrMongoDBNotConnected
	}
	ret, err := mc.Collection(collection).DeleteMany(mc.ctx, filter)
	return ret.DeletedCount, err
}

// DeleteOneAsync 异步删除一条
func (mc *MongoClient) DeleteOneAsync(collection string, filter interface{}, needRet bool, event Event, context interface{}) (bool, error) {
	if mc.client == nil {
		return false, ErrMongoDBNotConnected
	}
	if !mc.async {
		return false, ErrMongoDBNotAsync
	}
	req := NewOpReq(opTypeDelete, collection, "", filter, nil, needRet, event, context)
	mc.reqChan <- req
	return false, nil
}

func (mc *MongoClient) IsGridFileExist(filename, bucketName string) bool {
	opts := options.GridFSBucket().SetName(bucketName)
	bucket, err := gridfs.NewBucket(mc.client.Database(mc.dbName), opts)
	if err != nil {
		panic(err)
	}
	cursor, err := bucket.Find(bson.M{"filename": filename})
	if err != nil {
		panic(err)
	}
	return cursor.RemainingBatchLength() == 0
}

// ReadGridFile 读gridFS
func (mc *MongoClient) ReadGridFile(filename, bucketName string) ([]byte, error) {
	if mc.client == nil {
		return nil, ErrMongoDBNotConnected
	}
	opts := options.GridFSBucket().SetName(bucketName)
	bucket, err := gridfs.NewBucket(mc.client.Database(mc.dbName), opts)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(nil)
	if _, err = bucket.DownloadToStreamByName(filename, buffer); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// doWriteGridFile 写gridFS
func (mc *MongoClient) doWriteGridFile(filename, bucketName string, data []byte) error {
	opts := options.GridFSBucket().SetName(bucketName)
	bucket, err := gridfs.NewBucket(mc.client.Database(mc.dbName), opts)
	if err != nil {
		return err
	}
	var newFileID primitive.ObjectID
	newFileID, err = bucket.UploadFromStream(filename, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	cursor, err := bucket.Find(bson.M{"filename": filename, "_id": bson.M{"$ne": newFileID}})
	if err != nil {
		return err
	}
	defer func() {
		if err := cursor.Close(mc.ctx); err != nil {
			log.Errorf("close mongo cursor:%s", err)
		}
	}()
	type gridFS struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	var fileList []gridFS
	if err = cursor.All(mc.ctx, &fileList); err != nil {
		return err
	}
	for _, file := range fileList {
		if err = bucket.Delete(file.ID); err != nil {
			log.Errorf("delete grid fs err:%s", err)
		}
	}
	return err
}

// WriteGridFile 同步写gridFS
func (mc *MongoClient) WriteGridFile(filename, bucketName string, data []byte) error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	return mc.doWriteGridFile(filename, bucketName, data)
}

// WriteGridFileAsync 异步写gridFS
func (mc *MongoClient) WriteGridFileAsync(filename, bucketName string, data []byte, needRet bool, event Event, context interface{}) error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	if !mc.async {
		return ErrMongoDBNotAsync
	}
	req := NewOpReq(opTypeUploadGridFS, bucketName, filename, nil, data, needRet, event, context)
	mc.reqChan <- req
	return nil
}

func (mc *MongoClient) DeleteGridFile(filename, bucketName string) error {
	if mc.client == nil {
		return ErrMongoDBNotConnected
	}
	opts := options.GridFSBucket().SetName(bucketName)
	bucket, err := gridfs.NewBucket(mc.client.Database(mc.dbName), opts)
	if err != nil {
		return err
	}
	cursor, err := bucket.Find(bson.M{"filename": filename})
	if err != nil {
		return err
	}
	defer func() {
		if err := cursor.Close(mc.ctx); err != nil {
			log.Errorf("close mongo cursor:%s", err)
		}
	}()
	type gridFS struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	var fileList []gridFS
	if err = cursor.All(mc.ctx, &fileList); err != nil {
		return err
	}
	if fileList == nil || len(fileList) != 1 {
		return nil
	}
	if err = bucket.Delete(fileList[0].ID); err != nil {
		return err
	}
	return nil
}
