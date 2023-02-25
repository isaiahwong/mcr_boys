package internal

type Account struct {
	AccountId int     `json:"account_id"`
	Balance   float64 `json:"balance"`
	Owner     string  `json:"owner"`
}
