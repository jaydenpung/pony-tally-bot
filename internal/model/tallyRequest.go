package model

type Payload struct {
	Query     string    `json:"query"`
	Variables Variables `json:"variables"`
}

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type Sort struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

type Variables struct {
	Pagination    Pagination `json:"pagination"`
	Sort          Sort       `json:"sort"`
	ChainID       string     `json:"chainId"`
	GovernanceID  string     `json:"governanceId"`
	GovernanceIds []string   `json:"governanceIds"`
}
