package topic

import "github.com/stonksdex/externalapi/log"

func init() {
	// go ConsumeToken()
	go ConsumeFlow()
}

func StartSubscription() {
	log.Info("[Sub] system started")
}
