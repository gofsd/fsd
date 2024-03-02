package account

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"

	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofsd/fsd/pkg/cmd"
	"github.com/gofsd/fsd/pkg/meta"

	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
)

var ac = Account{}

func Get() Account {
	return ac
}

func Add(cmd cmd.Command) (ac byte, e error) {
	sha_512 := sha512.New()
	sha_512.Write([]byte(cmd.Args[0] + cmd.Args[1]))
	return
}

var accounts = make([]Account, 0)
var sessions = map[string]string{}

type Account struct {
	ID       int    `json:"id,omitempty" validate:"gt=0,lt=65000"`
	Email    string `json:"email" validate:"required,max=64,min=6,email"`
	Password string `json:"password" validate:"required,length=128"`
}

// AuthStruct -
type AuthStruct struct {
	Email    string `json:"username" validate:"required,email,max=50"`
	Password string `json:"password" validate:"required,min=8,max=25"`
}

// AuthResponse -
type AuthResponse struct {
	UserID int    `json:"user_id" validate:"required,email,max=50"`
	Email  string `json:"username" validate:"required,email,max=50"`
	Token  string `json:"token" validate:"required,email,max=50"`
}

func SignUps(ctx *gin.Context) {
	var authReq AuthStruct
	var authResp AuthResponse
	var account Account

	_, exist := lo.Find[Account](accounts, func(item Account) bool {
		if item.Email == authReq.Email {
			return true
		} else {
			return false
		}
	})

	if !exist {
		account.Email = authReq.Email
		account.Password = ToSh512(authReq.Password)
		account.ID = len(accounts) + 1
		accounts = append(accounts, account)
		token := GenerateSecureToken(100)
		sessions[token] = account.Email
		authResp.Token = token
		authResp.Email = account.Email
		authResp.UserID = account.ID
		ctx.JSON(http.StatusOK, authResp)
	} else {
		HandleError(ctx, AccountExist)
	}

}

func ToSh512(s string) string {
	sha_512 := sha512.New()
	sha_512.Write([]byte(s))
	var hashedPasswordHex = hex.EncodeToString(sha_512.Sum(nil))
	return hashedPasswordHex
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func HandleError(ctx *gin.Context, err error) {
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		panic(err.Error())
	}
}

func SignUp(email, password string) (auth []byte, e error) {
	sha512 := ToSh512(fmt.Sprintf("%s_%s", email, password))
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, e := token.SignedString([]byte(sha512))
	auth, e = meta.New(
		meta.
			Op().
			NDV("Email", email, "required,email").
			ReturnKeyAs("user_id").
			ReturnAs("email").
			NDV("Password", password, "required,min=8").
			NDV("SHA", sha512, "required,sha512").
			NDV("JWT", tokenString, "required,jwt").
			ReturnAs("jwt").
			Root("account", "Email", "SHA", "JWT").
			Crown("SHA").
			Crown("JWT").
			Create("Email", "SHA", "JWT").
			Set(),
	).ToJson()

	return auth, e
}
