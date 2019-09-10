package etcd

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	cl, err := NewClient(clientv3.Config{
		Endpoints: []string{"144.202.13.31:2379"},
		DialTimeout: 50 * time.Second,
	})
	fmt.Println(err)
	if err != nil {
		return
	}
	cl.Set("/ashe/jobs/wuranxu", "running")
	fmt.Println(cl.GetSingle("/ashe/jobs/wuranxu"))

}
