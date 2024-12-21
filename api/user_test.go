package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/auronvila/simple-bank/db/mock"
	simplebank "github.com/auronvila/simple-bank/db/sqlc"
	"github.com/auronvila/simple-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Created",
			body: gin.H{
				"username":  util.RandomString(8),
				"full_name": util.RandomString(8),
				"email":     util.RandomEmail(),
				"password":  util.RandomString(8),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(simplebank.User{
						Username:          util.RandomString(8),
						FullName:          util.RandomString(8),
						Email:             util.RandomEmail(),
						PasswordChangedAt: time.Time{},
						CreatedAt:         time.Time{},
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "DuplicatedUsername",
			body: gin.H{
				"username":  "existingUser",
				"full_name": util.RandomString(8),
				"email":     util.RandomEmail(),
				"password":  util.RandomString(8),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(simplebank.User{}, &pq.Error{Code: "23505", Message: "username already exists"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

				var response map[string]interface{}
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Contains(t, response["error"], "username already exists")
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"username":  nil,
				"full_name": nil,
				"email":     nil,
				"password":  nil,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username":  util.RandomString(8),
				"full_name": util.RandomString(8),
				"email":     util.RandomEmail(),
				"password":  util.RandomString(8),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(simplebank.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(mockStore)

			server := NewServer(mockStore)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/create-user")

			body, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}

}
