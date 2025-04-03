package model

type Balance struct {
	UUIDPrimaryKey
	UserId          string  `gorm:"index" json:"user_id"`
	RecipentBank    string  `json:"recipient_bank"`
	RecipentAccount string  `json:"recipient_account"`
	Amount          float64 `json:"amount"`
	PartnerTxId     string  `json:"partner_trx_id"`
	Email           string  `json:"email"`
}
