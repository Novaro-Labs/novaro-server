package router

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"novaro-server/config"
	"novaro-server/model"
	"strings"
	"time"

	"github.com/casbin/casbin/v2/log"
	"github.com/gin-gonic/gin"
)

const (
	VerifierSize = 32
)

func AddAuthRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/auth")

	group.GET("/login", login)
	group.GET("/callback", callback)
}

func login(c *gin.Context) {
	redirectUri := redirectUrl(c)
	codeVerifier := generateCodeVerifier()

	fmt.Printf("codeVerifier: %v\n", codeVerifier)

	querys := url.Values{
		"response_type":         {"code"},
		"client_id":             {config.ClientId},
		"redirect_uri":          {redirectUri},
		"scope":                 {"tweet.read+users.read+follows.read+offline.access"},
		"state":                 {codeVerifier},
		"code_challenge":        {codeVerifier},
		"code_challenge_method": {"plain"},
	}

	url := "https://x.com/i/oauth2/authorize?" + querys.Encode()
	c.Redirect(302, url)
}

func callback(c *gin.Context) {
	redirectUri := redirectUrl(c)
	query := c.Request.URL.Query()
	code := query.Get("code")
	codeVerifier, _ := url.QueryUnescape(query.Get("state"))

	token, err := getToken(code, redirectUri, codeVerifier)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	user, err := getUserInfo(token)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	err = model.SaveUsers(user.toUsers())
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, user)
}

func redirectUrl(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	redirectUri := fmt.Sprintf("%s://%s/v1/auth/callback", scheme, c.Request.Host)
	return redirectUri
}

func getToken(code, redirectUri, codeVerifier string) (string, error) {
	form := url.Values{
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {redirectUri},
		"code_verifier": {codeVerifier},
	}
	body := form.Encode()
	fmt.Printf("token request body: %v\n", body)
	request, err := http.NewRequest("POST", "https://api.x.com/2/oauth2/token", strings.NewReader(body))
	if err != nil {
		log.LogError(err, "new token request error")
		return "", err
	}

	request.SetBasicAuth(config.ClientId, config.ClientSecret)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := newHttpClient()
	response, err := client.Do(request)
	if err != nil {
		log.LogError(err, "get token error")
		fmt.Printf("get token error: %v\n", err)
		return "", err
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		rbody, _ := io.ReadAll(response.Body)
		return "", fmt.Errorf("request token failed: %s %v", response.Status, string(rbody))
	}

	var token Token
	err = json.NewDecoder(response.Body).Decode(&token)
	if err != nil {
		log.LogError(err, "decode body error")
		return "", err
	}
	return token.AccessToken, nil
}

func getUserInfo(token string) (*UserInfo, error) {
	request, err := http.NewRequest("GET", "https://api.x.com/2/users/me?user.fields=created_at,profile_image_url", nil)
	if err != nil {
		log.LogError(err, "new userinfo request error")
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+token)

	client := newHttpClient()
	response, err := client.Do(request)
	if err != nil {
		log.LogError(err, "get user info error")
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("get user info failed: %s", response.Status)
	}
	defer response.Body.Close()

	var user UserData
	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		log.LogError(err, "decode body error")
		return nil, err
	}
	return user.Data, nil
}

func generateCodeVerifier() string {
	b := make([]byte, VerifierSize)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func newHttpClient() *http.Client {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		// We use ABSURDLY large keys, and should probably not.
		TLSHandshakeTimeout: 60 * time.Second,
	}
	if config.Proxy != "" {
		proxyUrl, err := url.Parse(config.Proxy)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
	}
	client := &http.Client{
		Transport: transport,
	}
	return client
}

type Token struct {
	TokenType    string `json:"token_type"`
	Expires      int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

type UserData struct {
	Data *UserInfo `json:"data"`
}

type UserInfo struct {
	Avatar   string    `json:"profile_image_url"`
	Name     string    `json:"name"`
	Created  time.Time `json:"created_at"`
	Id       string    `json:"id"`
	Username string    `json:"username"`
}

func (userInfo *UserInfo) toUsers() *model.Users {
	users := &model.Users{
		TwitterId: userInfo.Id,
		UserName:  userInfo.Username,
		CreatedAt: userInfo.Created,
		Avatar:    &userInfo.Avatar,
	}
	return users
}
