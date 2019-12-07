package gocache

import (
	alg "four-seasons/algorithm"
	"sync"
	"time"
)

type item struct {
	obj interface{}
	expiration int64
}

var  (
	NoExpiration  = time.Second * 0
	DefaultExpiration = time.Hour * 1
	DefaultCapacity = 2048
	DefaultCommonsChannelSize = 4096
	DefaultCleatStep = time.Second * 1
)


type CacheType string

var CacheForLFU CacheType = "Type_LFU"
var CacheForLRU CacheType = "Type_LRU"

var once sync.Once

type CacheManager struct {
	mu sync.RWMutex
	m map[string]*Cache
}

type SignalCache Cache

type Cache struct {
	cache      CacheInter
	dp         dispatcher
	expiration time.Duration
	liquidator *Liquidator
	isClose bool
}

type Liquidator struct {
	elements *alg.PriorityQueue
}

type commons struct {
	fn func()
	setFn func( string,  *item , *Cache, chan interface{}) interface{}
}

type dispatcher struct {
	queue chan commons
	stateCh chan string
}


type CacheInter interface {
	Get(string) interface{}
	Put(string, interface{}) interface{}
	Delete(string) interface{}
}

type strategy interface {
	update(*item)
	isClear(*item) bool
}


/*************************************************************
					LFU
 *************************************************************/
type lfuNode struct {
	prev   *lfuNode
	next   *lfuNode
	parent *lfuNodeChain
	key string
	data interface{}
	freq int
}

type lfuNodeChain struct {
	prev *lfuNodeChain
	next *lfuNodeChain
	head *lfuNode
	tail *lfuNode
	freq int
}

type LFUCache struct {
	capacity int
	size int
	elements         map[string]*lfuNode
	manager lfuChainManager
	master *Cache
}
type lfuChainManager struct {
	firstLinkedList  *lfuNodeChain
	lastLinkedList   *lfuNodeChain
}


type needClearNode struct {
	masterCache *Cache
	key string
	expiration int64
}

type lruNode struct {
	key string
	data interface{}
	next *lruNode
	prev *lruNode
}

type lruNodeChain struct {
	head *lruNode
	tail *lruNode
}

type LRUCache struct {
	capacity int
	size int
	elements         map[string]*lruNode
	manager *lruNodeChain
	master *Cache
}


