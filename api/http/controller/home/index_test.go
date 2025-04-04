package home

import (
	"fmt"
	"strings"
	"testing"
)

func TestJoin(t *testing.T) {
	s := []string{
		"Grass7B4RdKfBCjTKgSqnXkqjwiGvQyFbuSCUJr3XXjs",
		"cbbtcf3aa214zXHbiAZQwf4122FBYbraNdFqgw4iMij",
		"LAYER4xPpTCb3QL8S9u41EAhAX7mhBn8Q6xMTwY2Yzc",
		"J1toso1uCk3RLmjorhTtrVwY9HJ7X8V9yYac6Y7kGCPn",
		"SonicxvLud67EceaEzCLRnMTBqzYUUYNr93DBkBdDES",
		"TNSRxcUxoT9xBG3de7PiJyTDYu7kskLqcpddxnEJAS6",
	}

	sql := fmt.Sprintf("insert into token_list_temp(token0) values%s", "('"+strings.Join(s, `'),('`)+"')")
	fmt.Println(sql)
}

//insert into token_list_temp(token0) values('Grass7B4RdKfBCjTKgSqnXkqjwiGvQyFbuSCUJr3XXjs'),cbbtcf3aa214zXHbiAZQwf4122FBYbraNdFqgw4iMij'),LAYER4xPpTCb3QL8S9u41EAhAX7mhBn8Q6xMTwY2Yzc'),J1toso1uCk3RLmjorhTtrVwY9HJ7X8V9yYac6Y7kGCPn'),SonicxvLud67EceaEzCLRnMTBqzYUUYNr93DBkBdDES'),TNSRxcUxoT9xBG3de7PiJyTDYu7kskLqcpddxnEJAS6'

// insert into token_list_temp(token0) values('Grass7B4RdKfBCjTKgSqnXkqjwiGvQyFbuSCUJr3XXjs'),('cbbtcf3aa214zXHbiAZQwf4122FBYbraNdFqgw4iMij'),('LAYER4xPpTCb3QL8S9u41EAhAX7mhBn8Q6xMTwY2Yzc'),('J1toso1uCk3RLmjorhTtrVwY9HJ7X8V9yYac6Y7kGCPn'),('SonicxvLud67EceaEzCLRnMTBqzYUUYNr93DBkBdDES'),('TNSRxcUxoT9xBG3de7PiJyTDYu7kskLqcpddxnEJAS6'
