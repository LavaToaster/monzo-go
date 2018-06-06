package monzo

import (
	"time"
	"encoding/json"
	"log"
)

type Owners struct {
	UserId        string `json:"user_id"`
	PreferredName string `json:"preferred_name"`
}

type Account struct {
	Id          string    `json:"id"`
	Closed      bool      `json:"closed"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Type        string    `json:"type"`
	Owners      []Owners  `json:"owners"`
}

type AccountsList struct {
	Accounts []Account `json:"accounts"`
}

func (c *Client) GetAccounts() []Account {
	resp, _ := c.Do("GET", "/accounts", nil)
	defer resp.Body.Close()

	var accounts AccountsList

	if err := json.NewDecoder(resp.Body).Decode(&accounts); err != nil {
		log.Println(err)

		return nil
	}

	return accounts.Accounts
}
