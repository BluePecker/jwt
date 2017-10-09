package storage

import (
	"fmt"
)

type (
	Engine interface {
		Read(key string) (v interface{}, err error)

		ReadInt(key string) (v int, err error)

		ReadString(key string) (string, error)

		Upgrade(key string, expire int)

		Initializer(opts string) error

		Set(key string, value interface{}, expire int) error

		TTL(key string) float64

		SetImmutable(key string, value interface{}, expire int) error

		Remove(key string)

		LKeep(key string, value interface{}, maxLen, expire int) error

		LRange(key string, start, stop int) ([]string, error)

		LExist(key string, value interface{}) bool
	}
)

var provider = make(map[string]Engine)

func Register(name string, driver Engine) {
	if driver == nil {
		panic("storage: register driver is nil")
	}
	if _, find := provider[name]; find {
		panic("storage: register called twice for " + name)
	}

	provider[name] = driver
}

func New(name string, opts string) (*Engine, error) {
	if storage, find := provider[name]; !find {
		return nil, fmt.Errorf("storage: unknown driver %q (forgotten import?)", name)
	} else {
		if err := storage.Initializer(opts); err != nil {
			return nil, fmt.Errorf("storage: %q driver init failed", name);
		}
		return &storage, nil;
	}
}
