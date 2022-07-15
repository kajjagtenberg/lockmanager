package lockmanager

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ErrResourceLocked = fmt.Errorf("resource already locked")
)

type LockManager struct {
	rds    *redis.Client
	prefix string
}

func (m *LockManager) Lock(ctx context.Context, key string, fn func() error) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	fullKey := fmt.Sprintf("%s%s", m.prefix, key)

	locked, err := m.rds.SetNX(ctx, fullKey, hostname, time.Second*10).Result()
	if err != nil {
		return err
	}

	if !locked {
		return ErrResourceLocked
	}

	defer m.rds.Del(ctx, fullKey)

	c := make(chan error, 1)

	go func() {
		c <- fn()
	}()

	for {
		select {
		case err := <-c:
			return err
		case <-time.After(time.Second * 5):
			if err := m.rds.Expire(ctx, fullKey, time.Second*10).Err(); err != nil {
				return err
			}
		}
	}

}

func (m *LockManager) SetPrefix(prefix string) {
	m.prefix = prefix
}

func (m *LockManager) Prefix() string {
	return m.prefix
}

func NewLockManager(rds *redis.Client) (*LockManager, error) {
	return &LockManager{rds: rds}, nil
}
