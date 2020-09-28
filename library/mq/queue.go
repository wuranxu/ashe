package mq
//
//import (
//	"context"
//	"github.com/apache/rocketmq-client-go/primitive"
//	"github.com/apache/rocketmq-client-go/producer"
//)
//
//func Push() error {
//	p, err := producer.NewDefaultProducer(
//		producer.WithNameServer([]string{"127.0.0.1:9876"}))
//	if err != nil {
//		return err
//	}
//	err = p.Start()
//	if err != nil {
//		return err
//	}
//	err = p.SendOneWay(context.Background(), &primitive.Message{
//		Topic: "test",
//		Body:  []byte("Hello RocketMQ Go Client!"),
//	})
//	if err != nil {
//		return err
//	}
//	defer p.Shutdown()
//	return nil
//}
