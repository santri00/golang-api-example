package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/santri00/golang-api-example/pkg/api"
	"github.com/santri00/golang-api-example/pkg/apiconfig"
	"github.com/santri00/golang-api-example/pkg/apiconfig/configkeys"
	"github.com/santri00/golang-api-example/pkg/services/vehicle"
	"log"
)

// initializes the server and dependencies
func initialize() *api.Server {
	c := apiconfig.New()
	e := c.Load("config.yaml")
	if e != nil {
		log.Fatalf("failed to load the api configuration. error: %v", e.Error())
	}

	rdb := initializeRedisClient(fmt.Sprintf("%s:%s",
		c.Store.Get(configkeys.RedisAddr),
		c.Store.Get(configkeys.RedisPort)))
	defer rdb.Close()

	return api.NewServer(context.Background(),
		c,
		rdb,
		vehicle.NewService())
}

// initializes the redis client and connection
func initializeRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // TODO
		DB:       0,  // use default DB
	})
}
