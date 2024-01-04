// -------------------------------------------
// @file      : neighbor.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 上午10:13
// -------------------------------------------

package gscluster

// Neighbor 集群中的邻居
type Neighbor struct {
	clusterRemote *ClusterRemote
	services      map[string]IRemoteService // 按名字索引的邻居节点上的远程服务
	servicesByID  map[ID]IRemoteService     // 按ID索引的邻居节点上的远程服务
}

// NewNeighbor 新建邻居节点
func NewNeighbor(remote *ClusterRemote) *Neighbor {
	return &Neighbor{
		clusterRemote: remote,
		services:      make(map[string]IRemoteService),
		servicesByID:  make(map[ID]IRemoteService),
	}
}
