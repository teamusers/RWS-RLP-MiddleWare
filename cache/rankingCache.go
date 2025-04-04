package mycache

import (
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/stonksdex/externalapi/log"
	"github.com/stonksdex/externalapi/model"
	"github.com/stonksdex/externalapi/system"
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
func RankingCacheGet(key string) (*[]model.RankingItem, error) {

	get, err := system.ObjectGet(key, &[]model.RankingItem{})
	return get, err
}
func RankingCacheSet(key string, value []model.RankingItem, expire time.Duration) {

	err := system.ObjectSet(key, value, expire)
	if err != nil {
		log.Errorf("RankingCacheSet cache error:%s", err.Error())
	}

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
