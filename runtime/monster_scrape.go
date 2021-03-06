package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	osrs_api "github.com/cloakd/osrs-api"
	"github.com/cloakd/osrs-api/cmd"
)

const (
	CACHE_FILE = "cache/monster_cache.json"
)

func main() {
	start := time.Now()

	monMap, err := cmd.CacheMonsters()
	if err != nil {
		log.Fatal(err)
	}

	mWrite := make(chan *osrs_api.Monster, 100)
	go func(mChan chan *osrs_api.Monster) {
		log.Printf("Opening file: %s", CACHE_FILE)
		f, err := os.Create(CACHE_FILE)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		for m := range mChan {
			j, err := json.Marshal(m)
			if err != nil {
				log.Fatal(err)
			}

			_, err = f.Write(j)
			if err != nil {
				log.Fatal(err)
			}

			f.WriteString("\n")
		}
	}(mWrite)

	ms := cmd.NewMonsterScraper(&monMap)
	ms.Scrape(mWrite)

	close(mWrite)

	fin := time.Since(start)
	log.Printf("Monster scrape took %s", fin)
}
