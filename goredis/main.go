package main

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-redis/redis/v8"
)

type Model struct {
	Str1    string   `redis:"str1"`
	Str2    string   `redis:"str2"`
	Int     int      `redis:"int"`
	Bool    bool     `redis:"bool"`
	Ignored struct{} `redis:"-"`
}

func main() {
	ctx := context.Background()

	rdb := redis.NewClusterClient(&redis.ClusterOptions{Addrs: []string{"127.0.0.1:6380", "127.0.0.1:6381", "127.0.0.1:6382", "127.0.0.1:6383", "127.0.0.1:6384", "127.0.0.1:6385"}})

	if _, err := rdb.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, "key", "str1", "hello")
		rdb.HSet(ctx, "key", "str2", "hello")
		rdb.HSet(ctx, "key", "int", 123)
		return nil
	}); err != nil {
		panic(err)
	}

	var model1, model2 Model

	if err := rdb.HGetAll(ctx, "key").Scan(&model1); err != nil {
		panic(err)
	}

	if err := rdb.HMGet(ctx, "key", "str1", "int").Scan(&model2); err != nil {
		panic(err)
	}

	spew.Dump(model1)
	spew.Dump(model2)
}
