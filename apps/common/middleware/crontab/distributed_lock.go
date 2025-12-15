package crontab

import (
	"context"
	"time"

	// "github.com/beego/beego/v2/client/cache/redis"

	"github.com/redis/go-redis/v9"
)

type DistributedLock struct {
	client *redis.Client
	key    string
	ttl    time.Duration
}

func NewDistributedLock(client *redis.Client, key string, ttl time.Duration) *DistributedLock {
	return &DistributedLock{
		client: client,
		key:    key,
		ttl:    ttl,
	}
}

func (dl *DistributedLock) Acquire(ctx context.Context) (bool, error) {
	return dl.client.SetNX(ctx, dl.key, "locked", dl.ttl).Result()
}

func (dl *DistributedLock) Release(ctx context.Context) error {
	return dl.client.Del(ctx, dl.key).Err()
}

// // controllers/task.go 中的任务执行函数
// func (c *TaskController) backupDatabase() {
//     ctx := context.Background()
//     lock := utils.NewDistributedLock(redisClient, "lock:backup_database", time.Minute*10)

//     acquired, err := lock.Acquire(ctx)
//     if err != nil || !acquired {
//         // 获取锁失败，说明其他节点正在执行，直接返回
//         return
//     }
//     defer lock.Release(ctx)

//     // 执行实际的备份任务
//     // ...
// }
