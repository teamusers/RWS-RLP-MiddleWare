package topic

import (
	"context"
	"encoding/json"
	"runtime"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stonksdex/externalapi/codes"
	"github.com/stonksdex/externalapi/log"
	"github.com/stonksdex/externalapi/model"
	"github.com/stonksdex/externalapi/system"
	"github.com/stonksdex/externalapi/wsm"
)

var lastPrintLogStat = 0
var zero = decimal.NewFromInt(0)

var tokenStatCache = make(map[string]*model.StatWalletToken)
var cacheLock sync.Mutex

func init() {
	go batchUpdateWalletStats()
}

func ConsumeFlow() {
	pubsub := system.GetRedis().Subscribe(context.Background(), codes.TOPIC_TOKEN_FLOW_SUB)
	ch := pubsub.Channel()

	for msg := range ch {
		var jsonFlow SolTokenFlow
		err := json.Unmarshal([]byte(msg.Payload), &jsonFlow)
		if err != nil {
			log.Errorf("[TokenFlowSub] TokenFlow Data parse error: %s, %v", msg.Payload, err)
			continue
		}

		if jsonFlow.Amount0In > 0 {
			jsonFlow.TradeDirection = "sell"
		} else {
			jsonFlow.TradeDirection = "buy"
		}

		strFlow, err := json.Marshal(jsonFlow)
		if err != nil {
			log.Error("[TokenFlowSub] json object error: ", err)
			continue
		}

		// updateTokenWalletStat(&jsonFlow)
		go updateTxDiagram(jsonFlow.Tx)

		wsmanager := wsm.RetrieveWsManager()
		wsmanager.Broadcast("solana", jsonFlow.Token0, string(strFlow))
		all, audit := wsmanager.Stat()
		if time.Now().Second()-lastPrintLogStat >= 60 {
			if all > 0 || audit > 0 {
				log.Infof("current connect audition result: total:%d, align:%d", all, audit)
			}
			lastPrintLogStat = time.Now().Second()
		}
	}
}

func updateTxDiagram(txhash string) {
	redis := system.GetRedis()
	if redis == nil {
		return
	}

	exist, err := redis.Exists(context.Background(), "txnotify:"+txhash).Result()
	if err != nil {
		log.Error("[txnotify] fetch tx from redis error: ", txhash, err)
		return
	} else {
		if exist == 0 {
			// log.Info("[txnotify] hash not exist: ", txhash)
			return
		}
	}

	db := system.GetDb()
	var tx model.TxsDiagram
	db.Model(&model.TxsDiagram{}).Where("chain = ? and tx_hash = ?", "solana", txhash).First(&tx)
	if tx.ID > 0 {
		tx.UpdateTime = time.Now()
		tx.Status = "20"
		db.Save(&tx)
		log.Infof("update txdiagram tx: %s", txhash)
	}
}

func updateTokenWalletStat(tf *SolTokenFlow) {
	defer func() {
		if r := recover(); r != nil {
			stackBuf := make([]byte, 4096)
			n := runtime.Stack(stackBuf, false)
			stackTrace := string(stackBuf[:n])
			log.Errorf("[updateTokenWalletStat error] %s", stackTrace)
		}
	}()
	if tf == nil {
		return
	}

	// **构造 key**
	cacheKey := tf.Payer + "_" + tf.Token0

	cacheLock.Lock()
	defer cacheLock.Unlock()

	statwallet, exists := tokenStatCache[cacheKey]
	if !exists {
		statwallet = &model.StatWalletToken{
			Chain:  "solana",
			Wallet: tf.Payer,
			Token:  tf.Token0,
		}
		tokenStatCache[cacheKey] = statwallet
	}

	if tf.Amount0In > 0 { // sell transaction
		if tf.Amount0In <= statwallet.Balance {
			statwallet.Balance = statwallet.Balance - tf.Amount0In
		}
		statwallet.CountSell += 1
		statwallet.TotalSell += tf.Amount0In
		statwallet.LastDirect = "sell"
		if tf.Price.Cmp(zero) == 1 && tf.Token0Decimals > 0 {
			statwallet.TotalSellValue = statwallet.TotalSellValue.Add(decimal.NewFromUint64(tf.Amount0In).Mul(tf.Price).Div(decimal.NewFromUint64(uint64(tf.Token0Decimals))))
		}
	} else { // buy transaction
		statwallet.Balance += tf.Amount0Out
		statwallet.CountBuy += 1
		statwallet.TotalBuy += tf.Amount0Out
		if statwallet.FirstTx == "" {
			statwallet.FirstTs = tf.TxTime
			statwallet.FirstTx = tf.Tx
		}
		statwallet.LastDirect = "buy"
		if tf.Price.Cmp(zero) == 1 && tf.Token0Decimals > 0 {
			statwallet.TotalBuyValue = statwallet.TotalBuyValue.Add(decimal.NewFromUint64(tf.Amount0Out).Mul(tf.Price).Div(decimal.NewFromUint64(uint64(tf.Token0Decimals))))
		}
	}
	statwallet.LastTs = tf.TxTime
	statwallet.LastTx = tf.Tx
}

func batchUpdateWalletStats() {
	for {
		time.Sleep(1 * time.Second)

		cacheLock.Lock()
		if len(tokenStatCache) == 0 {
			cacheLock.Unlock()
			continue
		}

		updates := make([]*model.StatWalletToken, 0, len(tokenStatCache))
		for key, stat := range tokenStatCache {
			updates = append(updates, stat)
			delete(tokenStatCache, key)
		}
		tokenStatCache = make(map[string]*model.StatWalletToken)
		cacheLock.Unlock()

		db := system.GetDb()
		tx := db.Begin()
		for _, stat := range updates {
			tx.Save(stat)
		}
		tx.Commit()
	}
}
