package model

type TallyResponse struct {
	Data Data `json:"data"`
}

type Data struct {
	Proposals []Proposal `json:"proposals"`
}
