package gapi

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jasonwebb3152/webb-ai/api/pb"
	mockdb "github.com/jasonwebb3152/webb-ai/db/mock"
	db "github.com/jasonwebb3152/webb-ai/db/sqlc"
	"github.com/jasonwebb3152/webb-ai/util"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type eqCreateUserTxParamsMatcher struct {
	arg      db.CreateUserParams
	password string
	user     db.User
}

func (expected eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	// In case, some value is nil
	actualArg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(expected.password, actualArg.HashedPassword)
	if err != nil {
		return false
	}

	expected.arg.HashedPassword = actualArg.HashedPassword

	if !reflect.DeepEqual(expected.arg, actualArg) {
		return false
	}

	return err == nil
}

func (e eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserTxParams(arg db.CreateUserParams, password string, user db.User) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, password, user}
}

func randomUser(t *testing.T) (db.User, string) {
	password := util.RandomPassword(10)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user := db.User{
		Username:          util.RandomString(8),
		FullName:          util.RandomString(8),
		Email:             util.RandomEmail(),
		HashedPassword:    hashedPassword,
		CreatedAt:         time.Now(),
		PasswordChangedAt: time.Now(),
	}
	return user, password
}

func TestCreateUser(t *testing.T) {
	user, password := randomUser(t)
	testCases := []struct {
		name          string
		req           *pb.CreateUserRequest
		buildStubs    func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, rsp *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				mockStore.
					EXPECT().
					CreateUser(gomock.Any(), EqCreateUserTxParams(arg, password, user)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, rsp *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, rsp)
				createdUser := rsp.GetUser()
				require.Equal(t, user.Username, createdUser.Username)
				require.Equal(t, user.Email, createdUser.Email)
				require.Equal(t, user.FullName, createdUser.FullName)
			},
		},
		{
			name: "BadRequestFields",
			req: &pb.CreateUserRequest{
				Username: "jw",
				FullName: "he",
				Email:    "hello",
				Password: "something",
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				mockStore.
					EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, rsp *pb.CreateUserResponse, err error) {
				require.Nil(t, rsp)
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())

				fields := []string{"username", "fullname", "password", "email"}
				descs := []string{"not between", "not between", "does not meet", "missing '@'"}
				details := st.Details()
				for idx, dt := range details {
					violation := dt.(*errdetails.BadRequest).
						FieldViolations[0]
					require.Equal(t, fields[idx], violation.Field)
					require.Contains(t, violation.Description, descs[idx])
				}
			},
		},
		{
			name: "UniqueViolation",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				mockStore.
					EXPECT().
					CreateUser(gomock.Any(), EqCreateUserTxParams(arg, password, user)).
					Times(1).
					Return(db.User{}, &pgconn.PgError{Code: db.UniqueViolation})
			},
			checkResponse: func(t *testing.T, rsp *pb.CreateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, rsp)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.AlreadyExists, st.Code())
			},
		},
	}

	for _, tc := range testCases {

		// This "t" is different than global "t"
		// So that the test cases run independently
		t.Run(tc.name, func(t *testing.T) {
			storeController := gomock.NewController(t)
			defer storeController.Finish()

			mockStore := mockdb.NewMockStore(storeController)
			tc.buildStubs(mockStore)

			server := newTestServer(mockStore)

			rsp, err := server.CreateUser(context.Background(), tc.req)

			tc.checkResponse(t, rsp, err)
		})
	}
}
