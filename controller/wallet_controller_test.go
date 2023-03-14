package controller_test

import (
	"bytes"
	"coinbit-wallet/controller"
	"coinbit-wallet/dto/request"
	"coinbit-wallet/dto/response"
	"coinbit-wallet/middleware"
	"coinbit-wallet/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type WalletControllerSuite struct {
	suite.Suite
	walletService    *mocks.MockWalletService
	walletController controller.WalletController
	router           *gin.Engine
}

func TestWalletControllerSuite(t *testing.T) {
	suite.Run(t, new(WalletControllerSuite))
}

func (c *WalletControllerSuite) SetupSuite() {
	c.walletService = new(mocks.MockWalletService)
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

	testCases := []struct {
		name           string
		depositRequest request.WalletDepositRequest
		mockReturn     error
		expectedCode   int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Successful deposit",
			depositRequest: request.WalletDepositRequest{
				WalletId: "111-222",
				Amount:   2000,
			},
			mockReturn:   nil,
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": nil,
			},
		},
		{
			name: "Failed deposit",
			depositRequest: request.WalletDepositRequest{
				WalletId: "111-222",
				Amount:   2000,
			},
			mockReturn:   errors.New("failed processing deposit to wallet"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error":   "INTERNAL SERVER ERROR",
				"message": "failed processing deposit to wallet",
			},
		},
	}

	for _, tc := range testCases {
		c.Run(tc.name, func() {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			c.walletService.On("DepositWallet", tc.depositRequest).Return(tc.mockReturn).Once()

			bodyBytes, _ := json.Marshal(tc.depositRequest)
			req, _ := http.NewRequest(http.MethodPost, "/deposit", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = req

			c.router.ServeHTTP(w, req)

			require.Equal(c.T(), tc.expectedCode, w.Code)

			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)

			if tc.mockReturn != nil {
				require.Equal(c.T(), tc.expectedBody["message"], responseBody["message"])
			} else {
				require.Equal(c.T(), tc.expectedBody["data"], responseBody["data"])
			}
		})
	}
}

func (c *WalletControllerSuite) TestWalletController_GetDetails() {

	testCases := []struct {
		name            string
		walletId        string
		walletDetails   *response.GetWalletDetailsResponse
		returnError     error
		expectedCode    int
		expectedData    map[string]interface{}
		expectedMessage string
	}{
		{
			name:     "Get Wallet Details Success",
			walletId: "111-222",
			walletDetails: &response.GetWalletDetailsResponse{
				WalletId:       "111-222",
				Balance:        2000,
				AboveThreshold: false,
			},
			returnError:  nil,
			expectedCode: 200,
			expectedData: map[string]interface{}{
				"wallet_id":       "111-222",
				"balance":         float64(2000),
				"above_threshold": false,
			},
			expectedMessage: "SUCCESS",
		},
		{
			name:            "Get Wallet Details Failed",
			walletId:        "111-222",
			walletDetails:   nil,
			returnError:     errors.New("failed to get wallet details"),
			expectedCode:    500,
			expectedData:    nil,
			expectedMessage: "failed to get wallet details",
		},
	}

	for _, tc := range testCases {
		c.Run(tc.name, func() {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			c.walletService.On("GetWalletDetails", tc.walletId).Return(tc.walletDetails, tc.returnError).Once()

			req, _ := http.NewRequest(http.MethodGet, "/details/"+tc.walletId, nil)
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = req

			// call the endpoint
			c.router.ServeHTTP(w, req)

			// check response status code
			require.Equal(c.T(), tc.expectedCode, w.Code)

			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)

			require.Equal(c.T(), tc.expectedMessage, responseBody["message"].(string))

			if tc.returnError == nil {
				data := responseBody["data"].(map[string]interface{})
				require.Equal(c.T(), tc.expectedData["wallet_id"], data["wallet_id"])
				require.Equal(c.T(), tc.expectedData["balance"], data["balance"])
				require.Equal(c.T(), tc.expectedData["above_threshold"], data["above_threshold"])
			}
		})
	}
}
