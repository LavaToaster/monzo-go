package monzo

import (
	"encoding/json"
	"log"
)

type WhoAmI struct {
	Authenticated bool `json:"authenticated"`
	ClientId string `json:"client_id"`
	UserId string `json:"user_id"`
}

func (c *Client) WhoAmI() WhoAmI {
	resp, _ := c.Do("GET", "/ping/whoami", nil)
	defer resp.Body.Close()

	var whoAmI WhoAmI

	if err := json.NewDecoder(resp.Body).Decode(&whoAmI); err != nil {
		log.Println(err)
	}

	return whoAmI
}
