package etcd

import (
	"context"
	"github.com/coreos/etcd/clientv3/concurrency"
	"time"
)

const (
	// 上锁超时时间
	TIMEOUT = 10
	TTL     = 5
)

type DistributeLock struct {
	mutex   *concurrency.Mutex
	session *concurrency.Session
	timeout int
}

func NewLock(client *Client, prefix string, timeout ...int) (*DistributeLock, error) {
	session, err := concurrency.NewSession(client.cli, concurrency.WithTTL(TTL))
	if err != nil {
		return nil, err
	}
	tm := TIMEOUT
	if len(timeout) > 0 {
		tm = timeout[0]
	}
	return &DistributeLock{
		timeout: tm,
		session: session,
		mutex:   concurrency.NewMutex(session, prefix),
	}, nil
}

func (l *DistributeLock) Lock() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.timeout)*time.Second)
	defer cancel()
	return l.mutex.Lock(ctx)
}

func (l *DistributeLock) UnLock() error {
	return l.mutex.Unlock(context.Background())
}

func (l *DistributeLock) Close() error {
	return l.session.Close()
}
