package response

type GetWalletDetailsResponse struct {
	Balance        float32 `json:"balance"`
	AboveThreshold bool    `json:"above_threshold"`
}
