package request

type WalletDepositRequest struct {
	WalletId string  `json:"wallet_id" validate:"required" binding:"required"`
	Amount   float32 `json:"amount" validate:"required" binding:"required"`
}
