package gapi

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jasonwebb3152/webb-ai/api/pb"
	mockdb "github.com/jasonwebb3152/webb-ai/db/mock"
	db "github.com/jasonwebb3152/webb-ai/db/sqlc"
	"github.com/jasonwebb3152/webb-ai/util"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func randomSession(t *testing.T) db.Session {
	return db.Session{
		ID:           uuid.New(),
		Username:     util.RandomString(8),
		RefreshToken: util.RandomString(8),
		UserAgent:    util.RandomString(8),
		ClientIp:     util.RandomString(8),
		IsBlocked:    false,
		ExpiresAt:    time.Now(),
		CreatedAt:    time.Now(),
	}
}

func TestLoginUser(t *testing.T) {
	user, password := randomUser(t)
	session := randomSession(t)

	testCases := []struct {
		name          string
		req           *pb.LoginUserRequest
		buildStubs    func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, rsp *pb.LoginUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.LoginUserRequest{
				Username: user.Username,
				Password: password,
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				mockStore.
					EXPECT().
					GetUser(gomock.Any(), user.Username).
					Times(1).
					Return(user, nil)
				mockStore.
					EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(session, nil)
			},
			checkResponse: func(t *testing.T, rsp *pb.LoginUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, rsp)

				require.Equal(t, session.ID.String(), rsp.SessionId)
				require.Equal(t, session.RefreshToken, rsp.RefreshToken)
				require.NotNil(t, rsp.AccessToken)
				require.Equal(t, session.ExpiresAt.UTC(), rsp.RefreshTokenExpiresAt.AsTime())
				require.WithinDuration(t, time.Now().Add(time.Minute*30), rsp.AccessTokenExpiresAt.AsTime(), time.Second*10)

				require.Equal(t, user.Username, rsp.User.Username)
				require.Equal(t, user.Email, rsp.User.Email)
				require.Equal(t, user.FullName, rsp.User.FullName)
			},
		},
		{
			name: "UserNotFound",
			req: &pb.LoginUserRequest{
				Username: user.Username,
				Password: password,
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				mockStore.
					EXPECT().
					GetUser(gomock.Any(), user.Username).
					Times(1).
					Return(db.User{}, db.ErrRecordNotFound)
				mockStore.
					EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, rsp *pb.LoginUserResponse, err error) {
				require.Error(t, err)
				require.Empty(t, rsp)
				status, _ := status.FromError(err)
				require.Equal(t, codes.InvalidArgument, status.Code())
			},
		},
		{
			name: "CreateSessionFailure",
			req: &pb.LoginUserRequest{
				Username: user.Username,
				Password: password,
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				mockStore.
					EXPECT().
					GetUser(gomock.Any(), user.Username).
					Times(1).
					Return(user, nil)
				mockStore.
					EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Session{}, errors.New("failed create"))
			},
			checkResponse: func(t *testing.T, rsp *pb.LoginUserResponse, err error) {
				require.Error(t, err)
				require.Empty(t, rsp)
				status, _ := status.FromError(err)
				require.Equal(t, codes.Internal, status.Code())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeController := gomock.NewController(t)
			defer storeController.Finish()

			mockStore := mockdb.NewMockStore(storeController)
			tc.buildStubs(mockStore)

			server := newTestServer(mockStore)

			rsp, err := server.LoginUser(context.Background(), tc.req)

			tc.checkResponse(t, rsp, err)
		})
	}
}
