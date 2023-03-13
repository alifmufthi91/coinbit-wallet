package controller_test

import (
	"bytes"
	"coinbit-wallet/controller"
	"coinbit-wallet/dto/request"
	"coinbit-wallet/dto/response"
	"coinbit-wallet/middleware"
	"coinbit-wallet/service"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type WalletControllerSuite struct {
	suite.Suite
	walletService    *service.MockWalletService
	walletController controller.WalletController
	router           *gin.Engine
}

func TestWalletControllerSuite(t *testing.T) {
	suite.Run(t, new(WalletControllerSuite))
}

func (c *WalletControllerSuite) SetupSuite() {
	c.walletService = new(service.MockWalletService)
	c.walletController = controller.NewWalletController(c.walletService)

	c.router = gin.Default()
	c.router.Use(middleware.ErrorHandlingMiddleware())
	c.router.POST("/deposit", c.walletController.Deposit)
	c.router.GET("/details/:walletId", c.walletController.GetDetails)
}

func (s *WalletControllerSuite) AfterTest(_, _ string) {
	s.walletService.AssertExpectations(s.T())
}

func (c *WalletControllerSuite) TestWalletController_Deposit() {

	// create new context
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	depositRequest := request.WalletDepositRequest{
		WalletId: "111-222",
		Amount:   2000,
	}
	c.walletService.On("DepositWallet", depositRequest).Return(nil).Once()

	// assign payload to context request body
	req, _ := http.NewRequest(http.MethodPost, "/deposit", nil)
	bodyBytes, _ := json.Marshal(depositRequest)
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req

	// call the endpoint
	c.router.ServeHTTP(w, req)

	// check response status code
	require.Equal(c.T(), http.StatusOK, w.Code)

	var responseBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &responseBody)

	require.Empty(c.T(), responseBody["data"])

	w2 := httptest.NewRecorder()
	ctx2, _ := gin.CreateTestContext(w)
	expectedError := errors.New("failed processing deposit to wallet")
	c.walletService.On("DepositWallet", depositRequest).Return(expectedError).Once()

	req2, _ := http.NewRequest(http.MethodPost, "/deposit", nil)
	req2.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	req2.Header.Set("Content-Type", "application/json")
	ctx2.Request = req2

	c.router.ServeHTTP(w2, req2)
	require.Equal(c.T(), http.StatusInternalServerError, w2.Code)
}

func (c *WalletControllerSuite) TestWalletController_GetDetails() {

	// create new context
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	walletId := "111-222"
	walletDetails := response.GetWalletDetailsResponse{
		WalletId:       walletId,
		Balance:        2000,
		AboveThreshold: false,
	}
	c.walletService.On("GetWalletDetails", walletId).Return(&walletDetails, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/details/"+walletId, nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req

	// call the endpoint
	c.router.ServeHTTP(w, req)

	// check response status code
	require.Equal(c.T(), http.StatusOK, w.Code)

	var responseBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &responseBody)

	data := responseBody["data"].(map[string]interface{})
	require.NotEmpty(c.T(), data)

	require.Equal(c.T(), walletDetails.WalletId, data["wallet_id"])
	require.Equal(c.T(), walletDetails.Balance, float32(data["balance"].(float64)))
	require.Equal(c.T(), walletDetails.AboveThreshold, data["above_threshold"])
}
