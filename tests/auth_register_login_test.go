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
		return []byte(st.Cfg.Secret), nil
	})
	require.NoError(t, err)
	claim, ok := tokenParse.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, respR.GetUserId(), int64(claim["uid"].(float64)))
	assert.Equal(t, email, claim["email"])
}

func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()
	name := gofakeit.Name()
	respReg, err := st.AuthClient.Register(ctx, &auth.RegisterRequest{
		Email:    email,
		Password: pass,
		Name:     name,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respReg.GetUserId())

	respReg, err = st.AuthClient.Register(ctx, &auth.RegisterRequest{
		Email:    email,
		Password: pass,
		Name:     name,
	})
	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestRegisterLogin_InvalidParam(t *testing.T) {
	ctx, st := suite.New(t)
	respReg, err := st.AuthClient.Register(ctx, &auth.RegisterRequest{
		Email:    "",
		Password: "",
		Name:     "",
	})
	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "invalid params")
	respReg, err = st.AuthClient.Register(ctx, &auth.RegisterRequest{
		Email:    "",
		Password: randomFakePassword(),
		Name:     gofakeit.Name(),
	})
	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "invalid params")
}

func TestRefreshToken_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)
	email := gofakeit.Email()
	pass := randomFakePassword()
	name := gofakeit.Name()
	respReg, err := st.AuthClient.Register(ctx, &auth.RegisterRequest{
		Email:    email,
		Password: pass,
		Name:     name,
	})
	require.NoError(t, err)
	userId := respReg.GetUserId()
	rToken := respReg.GetRefreshToken()

	assert.NotEmpty(t, userId)
	assert.NotEmpty(t, rToken)

	respRef, err := st.AuthClient.Refresh(ctx, &auth.RefreshRequest{
		UserId:       userId,
		RefreshToken: rToken,
	})
	require.NoError(t, err)

	rToken = respRef.GetRefreshToken()
	aToken := respRef.GetAccessToken()
	assert.NotEmpty(t, rToken)
	assert.NotEmpty(t, aToken)
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passwordLength)
}
