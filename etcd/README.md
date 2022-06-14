# go-pkg ETCD Package

[![Go Report Card](https://goreportcard.com/badge/github.com/wwwangxc/go-pkg/etcd)](https://goreportcard.com/report/github.com/wwwangxc/go-pkg/etcd)
[![GoDoc](https://pkg.go.dev/badge/github.com/wwwangxc/go-pkg/etcd?status.svg)](https://pkg.go.dev/github.com/wwwangxc/go-pkg/etcd)

go-pkg/etcd is an componentized etcd client package.

It providels:

- An easy way to configre and manage etcd v3 client.

Based on [go.etcd.io/etcd/client/v3](https://github.com/etcd-io/etcd/tree/main/client/v3)


## Install

```sh
go get github.com/wwwangxc/go-pkg/etcd
```

## Quick Start

### Client Proxy

```go
package main

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"

	// go-pkg/etcd will automatically read configuration
	// files (./app.yaml) when package loaded
	"github.com/wwwangxc/go-pkg/etcd"
)

func ExampleNewClientProxy() {
	_ = etcd.NewClientProxy("etcd1",
		etcd.WithTarget("127.0.0.1:2379,127.0.0.1:2380"),       // set target, target is a list of URLs, multiple URL split by ','.
		etcd.WithTimeout(1000),                                 // set timeout, unit millisecond, default 3000.
		etcd.WithUsername("username"),                          // set user name for authentication.
		etcd.WithPassword("password"),                          // set password for authentication.
		etcd.WithTLSKeyPath("/usr/local/etcd_conf/key.pem"),    // set tls key file path.
		etcd.WithTLSCertPath("/usr/local/etcd_conf/cert.pem"),  // set tls cert file path.
		etcd.WithCACertPath("/usr/local/etcd_conf/cacert.pem"), // set ca cert file path.
	)
}

func ExampleClientProxy_Put() {
	// do etcd put operation
	_, err := etcd.NewClientProxy("etcd1").Put(context.Background(), "key", "val")
	if err != nil {
		fmt.Printf("put operation fail. error:%v", err)
		return
	}

	// or

	// do etcd put operation and convert result to an error
	if err = etcd.PutResult(etcd.NewClientProxy("etcd").Put(context.Background(), "key", "val")); err != nil {
		fmt.Printf("put operation fail. error:%v", err)
		return
	}
}

func ExampleClientProxy_PutWithLease() {
	cli := etcd.NewClientProxy("etcd1")
	lease := cli.Lease()

	// create a lease
	id, err := etcd.LeaseGrantResult(lease.Grant(context.Background(), 10))
	if err != nil {
		fmt.Printf("lease grant fail. error:%v", err)
		return
	}

	// put with lease
	err = etcd.PutResult(cli.Put(context.Background(), "key", "val", clientv3.WithLease(id)))
	if err != nil {
		fmt.Printf("put operation fail. error:%v", err)
	}
}

func ExampleClientProxy_Get() {
	// do etcd get operation
	resp, err := etcd.NewClientProxy("etcd1").Get(context.Background(), "key")
	if err != nil {
		fmt.Printf("get operation fail. error:%v", err)
		return
	}

	for k, v := range resp.Kvs {
		fmt.Printf("key: %s\n", k)
		fmt.Printf("val: %s\n", v)
	}

	// or

	// do etcd get operation and convert result to map[string]string and an error
	m, err := etcd.GetResult(etcd.NewClientProxy("etcd").Get(context.Background(), "key"))
	if err != nil {
		fmt.Printf("get operation fail. error:%v", err)
		return
	}

	for k, v := range m {
		fmt.Printf("key: %s\n", k)
		fmt.Printf("val: %s\n", v)
	}
}

func ExampleClientProxy_Delete() {
	// do etcd delete operation
	resp, err := etcd.NewClientProxy("etcd1").Delete(context.Background(), "key")
	if err != nil {
		fmt.Printf("delete operation fail. error:%v", err)
		return
	}

	fmt.Printf("number of keys deleted: %d\n", resp.Deleted)

	// or

	// do etcd delte operation and convert result to map[string]string and an error
	deletedNumber, err := etcd.DeleteResult(etcd.NewClientProxy("etcd").Delete(context.Background(), "key"))
	if err != nil {
		fmt.Printf("get operation fail. error:%v", err)
		return
	}

	fmt.Printf("number of keys deleted: %d\n", deletedNumber)
}

func ExampleClientProxy_Watch() {
	// do etcd watch operation
	watchChan, err := etcd.NewClientProxy("etcd1").Watch(context.Background(), "key")
	if err != nil {
		fmt.Printf("watch operation fail. error:%v", err)
		return
	}

	for v := range watchChan {
		// do something...
	}
}

func ExampleClientProxy_Txn() {
	_, err := etcd.NewClientProxy("etcd1").Txn(context.Background(),
		[]clientv3.Cmp{clientv3.Compare(clientv3.Value("key"), "=", "val")}, // if key's value == val
		[]clientv3.Op{clientv3.OpPut("key", "val1")},                        // then put key's value = val1
		[]clientv3.Op{clientv3.OpPut("key", "val")})                         // else put key's value = val
	if err != nil {
		fmt.Printf("txn fail. error:%v", err)
		return
	}

	// or

	err = etcd.TxnResult(etcd.NewClientProxy("etcd1").Txn(context.Background(),
		[]clientv3.Cmp{clientv3.Compare(clientv3.Value("key"), "=", "val")}, // if key's value == val
		[]clientv3.Op{clientv3.OpPut("key", "val1")},                        // then put key's value = val1
		[]clientv3.Op{clientv3.OpPut("key", "val")}))                        // else put key's value = val
	if err != nil {
		fmt.Printf("txn fail. error:%v", err)
		return
	}

}
```

### Lease Proxy

```go
package main

import (
	"context"
	"fmt"
	"time"

	// go-pkg/etcd will automatically read configuration
	// files (./app.yaml) when package loaded
	"github.com/wwwangxc/go-pkg/etcd"
)

func ExampleLeaseProxy_Grant() {
	lease := etcd.NewClientProxy("etcd1").Lease()

	// create a lease
	id, err := etcd.LeaseGrantResult(lease.Grant(context.Background(), 10))
	if err != nil {
		fmt.Printf("lease grant fail. error:%v", err)
		return
	}

	fmt.Printf("lease:0x%x\n", id)
}

func ExampleLeaseProxy_RevokeResult() {
	lease := etcd.NewClientProxy("etcd1").Lease()

	// create a lease
	id, err := etcd.LeaseGrantResult(lease.Grant(context.Background(), 10))
	if err != nil {
		fmt.Printf("lease grant fail. error:%v", err)
		return
	}

	// revoke a lease
	err = etcd.LeaseRevokeResult(lease.Revoke(context.Background(), id))
	if err != nil {
		fmt.Printf("lease revoke fail. error:%v", err)
		return
	}
}

func ExampleLeaseProxy_LeaseTimeToLiveResult() {
	lease := etcd.NewClientProxy("etcd1").Lease()

	// create a lease
	id, err := etcd.LeaseGrantResult(lease.Grant(context.Background(), 10))
	if err != nil {
		fmt.Printf("lease grant fail. error:%v", err)
		return
	}

	for {
		// get lease ttl
		ttl, err := etcd.LeaseTimeToLiveResult(lease.TimeToLive(context.Background(), id))
		if err != nil {
			fmt.Printf("get lease ttl fail. error:%v", err)
			return
		}

		if ttl == -1 {
			fmt.Printf("lease:0x%x expired\n", id)
			break
		}

		time.Sleep(time.Second)
	}
}

func ExampleLeaseProxy_KeepAlive() {
	lease := etcd.NewClientProxy("etcd1").Lease()

	// create a lease
	id, err := etcd.LeaseGrantResult(lease.Grant(context.Background(), 10))
	if err != nil {
		fmt.Printf("lease grant fail. error:%v", err)
		return
	}

	ch, err := lease.KeepAlive(context.Background(), id)
	if err != nil {
		fmt.Printf("lease keep alive fail. error:%v", err)
		return
	}

	for {
		ka := <-ch
		if ka == nil {
			fmt.Println("lease timeout")
			return
		}
		fmt.Println("ttl:", ka.TTL)
	}
}
```

### config

app.yaml

```yaml
client:
  timeout: 3000
  service:
    - name: etcd1
      target: 127.0.0.1:2379,127.0.0.1:2380
      timeout: 1000
      username: username
      password: password
      tls_key: /usr/local/etcd_conf/key.pem
      tls_cert: /usr/local/etcd_conf/cert.pem
      ca_cert: /usr/local/etcd_conf/cacert.pem
```

## Hot To Mock

### Client & Lease Proxy

```go
package tests

import (
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/golang/mock/gomock"
	clientv3 "go.etcd.io/etcd/client/v3"

	// go-pkg/etcd will automatically read configuration
	// files (./app.yaml) when package loaded
	"github.com/wwwangxc/go-pkg/etcd"
	"github.com/wwwangxc/go-pkg/etcd/mocketcd"
)

func TestMockClientProxy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock Lease Proxy
	mockLease := mocketcd.NewMockLeaseProxy(ctrl)
	mockLease.EXPECT().Grant(gomock.Any(), gomock.Any()).Return(&clientv3.LeaseGrantResponse{}, nil).AnyTimes()
	mockLease.EXPECT().Revoke(gomock.Any(), gomock.Any()).Return(&clientv3.LeaseRevokeResponse{}, nil).AnyTimes()
	mockLease.EXPECT().TimeToLive(gomock.Any(), gomock.Any(), gomock.Any()).Return(&clientv3.LeaseTimeToLiveResponse{}, nil).AnyTimes()
	mockLease.EXPECT().Leases(gomock.Any()).Return(&clientv3.LeaseLeasesResponse{}, nil).AnyTimes()
	mockLease.EXPECT().KeepAlive(gomock.Any(), gomock.Any()).Return(make(chan *clientv3.LeaseKeepAliveResponse), nil).AnyTimes()

	// Mock Client Proxy
	mockCli := mocketcd.NewMockClientProxy(ctrl)
	mockCli.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&clientv3.PutResponse{}, nil).AnyTimes()
	mockCli.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(&clientv3.GetResponse{}, nil).AnyTimes()
	mockCli.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(&clientv3.DeleteResponse{}, nil).AnyTimes()
	mockCli.EXPECT().Watch(gomock.Any(), gomock.Any(), gomock.Any()).Return(make(chan clientv3.WatchChan), nil).AnyTimes()
	mockCli.EXPECT().Txn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&clientv3.TxnResponse{}, nil).AnyTimes()
	mockCli.EXPECT().Lease().Return(mockLease).AnyTimes()

	patches := gomonkey.ApplyFunc(etcd.NewClientProxy,
		func(string, ...etcd.ClientOption) etcd.ClientProxy {
			return mockCli
		})
	defer patches.Reset()

	// do something...
}
```
