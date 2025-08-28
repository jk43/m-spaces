package service

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/moly-space/molylibs/service"
	"github.com/moly-space/molylibs/utils"
)

type TokenPairs struct {
	Token        string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func GenerateTokenPair(user *service.User, orgSettings *service.Settings, r *http.Request) (TokenPairs, utils.MapStringAny, error) {
	host, err := utils.GetXHost(r)
	if err != nil {
		return TokenPairs{}, nil, err
	}
	config := utils.GetJWTConfig(host)
	claims := utils.JWTClaims{}

	claims.Subject = user.ID.Hex()
	claims.Audience = []string{config.Domain()}
	claims.Issuer = config.Domain()
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(config.TokenExpiry()))
	//custom values
	claims.ID = user.ID.Hex()
	claims.Email = user.Email
	claims.FirstName = user.FirstName
	claims.LastName = user.LastName
	claims.ProfileImage = user.ProfileImage
	claims.Role = user.Role
	claims.Metadata = user.Metadata
	claims.OrganizationID = user.OrganizationID.Hex()
	claims.IP = r.RemoteAddr
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	singedAccessToken, err := token.SignedString([]byte(config.Secret()))
	if err != nil {
		return TokenPairs{}, nil, err
	}
	//create the refresh token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)

	refreshTokenClaims["sub"] = user.ID
	//set expiry ; must be longer than jwt expiry
	refreshTokenClaims["exp"] = time.Now().Add(config.RefreshTokenExpiry()).Unix()

	singedRefreshToken, err := refreshToken.SignedString([]byte(config.RefreshSecret()))
	if err != nil {
		return TokenPairs{}, nil, err
	}

	var tokenPairs = TokenPairs{
		Token:        singedAccessToken,
		RefreshToken: singedRefreshToken,
	}
	sharableClaims := utils.MapStringAny{
		"id":             user.ID.Hex(),
		"firstName":      user.FirstName,
		"lastName":       user.LastName,
		"email":          user.Email,
		"profileImage":   user.ProfileImage,
		"registerMethod": user.RegisterMethod,
		"metadata":       user.Metadata,
	}

	return tokenPairs, sharableClaims, nil
}
