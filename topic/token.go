package topic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/stonksdex/externalapi/codes"
	"github.com/stonksdex/externalapi/fullindex"
	"github.com/stonksdex/externalapi/log"
	"github.com/stonksdex/externalapi/model"
	"github.com/stonksdex/externalapi/system"
	"github.com/stonksdex/externalapi/thirdpart"
	"gorm.io/gorm"
)

var solanaTokenQueue = system.NewRichQueue[string]()

func ConsumeToken() {
	pubsub := system.GetRedis().Subscribe(context.Background(), codes.TOPIC_TOKEN_SUB)
	ch := pubsub.Channel()
	go parseTokens(50)
	for msg := range ch {
		var addrs []string
		err := json.Unmarshal([]byte(msg.Payload), &addrs)
		if err != nil {
			log.Errorf("[TokenSub] Token Data Parse Error: %v", err)
			continue
		}

		solanaTokenQueue.BatchEnqueue(addrs)
	}
}

func parseTokens(size int) {
	mdb := system.GetDb()
	chain := "solana"
	solanaTokenQueue.Consumer(size, func(ca string, lock *sync.WaitGroup) {
		var existMeta model.TokenMeta
		err := mdb.Model(&model.TokenMeta{}).Where("chain = ? and ca = ?", chain, ca).First(&existMeta).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("[TokenSub] Unexpected:", err)
		}
		if existMeta.ID > 0 {
			return
		}

		result, err := thirdpart.GetMetaInfoFromGoPlus(ca)
		if err != nil {
			log.Error("[TokenSub] GOPLUS getMetaInfoFromGoPlus failed, and still insert it into Queue: ", ca, err)
			solanaTokenQueue.Enqueue(ca)
			return
		}

		goplusToken := result.Result[ca]
		jsonGoplusToken, err := json.Marshal(goplusToken)
		if err == nil {
			// log.Infof("[TokenSub] GOPLUS getMetaInfoFromGoPlus success; %s, %s", ca, string(jsonGoplusToken))
			_ = jsonGoplusToken
		} else {
			log.Infof("[TokenSub] GOPLUS getMetaInfoFromGoPlus success; %s, %v", ca, goplusToken)
		}
		existMeta, existPairs := thirdpart.CopyMetaIntoObject(goplusToken)
		existMeta.CA = ca
		existMeta.Chain = chain

		existMeta.CreateTime = time.Now()
		err = mdb.Save(&existMeta).Error

		if err == nil {
			log.Info("[TokenSub][Search] create index for token ", existMeta)
			fullindex.AddToIndex(fmt.Sprintf("token-%d", existMeta.ID), map[string]interface{}{
				"ca":     existMeta.CA,
				"name":   existMeta.Name,
				"symbol": existMeta.Symbol,
			})
		}
		for _, vs := range existPairs {
			vs.TokenID = existMeta.ID
			mdb.Save(&vs)
		}
		time.Sleep(3 * time.Second)
	})
}
