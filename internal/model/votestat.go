package model

type VoteStat struct {
	Votes   string  `json:"votes"`
	Weight  string  `json:"weight"`
	Support string  `json:"support"`
	Percent float32 `json:"percent"`
}
