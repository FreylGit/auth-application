package tests

import (
	"auth-application/tests/suite"
	auth "github.com/FreylGit/protoModel/gen/go/auth"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const passwordLength = 10

var appSecret = "secretTest"

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()
	name := gofakeit.Name()
	respR, err := st.AuthClient.Register(ctx, &auth.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     name,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respR.GetUserId())

	respL, err := st.AuthClient.Login(ctx, &auth.LoginRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)

	token := respL.GetAccessToken()
	assert.NotEmpty(t, token)
	tokenParse, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)
	claim, ok := tokenParse.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, respR.GetUserId(), int64(claim["uid"].(float64)))
	assert.Equal(t, email, int64(claim["email"].(float64)))
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passwordLength)
}
