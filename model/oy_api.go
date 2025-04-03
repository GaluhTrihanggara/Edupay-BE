package model

type OyBillerApi struct {
	CustomerId  string  `json:"customer_id"`
	ProductId   string  `json:"product_id"`
	PartnerTxId string  `json:"partner_tx_id"`
	Amount      float64 `json:"amount"`
}

type OyBillerStatus struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type OyBillerApiResponse struct {
	OyBillerStatus `json:"status"`
	OyBillerData   `json:"data"`
}

type OyBillerData struct {
	TxID        string  `json:"tx_id"`
	CustomerID  string  `json:"customer_id"`
	ProductID   string  `json:"product_id"`
	PartnerTxID string  `json:"partner_tx_id"`
	Amount      float64 `json:"amount"`
	AdminFee    float64 `json:"admin_fee"`
	Description string  `json:"description"`
}

type GenerateVirtualAgregator struct {
	PatnerUserId string  `json:"partner_user_id"`
	BankCode     string  `json:"bank_code"`
	Amount       float64 `json:"amount"`
	Username     string  `json:"username_display"`
	IsOpen       bool    `json:"is_open"`
	SingleUse    bool    `json:"is_single_use"`
}

type PartnerCallbackVirtualAggregator struct {
	Amount            float64 `json:"amount"`
	VaNumber          string  `json:"va_number"`
	Success           bool    `json:"success"`
	PartnerUserId     string  `json:"partner_user_id"`
	TxDate            string  `json:"tx_date"`
	UsernameDisplay   string  `json:"username_display"`
	TrxExpirationDate string  `json:"trx_expiration_date"`
	TrxId             string  `json:"trx_id"`
	SettlementTime    string  `json:"settlement_time"`
	SettlementStatus  string  `json:"settlement_status"`
}
