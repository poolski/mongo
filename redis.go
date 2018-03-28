package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

func aaa() {
	rc := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	pong, err := rc.Ping().Result()
	fmt.Println(pong, err)
}
