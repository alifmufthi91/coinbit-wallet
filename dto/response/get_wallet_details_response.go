package response

type GetWalletDetailsResponse struct {
	WalletId       string  `json:"wallet_id"`
	Balance        float32 `json:"balance"`
	AboveThreshold bool    `json:"above_threshold"`
}
