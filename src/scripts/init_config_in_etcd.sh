echo ""
ADDR=127.0.0.1:2379
etcdctl --endpoints=$ADDR put GOGS/LOGIN_CONFIG/1 '{"Type":"LOGIN_CONFIG","ID":1,"MongoID":1,"IsBlockLogin":0}'
etcdctl --endpoints=$ADDR put GOGS/GAME_CONFIG/1 '{"Type":"GAME_CONFIG","ID":1,"MongoID":1}'
etcdctl --endpoints=$ADDR put GOGS/MONGO/1 '{"Type":"MONGO","ID":1,"AddrList":["127.0.0.1:27017"],"User":"","Password":"","ReplicaSet":""}'