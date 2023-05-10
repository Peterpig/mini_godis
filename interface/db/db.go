package db

import "github.com/Peterpig/mini_godis/interface/redis"

type DB interface {
	Exec(client redis.Client, args []byte) redis.Reply
	AfterClientClose(client redis.Client)
	Close()
}
