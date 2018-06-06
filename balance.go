package monzo

import (
	"encoding/json"
	"log"
)

type Spend struct {
	Currency   string `json:"currency"`
	SpendToday int64  `json:"spend_today"`
}

type Balance struct {
	*Spend
	Balance           int64   `json:"balance"`
	TotalBalance      int64   `json:"total_balance"`
	LocalCurrency     string  `json:"local_currency"`
	LocalExchangeRate int64   `json:"local_exchange_rate"`
	LocalSpend        []Spend `json:"local_spend"`
}

func (c *Client) GetBalance(accountId string) Balance {
	resp, _ := c.Do("GET", "/balance?account_id="+accountId, nil)
	defer resp.Body.Close()

	var balance Balance

	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
		log.Println(err)
	}

	return balance
}
