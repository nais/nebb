package storage

import (
	"github.com/gomodule/redigo/redis"
	rg "github.com/redislabs/redisgraph-go"
	"log"
	"os"
	"time"
)

type Datastore struct {
	graph rg.Graph
}

type App struct {
	Name    string
	Version string
	Libs    []Library
}

type Library struct {
	Name    string
	Version string
}

func redisAddress() string {
	address, exists := os.LookupEnv("REDIS_ADDRESS")
	if !exists {
		log.Fatal("env var REDIS_ADDRESS must be set")
	}
	return address
}

func redisPassword() string {
	password, exists := os.LookupEnv("REDIS_PASSWORD")
	if !exists {
		log.Fatal("env var REDIS_PASSWORD must be set")
	}
	return password
}

var connectionPool = redis.Pool{
	MaxIdle:     3,
	IdleTimeout: 240 * time.Second,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp",
			redisAddress(),
			redis.DialPassword(redisPassword()))
	},

	TestOnBorrow: func(c redis.Conn, t time.Time) error {
		_, err := c.Do("PING")
		return err
	},
}

func NewDatastore() (*Datastore, error) {
	conn, err := connectionPool.Dial()
	if err != nil {
		return &Datastore{}, err
	}

	return &Datastore{
		graph: rg.GraphNew("dependencies", conn),
	}, nil
}

func (d *Datastore) Upsert(app App) error {
	for _, lib := range app.Libs {
		params := make(map[string]interface{})
		params["appname"] = app.Name
		params["appversion"] = app.Version
		params["libname"] = lib.Name
		params["libversion"] = lib.Version
		_, err := d.graph.ParameterizedQuery(`MERGE (a:APP {name: $appname, version: $appversion}) MERGE (l:LIBRARY {name: $libname, version: $libversion}) MERGE (a)-[u:USES]->(l)`, params)
		if err != nil {
			return err
		}
	}
	return nil
}
