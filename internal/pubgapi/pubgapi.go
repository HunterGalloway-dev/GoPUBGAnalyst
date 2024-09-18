package pubgapi

import (
	"GoPUBGAnalyst/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

// Max allowed # of players per player search, used for optimization
const BATCHSIZE = 10

type PUBGApi struct {
	client   http.Client
	apiToken string
}

func (p *PUBGApi) BuildCircleMapFromTelemUrl(telemURL string) (map[float64]models.Location, error) {
	response, err := p.Get(telemURL)
	circleMap := make(map[float64]models.Location)

	if err != nil {
		return circleMap, err
	}

	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		return circleMap, err
	}

	var responseObject []models.TelemEvent
	json.Unmarshal(responseData, &responseObject)

	for _, e := range responseObject {
		if e.Type == "LogGameStatePeriodic" {
			circleMap[e.C.IsGame] = e.State.SafetyZonePos
		}
	}
	return circleMap, nil
}

func (p *PUBGApi) BuildMatchInfo(matchId string) (*models.MatchInfo, error) {
	response, err := p.Get(fmt.Sprintf("https://api.pubg.com/shards/steam/matches/%s", matchId))
	matchInfo := new(models.MatchInfo) // Bader use of memory usage

	if err != nil {
		return matchInfo, err
	}

	responseData, err := io.ReadAll(response.Body) // DRY -> Make function you're repeating yourself

	if err != nil {
		return matchInfo, err
	}

	var responseObject models.MatchResponse
	json.Unmarshal(responseData, &responseObject)
	telemURL := ""

	for _, v := range responseObject.Included {
		if v.Type == "asset" && v.Attributes.Name == "telemetry" {
			telemURL = v.Attributes.URL
		}
	}

	matchInfo.MapName = responseObject.Data.Attributes.MapName
	circleMap, err := p.BuildCircleMapFromTelemUrl(telemURL)

	if err != nil {
		return matchInfo, err
	}
	matchInfo.CircleInfo = circleMap

	return matchInfo, err
}

// Move to lib/utils.go later not now, Im tryna get it working first!
func Batcher(input []string, size int) [][]string {
	var ret [][]string

	arr := input

	for len(arr) >= size {
		// Append slice to ret
		ret = append(ret, arr[0:size])
		arr = arr[size:]
	}

	if len(arr) > 0 {
		ret = append(ret, arr)
	}

	return ret
}

func (p *PUBGApi) GetMatchIdsFromPlayers(players []string) (map[string]bool, error) {
	matchIdUnion := make(map[string]bool)
	batches := Batcher(players, BATCHSIZE)

	for _, b := range batches {
		ids, err := p.getMatchIdsFromPlayerBatch(b)

		if err != nil {
			return matchIdUnion, err
		}

		for k, _ := range ids {
			matchIdUnion[k] = true
		}
	}

	return matchIdUnion, nil
}

// To do add error when length of players exceeds 10
func (p *PUBGApi) getMatchIdsFromPlayerBatch(players []string) (map[string]bool, error) {
	matchIdSet := make(map[string]bool)

	response, err := p.Get(fmt.Sprintf("https://api.pubg.com/shards/steam/players?filter[playerNames]=%s", strings.Join(players, ",")))

	if err != nil {
		return matchIdSet, err
	}

	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		return matchIdSet, err
	}

	var responseObject models.PlayersResponse
	json.Unmarshal(responseData, &responseObject)

	for _, p := range responseObject.Data {
		for _, m := range p.Relationships.Matches.Data {
			matchIdSet[m.Id] = true
		}
	}

	return matchIdSet, nil
}

func (p *PUBGApi) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return p.Do(req)
}

func (p *PUBGApi) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.apiToken))
	req.Header.Add("Accept", "application/vnd.api+json")

	return p.client.Do(req)
}

func NewPUBGApi() *PUBGApi {
	return &PUBGApi{
		apiToken: os.Getenv("PUBG_API_KEY"),
	}
}
