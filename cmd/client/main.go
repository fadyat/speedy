package main

import (
	"fmt"
	"github.com/fadyat/speedy/client"
	"log"
	"time"
)

// todo: delete me and write in tests
//  for fast testing I up servers in docker compose
//  and here connecting to them

func main() {
	dcl, err := client.NewClient("cmd/client/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		key, value := fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i)
		if e := dcl.Put(key, value); e != nil {
			log.Println(e)
		} else {
			log.Println("put", key, value)
		}

		time.Sleep(100 * time.Millisecond)
	}

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key%d", i)
		if v, e := dcl.Get(key); e != nil {
			log.Println(e)
		} else {
			log.Println("get", v)
		}

		time.Sleep(100 * time.Millisecond)
	}

	// checking cache miss response
	if v, e := dcl.Get("aboba"); e == nil {
		log.Println("cache miss:", v)
		return
	} else {
		log.Println("unknown error:", e)
	}
}
