package ram

import (
    "fmt"
    "reflect"
    "strconv"
    "errors"
    "time"
)

type (
    Entry struct {
        value     interface{}
        ttl       int64
        immutable bool
        version   int64
    }
    
    MemStore map[string]Entry
)

func (e *Entry) Value() interface{} {
    if !e.immutable {
        return e.value
    }
    vv := reflect.Indirect(reflect.ValueOf(e.value))
    switch vv.Type().Kind() {
    case reflect.Map:
        newMap := reflect.MakeMap(vv.Type())
        for _, k := range vv.MapKeys() {
            newMap.SetMapIndex(k, vv.MapIndex(k))
        }
        return newMap
    case reflect.Slice:
        newSlice := reflect.MakeSlice(vv.Type(), vv.Len(), vv.Cap())
        reflect.Copy(newSlice, vv)
        return newSlice
    default:
        return vv.Interface()
    }
}

func (MS *MemStore) Exist(key string) bool {
    if _, ok := (*MS)[key]; !ok {
        return false;
    }
    return true;
}

func (MS *MemStore) Len() int {
    return len(*MS)
}

func (MS *MemStore) Reset() {
    *MS = make(map[string]Entry)
}

func (MS *MemStore) Remove(key string) bool {
    args := *MS
    if _, find := args[key]; !find {
        return false
    } else {
        delete(args, key)
        return true
    }
}

func (MS *MemStore) Visit(visitor func(key string, value interface{})) {
    for key, value := range *MS {
        visitor(key, value)
    }
}

func (MS *MemStore) clear(key string, expire int, timestamp int64) {
    timer := time.Duration(expire)
    time.AfterFunc(time.Second * timer, func() {
        if _, ok := (*MS)[key]; ok {
            if (*MS)[key].version == timestamp {
                MS.Remove(key)
            }
        }
    })
}

func (MS *MemStore) save(key string, value interface{}, expire int, immutable bool) error {
    if expire >= 0 {
        if len(*MS) == 0 {
            *MS = make(map[string]Entry)
        }
        tm := time.Now().UnixNano()
        if entry, find := (*MS)[key]; find {
            if entry.immutable {
                return fmt.Errorf("this key(%s) write protection", key)
            }
        }
        (*MS)[key] = Entry{
            value: value,
            version: tm,
            ttl: tm + int64(expire) * 1e9,
            immutable: immutable,
        }
        if expire > 0 {
            MS.clear(key, expire, tm)
        }
    }
    return nil
}

func (MS *MemStore) Set(key string, value interface{}, expire int) error {
    return MS.save(key, value, expire, false)
}

func (MS *MemStore) SetImmutable(key string, value interface{}, expire int) error {
    return MS.save(key, value, expire, true)
}

func (MS *MemStore) Get(key string) (interface{}, error) {
    args := *MS
    if entry, find := args[key]; find {
        return entry.Value(), nil
    } else {
        return nil, fmt.Errorf("can not find value for %s", key)
    }
}

func (MS *MemStore) GetString(key string) (string, error) {
    if v, err := MS.Get(key); err != nil {
        return "", err
    } else {
        if value, ok := v.(string); !ok {
            return "", fmt.Errorf("can not convert %#v to string", v)
        } else {
            return value, nil
        }
    }
}

func (MS *MemStore) GetInt(key string) (int, error) {
    v, _ := MS.Get(key)
    if vInt, ok := v.(int); ok {
        return vInt, nil
    }
    if vString, ok := v.(string); ok {
        return strconv.Atoi(vString)
    }
    return -1, errors.New(fmt.Sprintf("unable to find or parse the integer, found: %#v", v))
}

func (MS *MemStore) TTL(key string) float64 {
    if _, ok := (*MS)[key]; !ok {
        return -1
    }
    return float64(((*MS)[key].ttl - time.Now().UnixNano()) / 1e9)
}