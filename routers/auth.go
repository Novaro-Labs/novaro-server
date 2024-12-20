package routers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/casbin/casbin/v2/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io"
	"net"
	"net/http"
	"net/url"
	"novaro-server/auth"
	"novaro-server/config"
	"novaro-server/model"
	"novaro-server/service"
	"strings"
	"time"
)

const (
	VerifierSize = 32
)

func AddAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")

	auth.GET("/login", login)
	auth.GET("/callback", callback)
	auth.GET("/wallet_login", WalletLogin)
	auth.POST("/connect_wallet", ConnectWallet)
}

func login(c *gin.Context) {
	queryParams := &url.Values{}
	query := c.Request.URL.Query()
	if query.Has("code") {
		code := query.Get("code")
		if len(code) != 6 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid invitation code, length should be 6"})
			return
		}
		if exist, err := service.NewInvitationCodesService().CheckInvitationCodes(code); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		} else if !exist {
			c.JSON(400, gin.H{"error": "invalid invitation code"})
			return
		}

		queryParams.Add("code", code)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "The code is required"})
		return
	}

	//redirectUri := redirectUrl(c, queryParams)
	// codeVerifier := generateCodeVerifier()

	// querys := url.Values{
	// 	"response_type": {"code"},
	// 	"client_id":     {config.Get().Client.ClientId},
	// 	//"redirect_uri":          {redirectUri},
	// 	"redirect_uri":          {config.Get().X.RedirectUri},
	// 	"scope":                 {config.Get().X.Scope},
	// 	"state":                 {codeVerifier},
	// 	"code_challenge":        {codeVerifier},
	// 	"code_challenge_method": {"plain"},
	// }

	//url := config.Get().X.AuthorizeUrl + querys.Encode()
	//c.Redirect(302, url)
	c.JSON(200, gin.H{"message": "the code verification passed"})
}

func callback(c *gin.Context) {
	query := c.Request.URL.Query()
	code := query.Get("code")

	queryParams := &url.Values{}
	if query.Has("code") {
		invitation_code := query.Get("code")
		queryParams.Add("code", invitation_code)
	}
	redirectUri := redirectUrl(c, queryParams)

	codeVerifier, _ := url.QueryUnescape(query.Get("state"))
	token, err := getToken(code, redirectUri, codeVerifier)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	user, err := getUserInfo(token)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	invitationsService := service.NewInvitationsService()
	if isCheckInvitee, _ := invitationsService.CheckInvitee(user.Id); !isCheckInvitee {
		invitationsService.Save(&model.Invitations{InviterId: user.Id, InviteeId: user.Id, Rewards: 100})
	} else if isCheckInvitationReward, _ := invitationsService.CheckInvitationRewards(user.Id, code); !isCheckInvitationReward {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid login attempt"})
		return
	}

	err = service.NewTwiitterUserService().SaveTwitterUsers(user.ToTwitterUsers())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	id, err := service.NewUserService().SaveUsers(user.ToUsers())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(c)
	session.Set(auth.Userkey, id)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(200, user)
}

func redirectUrl(c *gin.Context, query *url.Values) string {
	redirectUrl := &url.URL{
		Scheme:   "http",
		Host:     c.Request.Host,
		Path:     "/v1/auth/callback",
		RawQuery: query.Encode(),
	}
	if c.Request.TLS != nil {
		redirectUrl.Scheme = "https"
	}

	return redirectUrl.String()
}

func getToken(code, redirectUri, codeVerifier string) (string, error) {
	cli := config.Get().Client
	form := url.Values{
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {redirectUri},
		"code_verifier": {codeVerifier},
	}
	body := form.Encode()
	fmt.Printf("token request body: %v\n", body)
	request, err := http.NewRequest("POST", config.Get().X.Oauth2TokenUrl, strings.NewReader(body))
	if err != nil {
		log.LogError(err, "new token request error")
		return "", err
	}

	request.SetBasicAuth(cli.ClientId, cli.ClientSecret)
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

func getUserInfo(token string) (*model.TwitterUserInfo, error) {
	request, err := http.NewRequest("GET", config.Get().X.UserProfile, nil)
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

	proxy := config.Get().Client.Proxy

	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		// We use ABSURDLY large keys, and should probably not.
		TLSHandshakeTimeout: 60 * time.Second,
	}
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
	}
	client := &http.Client{
		Transport: transport,
	}
	return client
}

type ConnectWalletRequest struct {
	WalletAddr string `json:"wallet_addr" binding:"required"`
	Id         string `json:"id" binding:"required"`
}

func ConnectWallet(c *gin.Context) {
	var connectWallet ConnectWalletRequest
	err := c.ShouldBindJSON(&connectWallet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userService := service.NewUserService()
	if isExist, _ := userService.UserExists(connectWallet.Id); isExist {
		c.JSON(http.StatusOK, gin.H{"error": "wallet already connected"})
		return
	}
	_, err = userService.SaveUsers(&model.Users{Id: connectWallet.Id, WalletPublicKey: connectWallet.WalletAddr})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "wallet connected"})
	return

}

type WalletLoginRequest struct {
	WalletAddr string `json:"wallet_addr" binding:"required"`
}

func WalletLogin(c *gin.Context) {
	var walletLogin WalletLoginRequest
	err := c.ShouldBindJSON(&walletLogin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userService := service.NewUserService()
	if isExist, _ := userService.UserExists(walletLogin.WalletAddr); !isExist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid " + walletLogin.WalletAddr})
		return
	}

	session := sessions.Default(c)
	session.Set(auth.Userkey, walletLogin.WalletAddr)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

}

type Token struct {
	TokenType    string `json:"token_type"`
	Expires      int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

type UserData struct {
	Data *model.TwitterUserInfo `json:"data"`
}
