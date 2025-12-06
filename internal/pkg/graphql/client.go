package graphql

import (
	"beedance-mcp/configs"
	"net/http"
	"sync"
	"time"
)

var clientPool *sync.Pool

type GraphqlClient struct {
	url    string
	client *http.Client
}

func InitClient() {
	newFunc := func() interface{} {
		transport := &http.Transport{}
		transport.MaxIdleConns = 100
		transport.MaxIdleConnsPerHost = 10
		transport.IdleConnTimeout = 90 * time.Second
		transport.TLSHandshakeTimeout = 10 * time.Second

		httpClient := &http.Client{}
		httpClient.Timeout = 30 * time.Second
		httpClient.Transport = transport

		return &GraphqlClient{
			client: httpClient,
			url:    configs.GlobalConfig.Gateway.URL,
		}
	}

	pool := &sync.Pool{New: newFunc}
	clientPool = pool
}

func GetClient() *GraphqlClient {
	return clientPool.Get().(*GraphqlClient)
}

func PutClient(client *GraphqlClient) {
	if client == nil {
		return
	}
	clientPool.Put(client)
}
