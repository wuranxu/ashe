package etcd
//
//import (
//	"ashe/common"
//	"fmt"
//	"sync"
//	"testing"
//	"time"
//)
//
//func dosomething() {
//	fmt.Println("我正在shuimian")
//	time.Sleep(2 * time.Second)
//	fmt.Println("我睡眠完了")
//}
//
//func client(cl *Client, index int) error {
//	//lock, err := NewLock(cl, "/wuranxu/test/lock", 300)
//	lock, err := NewMutex("/wuranxu/test/lock", cl.cli)
//	if err != nil {
//		return err
//	}
//	if err = lock.Lock(); err != nil {
//		// 没有获取到锁或者其他错误
//		fmt.Println("ID: ", index, "没有获取到锁哦, error: ", err)
//		return err
//	}
//	fmt.Println("ID: ", index, "获取锁成功")
//	// 获取到了锁则执行这个任务撒
//	defer func() {
//		err := lock.Unlock()
//		if err != nil {
//			fmt.Println("删除锁失败")
//			return
//		}
//		fmt.Println("删除锁成功")
//	}()
//	dosomething()
//	return nil
//
//}
//
//func TestDistributeLock_Lock(t *testing.T) {
//	cli, err := NewClient(common.EtcdConfig{
//		Endpoints:   []string{"106.13.173.14:2371"},
//		DialTimeout: 10,
//	})
//	if err != nil {
//		t.Fatal("etcd连接失败")
//	}
//	var wg sync.WaitGroup
//	for i := 0; i < 10; i++ {
//		wg.Add(1)
//		go func(data int) {
//			defer wg.Done()
//			fmt.Println(client(cli, data))
//		}(i)
//	}
//	wg.Wait()
//}
