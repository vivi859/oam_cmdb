package util

import (
	conf "OAM/conf"
	"context"
	"errors"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/beego/beego/v2/client/cache"
)

const (
	CACHE_LOGIN_FAILED = "loginfailedcache"
	CACHE_DICT         = "dictitemcache"
	CACHE_COMMONS      = "_public_object"
	CACHE_FUN_GROUP    = "fun_group"
)

var cacheMap map[string]MyCache

//继承beego cache,简化缓存操作
type MyCache struct {
	cache.Cache
	//过期清理间隔,单位秒
	interval int
	//默认过期时间
	ExpireTime time.Duration
	*sync.RWMutex
}

func init() {
	cacheMap = make(map[string]MyCache)
	//先定义要用到的缓存配置, 具体缓存对象延迟初始化
	loginFailedCache := MyCache{interval: 30, ExpireTime: time.Minute * 10, RWMutex: &sync.RWMutex{}}
	cacheMap[CACHE_LOGIN_FAILED] = loginFailedCache

	dictCache := MyCache{interval: 60, ExpireTime: time.Hour * 2, RWMutex: &sync.RWMutex{}}
	cacheMap[CACHE_DICT] = dictCache

	commonsCache := MyCache{interval: 60, ExpireTime: time.Hour * 24, RWMutex: &sync.RWMutex{}}
	cacheMap[CACHE_COMMONS] = commonsCache

	funGroupCache := MyCache{interval: 60, ExpireTime: time.Hour * 24, RWMutex: &sync.RWMutex{}}
	cacheMap[CACHE_FUN_GROUP] = funGroupCache
}

func (m MyCache) Get(key string) (interface{}, error) {
	return m.Cache.Get(context.TODO(), key)
}

// 从缓存中取string值,如果不存在返回空字符串,请确保值一定是string类型
func (m MyCache) GetString(key string) string {
	val, err := m.Cache.Get(context.TODO(), key)
	if err != nil {
		return ""
	}
	return val.(string)
}

// 从缓存中取int类型值,如果不存在返回0,int
func (m MyCache) GetInt(key string) int {
	val, err := m.Cache.Get(context.TODO(), key)
	if err != nil {
		return 0
	}
	return val.(int)
}

// 缓存对象,并设置过期时间
func (m MyCache) PutWithExpireTime(key string, value interface{}, timeout time.Duration) {
	m.Cache.Put(context.TODO(), key, value, timeout)
}

//缓存对象,过期时间为缓存配置的默认值
func (m MyCache) Put(key string, value interface{}) {
	m.Cache.Put(context.TODO(), key, value, m.ExpireTime)
}

func (m MyCache) Incr(key string) error {
	return m.Cache.Incr(context.TODO(), key)
}

func (m MyCache) Delete(key string) {
	m.Cache.Delete(context.TODO(), key)
}

func (m MyCache) ClearAll() {
	m.Cache.ClearAll(context.TODO())
}

// 创建一个缓存. 参数interval表示清理缓存间隔时间
func CreateCache(name string, interval int, expireTime time.Duration) MyCache {
	var mycache MyCache
	c, err := cache.NewCache(conf.GlobalCfg.CACHE_TYPE, `{"interval":`+strconv.Itoa(interval)+`}`)
	if err != nil {
		panic(errors.New("缓存初始失败"))
	}
	mycache = MyCache{Cache: c, interval: interval, ExpireTime: expireTime}
	return mycache
}

func GetCache(name string) MyCache {
	mc, ok := cacheMap[name]
	if ok {
		if mc.Cache == nil {
			mc.Lock()
			if mc.Cache == nil {
				c, err := cache.NewCache(conf.GlobalCfg.CACHE_TYPE, `{"interval":`+strconv.Itoa(mc.interval)+`}`)
				if err != nil {
					mc.Unlock()
					panic(errors.New("缓存初始失败"))
				}
				mc.Cache = c
				cacheMap[name] = mc
			}
			mc.Unlock()
		}
		return mc
	} else {
		panic("试图获取未定义的缓存")
	}
}

func QueryCacheFirst2[T any](cacheName string, cacheKey string, queryFunc func() []T) []T {
	return QueryCacheFirst(GetCache(cacheName), cacheKey, queryFunc)
}

// 数据查询,优先从缓存中取数据,如果缓存不存在,则调用queryFunc函数取数据,并缓存
func QueryCacheFirst[T any](cache MyCache, cacheKey string, queryFunc func() T) T {
	data, err := cache.Get(cacheKey)
	var ts T
	isFetch := true
	if err == nil {
		ts = data.(T)
		if !IsZeroRef[T](ts) {
			isFetch = false
		}
	}
	if isFetch {
		ts = queryFunc()
		if !IsZeroRef[T](ts) {
			cache.Put(cacheKey, ts)
		}
	}
	return ts
}

// 判断一个变量是否为零值
func IsZeroRef[T any](v T) bool {
	return reflect.ValueOf(&v).Elem().IsZero()
}

func GetPublicCache() MyCache {
	return GetCache(CACHE_COMMONS)
}

func GetFunGroupCache() MyCache {
	return GetCache(CACHE_FUN_GROUP)
}
