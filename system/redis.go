package system

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"rlp-member-service/config"

	redis "github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx = context.Background()

var Nil = redis.Nil

func init() {
	var conf = config.GetConfig()
	if conf.AllStart > 0 {
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf(conf.Redis.Host+":%d", conf.Redis.Port),
			Password: conf.Redis.Password,
			DB:       conf.Redis.Db,
		})
	}
}

func GetRedis() *redis.Client {
	return rdb
}
func ObjectSet(key string, value interface{}, expire time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshaling struct to JSON: %v", err)
	}

	err = rdb.Set(ctx, key, jsonData, expire).Err()
	if err != nil {
		return fmt.Errorf("error setting key: %v", err)
	}

	return nil
}

func ObjectGet[T any](key string, value *T) (*T, error) {
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting key: %v", err)
	}
	err = json.Unmarshal([]byte(val), value)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON to struct: %v", err)
	}
	return value, nil
}

func PublishTokenSearch(tokenSearch []byte) {
	go func() {
		cmd := rdb.Publish(ctx, "tokensearch", tokenSearch)
		val, err := cmd.Result()
		if err != nil {
			log.Println("publish token search error", val, err)
			return
		}

	}()

}
func PublishToChan(channel string, msg []byte) {
	go func() {
		cmd := rdb.Publish(ctx, channel, msg)
		val, err := cmd.Result()
		if err != nil {
			log.Println("publish token search error", val, err)
			return
		}
		//log.Printf("message published to %d subscribers\n", val)
	}()
}
func GetCacheListByIndex(key string, index int64) (string, error) {
	return rdb.LIndex(context.Background(), key, int64(index)).Result()
}

// Left-push the data into the list and set an expiration time
func SetCacheObjectListByLeft(key string, dataList interface{}, size int, expireByDays int) (int64, error) {
	// Trim the list first
	SetTrim(key, size)

	// Add data to the left side
	count, err := rdb.LPush(ctx, key, dataList).Result()
	if err != nil {
		return 0, err
	}

	// Set expiration time if expireByDays is greater than 0
	if expireByDays > 0 {
		_, err = rdb.Expire(ctx, key, time.Duration(expireByDays)*24*time.Hour).Result()
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

// Trim the list to the specified size
func SetTrim(key string, size int) error {
	listSize, err := rdb.LLen(ctx, key).Result()
	if err != nil {
		return err
	}

	if listSize > int64(size) {
		_, err := rdb.LTrim(ctx, key, 0, int64(size)-1).Result()
		if err != nil {
			return err
		}
	}
	return nil
}
func SetCacheObjectListByIndex(key string, index int64, value interface{}) {
	err := rdb.LSet(ctx, key, index, value).Err()
	if err != nil {
		fmt.Println("Error setting value:", err)
		return
	}
}

// Get hash fields from Redis and deserialize into a struct
func HGet(hashKey string, field string, result interface{}) error {
	// Get the value of a hash field
	data, err := rdb.HGet(ctx, hashKey, field).Result()
	if err != nil {
		if err == redis.Nil {
			// If the field does not exist, return nil instead of an error
			return nil
		}
		return fmt.Errorf("error getting hash field: %v", err)
	}

	// Deserialize JSON data into a struct instance
	err = json.Unmarshal([]byte(data), result)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON to struct: %v", err)
	}

	return nil
}

// Serialize the struct into JSON and write it to a Redis hash field
func HSet(hashKey string, field string, data interface{}) error {
	// Serialize the struct into a JSON string
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling struct to JSON: %v", err)
	}

	// Write JSON data to a Redis hash field
	err = rdb.HSet(ctx, hashKey, field, jsonData).Err()
	if err != nil {
		return fmt.Errorf("error setting hash field: %v", err)
	}

	return nil
}
func RedisExpire(key string, expireTime time.Duration) {
	err := rdb.Expire(ctx, key, expireTime).Err()
	if err != nil {
		log.Printf("RedisExpire,key:%s, err:%v ", key, err)
		return
	}
}
