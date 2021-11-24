package watch

import (
	"k8s.io/client-go/util/workqueue"
)

func LimitQueue() workqueue.RateLimitingInterface {
	limiter := workqueue.DefaultControllerRateLimiter()
	queue := workqueue.NewRateLimitingQueue(limiter) //队列支持重试，同时会记录重试次数
	return queue
}
