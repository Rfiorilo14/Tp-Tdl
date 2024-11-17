package shared

type Message struct {
	Type      string `json:"type"`
	PlayerID  string `json:"player_id"`
	Direction string `json:"direction"`
}
