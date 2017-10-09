package uri

import (
	"net/url"
	"github.com/go-redis/redis"
	"fmt"
	"reflect"
	"strings"
	"strconv"
	"time"
)

func extra(options interface{}, parsed *url.URL) interface{} {
	indirect := reflect.Indirect(reflect.ValueOf(options))
	for name, arr := range parsed.Query() {
		extra := indirect.FieldByName(name)
		if !extra.CanSet() {
			continue
		}
		switch extra.Interface().(type) {
		case int, int8, int16, int32, int64:
			integer, err := strconv.ParseInt(arr[0], 10, 64)
			if err == nil {
				extra.SetInt(integer)
			}
			break;
		case string:
			extra.SetString(arr[0])
			break;
		case bool:
			boolean, err := strconv.ParseBool(arr[0])
			if err == nil {
				extra.SetBool(boolean)
			}
			break;
		case time.Duration:
			duration, _ := time.ParseDuration(arr[0] + "ms")
			extra.Set(reflect.ValueOf(duration))
		}
	}
	return options
}

func Parser(uri string) (interface{}, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("wrong uri format: %s", uri)
	}
	if parsed.Scheme != "redis" {
		return nil, fmt.Errorf("valid schema: %s", parsed.Scheme)
	}
	var pwd string
	if parsed.User != nil {
		pwd, _ = parsed.User.Password()
	}
	var addr = strings.Split(parsed.Host, ",")
	if len(addr) == 0 {
		return nil, fmt.Errorf("%s", "valid host")
	}
	if len(addr) == 1 {
		database, _ := strconv.Atoi(strings.Trim(parsed.Path, "/"))
		options := &redis.Options{
			Addr:     addr[0],
			DB:       database,
			Password: pwd,
		}
		return extra(options, parsed), nil
	}
	options := &redis.ClusterOptions{
		Addrs:    addr,
		Password: pwd,
	}
	return extra(options, parsed), nil
}
