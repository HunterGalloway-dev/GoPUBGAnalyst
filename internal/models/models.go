package models

type Bunt struct {
	Message string `json:"message"`
}

type MatchInfo struct {
	MapName    string
	CircleInfo map[float64]Location
}

// TELEM EVENT START
type TelemEvent struct {
	Type  string    `json:"_T"`
	C     Common    `json:"common"`
	State GameState `json:"gameState"`
}

type GameState struct {
	SafetyZonePos Location `json:"safetyZonePosition"`
}

type Location struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Common struct {
	IsGame float64 `json:"isGame"`
}

// TELEM EVENT END

// MATCH RESPONSE /matches:id START
type MatchResponse struct {
	Data     MatchData          `json:"data"`
	Included []IncludedResponse `json:"included"`
}

type MatchData struct {
	Attributes MatchAttributes `json:"attributes"`
}

type MatchAttributes struct {
	// Todo add other attrs for future filter
	MapName string `json:"mapName"`
}

type IncludedResponse struct {
	Type       string             `json:"type"`
	Attributes IncludedAttributes `json:"attributes"`
}

type IncludedAttributes struct {
	Name string `json:"name"`
	URL  string `json:"URL"`
}

// MATCH RESPONSE END

// TO DO SEPERATE EACH RESPONSE INTO SEPERATE FILE UNDER SAME PACKAGE
// PLAYER RESPONSE FOR /players endpoint START
type PlayersResponse struct {
	Data []Player `json:"data"`
}

type Player struct {
	Type          string              `json:"type"`
	Id            string              `json:"id"`
	Relationships PlayerRelationships `json:"relationships"`
}

type PlayerRelationships struct {
	Matches PlayerMatches `json:"matches"`
}

type PlayerMatches struct {
	Data []PlayerMatch `json:"data"`
}

type PlayerMatch struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

// PLAYER RESPONSE FOR /players endpoint END
