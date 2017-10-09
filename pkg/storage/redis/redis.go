package redis

import (
	"strconv"
	"sync"
	"time"
	"github.com/BluePecker/JwtAuth/pkg/storage"
	"github.com/go-redis/redis"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/BluePecker/JwtAuth/pkg/storage/redis/uri"
	"reflect"
	"github.com/Sirupsen/logrus"
)

type (
	Redis struct {
		create  time.Time
		storage Client
		mu      sync.RWMutex
	}

	Client interface {
		Ping() *redis.StatusCmd

		TTL(key string) *redis.DurationCmd

		LRange(key string, start int64, stop int64) *redis.StringSliceCmd

		HSet(key string, field string, value interface{}) *redis.BoolCmd

		HGet(key string, field string) *redis.StringCmd

		Del(key ...string) *redis.IntCmd

		Close() error

		Expire(key string, expiration time.Duration) *redis.BoolCmd

		Pipelined(fn func(redis.Pipeliner) error) ([]redis.Cmder, error)
	}
)

func inject(from, target reflect.Value) {
	indirect := reflect.Indirect(target.Elem())
	for index := 0; index < from.Elem().NumField(); index++ {
		name := from.Elem().Type().Field(index).Name
		f1 := from.Elem().FieldByName(name)
		f2 := indirect.FieldByName(name)
		if f2.IsValid() {
			if f1.Type() == f2.Type() && f2.CanSet() {
				f2.Set(f1)
			}
		}
	}
}

func (R *Redis) Initializer(opts string) error {
	generic, err := uri.Parser(opts)
	if err != nil {
		return err
	}
	switch reflect.ValueOf(generic).Interface().(type) {
	case *redis.ClusterOptions:
		options := &redis.ClusterOptions{}
		inject(reflect.ValueOf(generic), reflect.ValueOf(options))
		R.storage = redis.NewClusterClient(options)
		break
	case *redis.Options:
		options := &redis.Options{}
		inject(reflect.ValueOf(generic), reflect.ValueOf(options))
		R.storage = redis.NewClient(options)
		break
	}
	statusCmd := R.storage.Ping()
	if statusCmd.Err() != nil {
		logrus.Error(statusCmd.Err())
		defer R.storage.Close()
	}
	return statusCmd.Err()
}

func (R *Redis) TTL(key string) float64 {
	R.mu.RLock()
	defer R.mu.RUnlock()
	return R.storage.TTL(R.md5Key(key)).Val().Seconds()
}

func (R *Redis) Read(key string) (interface{}, error) {
	R.mu.RLock()
	defer R.mu.RUnlock()
	status := R.get(R.md5Key(key))
	return status.Val(), status.Err()
}

func (R *Redis) ReadInt(key string) (int, error) {
	R.mu.RLock()
	defer R.mu.RUnlock()
	status := R.get(R.md5Key(key))
	if status.Err() != nil {
		return 0, status.Err()
	}
	return strconv.Atoi(status.Val())
}

func (R *Redis) ReadString(key string) (string, error) {
	R.mu.RLock()
	defer R.mu.RUnlock()
	status := R.get(R.md5Key(key))
	if status.Err() != nil {
		return "", status.Err()
	}
	return status.Val(), nil
}

func (R *Redis) Upgrade(key string, expire int) {
	R.mu.Lock()
	defer R.mu.Unlock()
	key = R.md5Key(key)
	if v, err := R.Read(key); err != nil {
		R.Set(key, v, expire)
	}
}

func (R *Redis) Set(key string, value interface{}, expire int) error {
	R.mu.Lock()
	defer R.mu.Unlock()
	return R.save(R.md5Key(key), value, expire, false)
}

func (R *Redis) SetImmutable(key string, value interface{}, expire int) error {
	R.mu.Lock()
	defer R.mu.Unlock()
	return R.save(R.md5Key(key), value, expire, true)
}

func (R *Redis) Remove(key string) {
	R.mu.Lock()
	defer R.mu.Unlock()
	R.remove(R.md5Key(key))
}

func (R *Redis) LKeep(key string, value interface{}, maxLen, expire int) error {
	R.mu.Lock()
	defer R.mu.Unlock()
	key = R.md5Key(key)
	_, err := R.storage.Pipelined(func(pip redis.Pipeliner) error {
		pip.LPush(key, value)
		pip.LTrim(key, 0, int64(maxLen-1))
		pip.Expire(key, time.Duration(expire)*time.Second)
		return nil;
	})
	return err;
}

func (R *Redis) LRange(key string, start, stop int) ([]string, error) {
	R.mu.Lock()
	defer R.mu.Unlock()
	key = R.md5Key(key)
	cmd := R.storage.LRange(key, int64(start), int64(stop))
	return cmd.Val(), cmd.Err()
}

func (R *Redis) LExist(key string, value interface{}) bool {
	if strArr, err := R.LRange(key, 0, -1); err == nil {
		for _, v := range strArr {
			if v == value.(string) {
				return true
			}
		}
	}
	return false
}

func (R *Redis) remove(key string) error {
	status := R.storage.Del(key)
	return status.Err()
}

func (R *Redis) get(key string) *redis.StringCmd {
	return R.storage.HGet(R.md5Key(key), "v")
}

func (R *Redis) save(key string, value interface{}, expire int, immutable bool) error {
	key = R.md5Key(key)
	cmd := R.storage.HGet(key, "i")
	if find, _ := strconv.ParseBool(cmd.Val()); find {
		return fmt.Errorf("this key(%s) write protection", key)
	}
	R.storage.Pipelined(func(pipe redis.Pipeliner) error {
		pipe.HSet(key, "v", value)
		pipe.HSet(key, "i", immutable)
		pipe.Expire(key, time.Duration(expire)*time.Second)
		return nil
	})
	return nil
}

func (R *Redis) md5Key(key string) string {
	hash := md5.New()
	hash.Write([]byte(key))
	return hex.EncodeToString(hash.Sum([]byte("jwt#")))
}

func init() {
	storage.Register("redis", &Redis{})
}
