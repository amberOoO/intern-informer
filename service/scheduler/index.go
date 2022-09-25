package scheduler

import (
	"internInformer/service/informer"
	"internInformer/service/spider"
)

type Instance struct {
	Spider   spider.SpiderInterface
	Informer informer.InformerInterface
}

type SchedulerInterface interface {
	Register(Instance)
	Start()
}
