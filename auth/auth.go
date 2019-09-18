package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig *oauth2.Config

func init() {
	gotenv.Load()
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_OAUTH_CALL_BACK"),
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

type GoogleUserStruct struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	Hd            string `json:"hd"`
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

//OAuthGooleLogin redireciona para o tela de google
func OAuthGooleLogin(c *gin.Context) {
	oauthState := generateStateOauthCookie(c.Writer)
	u := googleOauthConfig.AuthCodeURL(oauthState)
	c.JSON(200, gin.H{"url": u})
	return
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

//OauthGoogleCallback callback
func OauthGoogleCallback(c *gin.Context) {
	// // Read oauthState from Cookie
	// // GetOrCreate User in your db.
	// // Redirect or response with a token.
	// // More code .....
	// // fmt.Fprintf(c.Writer, "UserInfo: %s\n", data)
	// user := &models.User{}
	token, err := googleOauthConfig.Exchange(context.Background(), c.Request.FormValue("code"))
	if err != nil {
		c.JSON(400, gin.H{"errors": []string{err.Error()}})
		c.Abort()
		return
	}
	idGoogle, err := getdadosToken(token.AccessToken)
	if err != nil {
		c.JSON(400, gin.H{"errors": []string{err.Error()}})
		c.Abort()
		return
	}

	url := c.Request.Host + c.Request.URL.String()
	b := make([]byte, 16)

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":         url,
		"iat":         time.Now().Unix(),
		"nbf":         time.Now().Unix(),
		"exp":         time.Now().Add(168 * time.Hour).Unix(),
		"jti":         fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]),
		"id":          idGoogle.Id,
		"email":       idGoogle.Email,
		"given_name":  idGoogle.GivenName,
		"picture":     idGoogle.Picture,
		"locale":      idGoogle.Locale,
		"family_name": idGoogle.FamilyName,
		"hd":          idGoogle.Hd,
		"tokenGoogle": token.AccessToken,
	})

	tokenString, err := tok.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(400, gin.H{"errors": []string{"Error on database connection"}})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"token": tokenString})
	return
}

func getdadosToken(token string) (*GoogleUserStruct, error) {

	response, err := http.Get(oauthGoogleUrlAPI + token)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var ret = &GoogleUserStruct{}
	if err := json.NewDecoder(response.Body).Decode(&ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func getIdInInt(user GoogleUserStruct) (int, error) {
	return strconv.Atoi(user.Id)
}

func GetDadosGoogle(c *gin.Context) {
	token := c.Params.ByName("id")

	response, err := resolveToken(token)
	if err != nil {
		c.JSON(400, gin.H{"authErrors": []string{"Invalid payload"}})
		c.Abort()
	}
	c.JSON(200, gin.H{"data": response})
	return
}

func resolveToken(token string) (*GoogleUserStruct, error) {
	tok, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := tok.Claims.(jwt.MapClaims); ok && tok.Valid {
		id := claims["id"].(string)
		email := claims["email"].(string)
		given_name := claims["given_name"].(string)
		Picture := claims["picture"].(string)
		Locale := claims["locale"].(string)
		family_name := claims["family_name"].(string)
		Hd := claims["hd"].(string)
		ret := &GoogleUserStruct{
			Email:      email,
			Id:         id,
			FamilyName: family_name,
			GivenName:  given_name,
			Hd:         Hd,
			Picture:    Picture,
			Locale:     Locale,
		}
		return ret, nil
	}
	return nil, errors.New("token invalido")
}
