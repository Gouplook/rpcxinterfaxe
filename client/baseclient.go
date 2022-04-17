package client

import (
	"context"
	"fmt"
	"github.com/Gouplook/dzbase/common/plugins/jaeger"
	"github.com/Gouplook/dzgin"
	"github.com/smallnest/rpcx/client"
	"net/http"

	"sync"
)

var (
	rpcPools map[string]map[string]*client.XClientPool
	lock     *sync.RWMutex
)

type BaseClient struct {
	ServiceName string
	ServicePath string
	discovery   client.ServiceDiscovery
	xClient     client.XClient
}

func init() {
	rpcPools = map[string]map[string]*client.XClientPool{}
	lock = new(sync.RWMutex)
}

func (c *BaseClient) getPools(serviceName string, servicePath string) client.XClient {
	if service, ok := rpcPools[serviceName]; ok {
		if rpcPool, ok := service[servicePath]; ok {
			return rpcPool.Get()
		} else {
			lock.Lock()
			rpcPool, ok := service[servicePath]
			if !ok {
				rpcPool = client.NewXClientPool(dzgin.AppConfig.DefaultInt("rpc_pool_count", 10), c.ServicePath, client.Failtry, client.RandomSelect, c.GetDiscovery(), client.DefaultOption)
				rpcPools[serviceName][servicePath] = rpcPool
			}
			lock.Unlock()
			return rpcPool.Get()
		}
	} else {
		lock.Lock()
		service, ok := rpcPools[serviceName]
		if !ok {
			service = map[string]*client.XClientPool{
				servicePath: client.NewXClientPool(dzgin.AppConfig.DefaultInt("rpc_pool_count", 10), c.ServicePath, client.Failtry, client.RandomSelect, c.GetDiscovery(), client.DefaultOption),
			}
			rpcPools[serviceName] = service
		}
		lock.Unlock()
		return service[servicePath].Get()
	}
}

func (c *BaseClient) getXClient() client.XClient {
	return c.getPools(c.ServiceName, c.ServicePath)
}
func (c *BaseClient) GetDiscovery() client.ServiceDiscovery {
	if c.discovery == nil {
		address := dzgin.AppConfig.String(c.ServiceName)
		c.discovery, _ = client.NewPeer2PeerDiscovery(address, "")
	}
	return c.discovery
}

func (c *BaseClient) Call(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error {

	span, ctx, spanErr := jaeger.RpcxSpanWithContext(ctx, fmt.Sprintf("调用%s服务的%s方法", c.ServicePath, serviceMethod), &http.Request{})
	if spanErr == nil {
		span.SetTag("参数", args)
		defer span.Finish()
	}

	err := c.getXClient().Call(ctx, serviceMethod, args, reply)
	if err != nil && spanErr == nil {
		span.SetTag("error", true)
		span.SetTag("错误信息", fmt.Sprint(err))
	}
	return err
}

func (c *BaseClient) Close() {

}
