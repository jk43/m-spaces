package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type JsonObject = map[string]any
type JsonArray = []any

type MapStringAny = map[string]any
type MapStringSlice[T any] map[string]T

type KeyValuePair struct {
	Key   string `bson:"key" json:"key"`
	Value any    `bson:"value" json:"value"`
}

type Hashed = string
type Salt = string

func HashPassword(pw string) (Hashed, Salt, error) {
	cost, _ := strconv.Atoi(os.Getenv("PW_COST"))
	pepper := os.Getenv("PW_PEPPER")
	salt := GetRandStrings(10)
	if pw == "" {
		return "", "", errors.New("The password does not meet our requirements.")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(pepper+pw+salt), cost)
	return string(bytes), salt, err
}

func CheckPasswordHash(pw, hashed Hashed, salt Salt) bool {
	pepper := os.Getenv("PW_PEPPER")
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pepper+pw+salt))
	return err == nil
}

type JWTClaims struct {
	ID             string   `json:"id"`
	Email          string   `json:"email"`
	FirstName      string   `json:"firstName"`
	LastName       string   `json:"lastName"`
	OrganizationID string   `json:"orgId"`
	ProfileImage   string   `json:"profileImage"`
	Role           UserRole `json:"role"`
	Metadata       any      `json:"metadata"`
	IP             string   `json:"ip"`
	jwt.RegisteredClaims
}

func GetJWTClaims(token string, host string) (*JWTClaims, error) {
	config := GetJWTConfig(host)
	claims := &JWTClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected singing method:%v", token.Header["alg"])
		}
		return []byte(config.Secret()), nil
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return nil, errors.New("expired token")
		}
		return nil, err
	}

	if claims.Issuer != config.Domain() {
		return nil, errors.New("incorrect issuer")
	}
	return claims, nil
}

var jwtConfigInstance *jwtConfig
var once sync.Once

type jwtConfig struct {
	secret             string
	refreshSecret      string
	cookieName         string
	domain             string
	tokenExpiry        time.Duration
	refreshTokenExpiry time.Duration
}

func (c *jwtConfig) Secret() string {
	return c.secret
}

func (c *jwtConfig) RefreshSecret() string {
	return c.refreshSecret
}

func (c *jwtConfig) Domain() string {
	return c.domain
}

func (c *jwtConfig) CookieName() string {
	return c.cookieName
}

func (c *jwtConfig) TokenExpiry() time.Duration {
	return c.tokenExpiry
}

func (c *jwtConfig) RefreshTokenExpiry() time.Duration {
	return c.refreshTokenExpiry
}

func GetJWTConfig(domain string) *jwtConfig {
	// once.Do(func() {
	secrets := os.Getenv("JWT_AUTH_SECRET")
	refreshSecrets := os.Getenv("JWT_REFRESH_SECRET")
	cookieName := os.Getenv("REFRESH_TOKEN_COOKIE_NAME")
	expiry, _ := strconv.Atoi(os.Getenv("JWT_TOKEN_EXPIRY"))
	refreshExpiry, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_EXPIRY"))

	jwtConfigInstance = &jwtConfig{
		secret:             secrets,
		refreshSecret:      refreshSecrets,
		cookieName:         cookieName,
		domain:             domain,
		tokenExpiry:        time.Minute * time.Duration(expiry),
		refreshTokenExpiry: time.Minute * time.Duration(refreshExpiry),
	}
	// })
	return jwtConfigInstance
}

func GetVerificationCode() string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	confirmCode := seededRand.Intn(9000) + 1000
	return strconv.Itoa(confirmCode)
}

const (
	GeneralGRPC    = "general"
	FileClientGRPC = "file-client"
)

func GetGRPCAddr(grpcType string) string {
	if grpcType == GeneralGRPC {
		return "0.0.0.0:5000"
	}
	if grpcType == FileClientGRPC {
		return "0.0.0.0:5001"
	}
	return "0.0.0.0:5000"
}

func TermDebugging(name string, v any) {
	_, callerFile, callerLine, _ := runtime.Caller(1)
	//funcName := runtime.FuncForPC(callerPc).Name()
	callerParts := strings.Split(callerFile, "/")
	caller := callerParts[len(callerParts)-1]

	fmt.Printf("\n\n\033[33m@@@@@ Debugging %s : %s @@@@@\033[0m \n\n\033[36m%+v\n\n", name, caller+":"+strconv.Itoa(callerLine), v)
}

// GetRandPassword generates a random password of length for user created by the admin
func GetRandStrings(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()_+`-=[]{}|;':,./<>?")

	rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetMongoFilters(src bson.M, search MapStringSlice[[]string]) bson.M {
	for k, v := range search {
		if len(v) == 1 {
			src[k] = bson.M{
				"$regex": primitive.Regex{
					Pattern: v[0],
					Options: "i", // case insensitive
				},
			}
		} else {
			src[k] = bson.M{"$in": k}
		}
	}
	return src
}

func GenerateSlug(f func(string, ...string) (string, error), slug string, params ...string) string {
	result, err := f(slug, params...)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return slug
	}
	parts := strings.Split(result, "-")
	lastPart := parts[len(parts)-1]

	if num, err := strconv.Atoi(lastPart); err == nil {
		parts[len(parts)-1] = strconv.Itoa(num + 1)
		slug = strings.Join(parts, "-")
	} else {
		slug = result + "-1"
	}
	return GenerateSlug(f, slug, params...)
}

type YesOrNo string

const (
	Yes YesOrNo = "Y"
	No  YesOrNo = "N"
)
