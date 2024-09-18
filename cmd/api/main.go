package main

import (
	"GoPUBGAnalyst/internal/models"
	"GoPUBGAnalyst/internal/pubgapi"
	"fmt"
	"log"
	"time"
)

func main() {

	playerList := [...]string{"BuntStreams", "TGLTN", "rogiw0w", "CowBoixx", "highegoplayerr", "SneakAttack", "NordSkif", "shane_doe", "DontJuiceMe", "Voxsic", "xnnnyy", "TATER-_-"}

	apiService := pubgapi.NewPUBGApi()
	matchIdsMap, err := apiService.GetMatchIdsFromPlayers(playerList[:])
	var matchIds []string

	for k, _ := range matchIdsMap {
		matchIds = append(matchIds, k)
	}

	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(len(matchIds))

	if err != nil {
		log.Fatal(err.Error())
	}

	freqMap := make(map[string]int)
	B_SIZE := 50
	fmt.Println(B_SIZE)
	batches := pubgapi.Batcher(matchIds, B_SIZE)

	startTime := time.Now()
	for _, batch := range batches {

		c := make(chan *models.MatchInfo, len(batch))
		for i := 0; i < len(batch); i++ {
			go buildMatchInfoWorker(matchIds[i], c, apiService)
		}

		for i := 0; i < len(batch); i++ {
			mi := <-c

			val, ok := freqMap[mi.MapName]

			if ok {
				freqMap[mi.MapName] = val + 1
			} else {
				freqMap[mi.MapName] = 1
			}
		}

	}
	duration := time.Since(startTime)
	fmt.Println("Finished batch operation")
	fmt.Println(duration)

	fmt.Println(freqMap)
}

func buildMatchInfoWorker(matchID string, c chan *models.MatchInfo, apiServce *pubgapi.PUBGApi) {
	mi, _ := apiServce.BuildMatchInfo(matchID)
	c <- mi
}
