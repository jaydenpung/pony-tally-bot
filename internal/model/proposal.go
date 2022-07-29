package model

import "fmt"

type Proposal struct {
	Id          string     `json:"id"`
	Description string     `json:"description"`
	VoteStats   []VoteStat `json:"voteStats"`
}

func (proposal *Proposal) GenerateMessage() string {
	message := fmt.Sprint(
		"Howdy <@&1001553848281866382>, a new proposal just went live.\n\n",
		"**Voting ends in 48 hours.**\n\n",
		"Click the link below and cast your votes.\n",
		"https://www.tally.xyz/governance/eip155:1:0x6CC90C97a940b8A3BAA3055c809Ed16d609073EA/proposal/",
		proposal.Id,
	)
	return message
}
