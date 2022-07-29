package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
	"github.com/jaydenpung/pony-tally-bot/internal/model"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")

	requestData := model.Payload{
		Query: "query GovernanceProposals($sort: ProposalSort, $chainId: ChainID!, $pagination: Pagination, $governanceIds: [AccountID!], $proposerIds: [AccountID!]) {   proposals(     sort: $sort     chainId: $chainId     pagination: $pagination     governanceIds: $governanceIds     proposerIds: $proposerIds   ) {     id     description     voteStats {       votes       weight       support       percent     }   } } ",
		Variables: model.Variables{
			Pagination: model.Pagination{
				Limit:  4,
				Offset: 0,
			},
			Sort: model.Sort{
				Field: "START_BLOCK",
				Order: "DESC",
			},
			ChainID:       "eip155:1",
			GovernanceID:  "eip155:1:0x6CC90C97a940b8A3BAA3055c809Ed16d609073EA",
			GovernanceIds: []string{"eip155:1:0x6CC90C97a940b8A3BAA3055c809Ed16d609073EA"},
		},
	}
	payloadBytes, err := json.Marshal(requestData)
	if err != nil {
		log.Fatal(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api.tally.xyz/query", body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	tallyResponse := new(model.TallyResponse)

	if err := json.NewDecoder(resp.Body).Decode(tallyResponse); err != nil {
		log.Fatalln("Error converting json: ", err)
	}

	//fmt.Printf("response: %v", tallyResponse)

	db := createBucket()
	for _, proposal := range tallyResponse.Data.Proposals {
		existingProposal := readProposal(db, proposal.Id)
		if len(existingProposal.Id) > 0 {
			fmt.Printf("Proposal exist: %v\n", existingProposal.Id)
		} else {
			message := proposal.GenerateMessage()
			sendDiscordMessage(message)
			saveProposal(db, &proposal)
			fmt.Printf("Saved proposal: %v\n", proposal.Id)
		}
	}

	defer db.Close()
}

func createBucket() *bolt.DB {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		if err != nil {
			log.Fatalf("create bucket: %s", err)
		}
		return nil
	})

	return db
}

func saveProposal(db *bolt.DB, proposal *model.Proposal) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		data, err := json.Marshal(proposal)
		if err != nil {
			log.Fatal(err)
		}
		err = b.Put([]byte(proposal.Id), []byte(data))
		return err
	})
}

func readProposal(db *bolt.DB, id string) *model.Proposal {
	var proposal model.Proposal
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		v := b.Get([]byte(id))

		if len(v) > 0 {
			err := json.Unmarshal(v, &proposal)
			if err != nil {
				log.Fatal(err)
			}
		}

		return nil
	})

	return &proposal
}

func sendDiscordMessage(message string) {

	type DiscordData struct {
		Content string `json:"content"`
	}
	discordWebHookUrl := os.Getenv("DISCORD_WEB_HOOK_URL")
	if len(discordWebHookUrl) < 1 {
		log.Fatalf("Please set environment variable DISCORD_WEB_HOOK_URL")
	}

	payloadBytes, err := json.Marshal(DiscordData{
		Content: message,
	})
	if err != nil {
		log.Fatalf("Error forming discord web hook data: %v", err)
	}
	body := bytes.NewReader(payloadBytes)

	http.Post(discordWebHookUrl, "application/json", body)
}
