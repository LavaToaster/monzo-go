package monzo

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
)

type CounterParty struct {
	Name              string `json:"name"`
	UserId            string `json:"user_id"`
	PreferredName     string `json:"preferred_name"`      // Present on Monzo P2P
	ServiceUserNumber string `json:"service_user_number"` // Present on BACS (DD's)
	AccountNumber     string `json:"account_number"`      // Present on Faster Payments and BACS (DD's)
	SortCode          string `json:"sort_code"`           // ^
}

func (cp *CounterParty) IsMonzoUser() bool {
	if cp.UserId == "" {
		return false
	}

	return strings.Contains(cp.UserId, "anon")
}

type Attachment struct {
	Id         string `json:"id"`
	Type       string `json:"type"`
	Url        string `json:"url"`
	ExternalId string `json:"external_id"`
	FileType   string `json:"file_type"`
	FileUrl    string `json:"file_url"`
	UserId     string `json:"user_id"`
	Created    string `json:"created"`
}

type MerchantAddress struct {
	ShortFormatted string  `json:"short_formatted"`
	Formatted      string  `json:"formatted"`
	Address        string  `json:"address"`
	City           string  `json:"city"`
	Country        string  `json:"country"`
	Postcode       string  `json:"postcode"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	ZoomLevel      int64   `json:"zoom_level"`
	Approximate    bool    `json:"approximate"`
}

type Merchant struct {
	Id              string            `json:"id"`
	GroupId         string            `json:"group_id"`
	Created         string            `json:"created"`
	Name            string            `json:"name"`
	Logo            string            `json:"logo"`
	Emoji           string            `json:"emoji"`
	Category        string            `json:"category"`
	Online          bool              `json:"online"`
	Atm             bool              `json:"atm"`
	Address         MerchantAddress   `json:"address"`
	Updated         string            `json:"updated"`
	Metadata        map[string]string `json:"metadata"`
	DisableFeedback bool              `json:"disable_feedback"`
}

// Unable to add because I don't have any reference:
// - Fees
type Transaction struct {
	Id                         string            `json:"id"`
	Amount                     int64             `json:"amount"`
	Description                string            `json:"description"`
	Created                    string            `json:"created"`
	Currency                   string            `json:"currency"`
	Merchant                   Merchant          `json:"merchant"`
	Notes                      string            `json:"notes"`
	Metadata                   map[string]string `json:"metadata"`
	Labels                     []string          `json:"labels"` // "withdrawal.atm.international" is the only one i've seen
	AccountBalance             int64             `json:"account_balance"`
	Attachments                []Attachment      `json:"attachments"`
	Category                   string            `json:"category"`
	IsLoad                     bool              `json:"is_load"`
	Settled                    string            `json:"settled"`
	LocalAmount                int64             `json:"local_amount"`
	LocalCurrency              string            `json:"local_currency"`
	Updated                    string            `json:"updated"`
	AccountId                  string            `json:"account_id"`
	UserId                     string            `json:"user_id"`
	CounterParty               CounterParty      `json:"counter_party"`
	Scheme                     string            `json:"scheme"`
	DedupeId                   string            `json:"dedupe_id"`
	Originator                 bool              `json:"originator"`
	IncludedInSpending         bool              `json:"included_in_spending"`
	CanBeExcludedFromBreakdown bool              `json:"can_be_excluded_from_breakdown"`
}

type TransactionsList struct {
	Transactions []Transaction `json:"transactions"`
}

type TransactionContainer struct {
	Transaction Transaction `json:"transaction"`
}

func (c *Client) Transactions(accountId string) []Transaction {
	resp, _ := c.Do("GET", "/transactions?expand[]=merchant&account_id="+accountId, nil)
	defer resp.Body.Close()

	var transactions TransactionsList

	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		log.Println(err)
	}

	return transactions.Transactions
}

func (c *Client) Transaction(transactionId string) Transaction {
	resp, _ := c.Do("GET", "/transactions/"+transactionId+"?expand[]=merchant", nil)
	defer resp.Body.Close()

	var transaction TransactionContainer

	if err := json.NewDecoder(resp.Body).Decode(&transaction); err != nil {
		log.Println(err)
	}

	return transaction.Transaction
}

func (c *Client) AnnotateTransaction(transactionId string, metadata map[string]string) Transaction {
	form := url.Values{}

	for key, value := range metadata {
		form.Set("metadata["+key+"]", value)
	}

	req, _ := c.NewRequest("PATCH", "/transactions/"+transactionId, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Print(err)
	}

	var transaction TransactionContainer

	if err := json.NewDecoder(resp.Body).Decode(&transaction); err != nil {
		log.Println(err)
	}

	return transaction.Transaction
}
