package main

import (
	"internInformer/service/informer"
	"internInformer/service/scheduler"
	"internInformer/service/spider"
)

func main() {
	sche := scheduler.NewDefaultScheduler()
	mihoyoSpider := spider.NewDefaultMihoyoInternSpider()
	pushdeer := informer.NewDefaultPushdeer()
	sche.Register(&scheduler.Instance{
		Spider:   mihoyoSpider,
		Informer: pushdeer,
	})

	sche.Start()
}
