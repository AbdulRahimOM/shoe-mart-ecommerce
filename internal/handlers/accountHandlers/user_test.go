package accounthandler

import (
	e "MyShoo/internal/domain/customErrors"
	request "MyShoo/internal/models/requestModels"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mockusecase "MyShoo/internal/mock/mockUseCase"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserSignup(t *testing.T) {
	testCase := map[string]struct {
		input         request.UserSignUpReq
		buildstub     func(useCaseMock *mockusecase.MockIUserUC, signupData request.UserSignUpReq)
		checkResponse func(t *testing.T, responserecorder *httptest.ResponseRecorder)
	}{
		"success": {
			input: request.UserSignUpReq{
				FirstName: "vajid",
				LastName:  "Ottapalam",
				Email:     "vajid44@gmail.com",
				Phone:     "+919876543210",
				Password:  "jhy78ij",
			},
			buildstub: func(useCaseMock *mockusecase.MockIUserUC, signupData request.UserSignUpReq) {
				useCaseMock.EXPECT().SignUp(&signupData).Times(1).Return(func() *string {
					text := "sampleToken"
					return &text
				}(),
					nil,
				)
			},
			checkResponse: func(t *testing.T, responserecorder *httptest.ResponseRecorder) {

				assert.Equal(t, http.StatusOK, responserecorder.Code)

				var responseBody map[string]interface{}
				if err := json.Unmarshal(responserecorder.Body.Bytes(), &responseBody); err != nil {
					t.Errorf("failed to decode response body: %v", err)
					return
				}

				value, ok := responseBody["token"].(string)
				if !ok || value != "sampleToken" {
					t.Errorf("expected JSON field 'token' with value 'sampleToken' not found")
				}
			},
		},
		"bad requst - email already in use": {
			input: request.UserSignUpReq{
				FirstName: "vajid",
				LastName:  "Ottapalam",
				Email:     "vajid44@gmail.com",
				Phone:     "+919876543210",
				Password:  "jhy78ij",
			},
			buildstub: func(useCaseMock *mockusecase.MockIUserUC, signupData request.UserSignUpReq) {
				useCaseMock.EXPECT().SignUp(&signupData).Times(1).Return(
					nil, e.ErrEmailAlreadyUsed_401,
				)
			},
			checkResponse: func(t *testing.T, responserecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, responserecorder.Code)

				var responseBody map[string]interface{}
				if err := json.Unmarshal(responserecorder.Body.Bytes(), &responseBody); err != nil {
					t.Errorf("failed to decode response body: %v", err)
					return
				}

				value, ok := responseBody["msg"].(string)
				if !ok || value != "email already used" {
					t.Errorf("expected JSON field 'msg' with value 'email already used' not found")
				}

				value, ok = responseBody["error"].(string)
				if !ok || value != "invalid req" {
					t.Errorf("expected JSON field 'error' with value 'invalid req' not found")
				}

			},
		},
	}

	for testname, test := range testCase {
		test := test
		t.Run(testname, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockUseCase := mockusecase.NewMockIUserUC(ctrl)
			test.buildstub(mockUseCase, test.input)
			userHandler := NewUserHandler(mockUseCase)

			server := gin.Default()
			server.POST("/signup", userHandler.PostSignUp)

			jsonData, err := json.Marshal(test.input)
			assert.NoError(t, err)
			body := bytes.NewBuffer(jsonData)

			mockRequst, err := http.NewRequest(http.MethodPost, "/signup", body)
			assert.NoError(t, err)
			responseRecord := httptest.NewRecorder()
			server.ServeHTTP(responseRecord, mockRequst)

			test.checkResponse(t, responseRecord)
		})
	}
}
