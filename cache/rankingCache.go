package mycache

import (
	"time"

	"rlp-member-service/system"

	"github.com/dgraph-io/ristretto/v2"
)

var rankingUpdateCache *ristretto.Cache[string, time.Time]
var db = system.GetRedis()

func init() {
	cache12, err := ristretto.NewCache[string, time.Time](&ristretto.Config[string, time.Time]{
		NumCounters: 3000,      // number of keys to track frequency of (10M).
		MaxCost:     200000000, // maximum cost of cache (1GB).
		BufferItems: 64,        // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}
	time.Now().Unix()
	rankingUpdateCache = cache12
}

func RankingCacheShouldUpdate(key string, t time.Duration) bool {
	rankingUpdateCache.Wait()
	s, b := rankingUpdateCache.Get(key)
	flag := false
	if !b || s.Before(time.Now().Add(-t)) {
		flag = true
		rankingUpdateCache.SetWithTTL(key, time.Now(), 0, t)
	}

	return flag
}
