package middleware

import (
	"errors"
	"family_budget/internal/entities/users"
	"family_budget/pkg/database"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"strings"
	"time"
)

// GinJWTMiddleware provides a Json-Web-Token authentication implementation. On failure, a 401 HTTP response
// is returned. On success, the wrapped middleware is called, and the userID is made available as
// c.Get("userID").(string).
// Users can get a token by posting a json request to LoginHandler. The token then needs to be passed in
// the Authentication header. Example: Authorization:Bearer XXX_TOKEN_XXX
type GinJWTMiddleware struct {
	// Realm name to display to the user. Required.
	Realm string

	// signing algorithm - possible values are HS256, HS384, HS512
	// Optional, default is HS256.
	SigningAlgorithm string

	// Secret access token key used for signing. Required.
	AccessKey []byte

	// Secret refresh token key used for signing. Required.
	RefreshKey []byte

	// Duration that a jwt token is valid. Optional, defaults to one hour.
	AccessTimeout time.Duration

	// Duration that a refresh jwt token is valid.
	RefreshTimeout time.Duration

	// This field allows clients to refresh their token until MaxRefresh has passed.
	// Note that clients can refresh their token in the last moment of MaxRefresh.
	// This means that the maximum validity timespan for a token is MaxRefresh + Timeout.
	// Optional, defaults to 0 meaning not refreshable.
	MaxRefresh time.Duration

	// Callback function that should perform the authentication of the user based on userID and
	// password. Must return true on success, false on failure. Required.
	// Option return user id, if so, user id will be stored in Claim Array.
	Authenticator func(login string, password string, c *gin.Context, ipadress, otp string) (string, bool)

	// Callback function that should perform the authorization of the authenticated user. Called
	// only after an authentication success. Must return true on success, false on failure.
	// Optional, default to success.
	Authorizator func(userID string, c *gin.Context) (string, bool)

	// Callback function that will be called during login.
	// Using this function it is possible to add additional payload data to the webtoken.
	// The data is then made available during requests via c.Get("JWT_PAYLOAD").
	// Note that the payload is not encrypted.
	// The attributes mentioned on jwt.io can't be used as keys for the map.
	// Optional, by default no additional data will be set.
	PayloadFunc func(userID string) map[string]interface{}

	// User can define own Unauthorized func.
	Unauthorized func(*gin.Context, int, string)

	// Set the identity handler function
	IdentityHandler func(jwt.MapClaims) string

	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	TokenLookup string

	// TokenHeadName is a string in the header. Default value is "Bearer"
	TokenHeadName string

	// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
	TimeFunc func() time.Time

	DB *gorm.DB
}

// Login form structure.
type Login struct {
	Login      string `form:"login" json:"login" binding:"required"`
	Password   string `form:"password" json:"password" binding:"required"`
	DeviceUUID string `form:"device_uuid" json:"device_uuid"`
	Captcha    string `form:"captcha" json:"captcha,omitempty"`
	Otp        string `form:"otp" json:"otp,omitempty"`
}

type ResponseStruct struct {
	UserId        string    `json:"user_id"`
	AccessToken   string    `json:"access_token"`
	RefreshToken  string    `json:"refresh_token"`
	AccessExpire  time.Time `json:"access_expire"`
	RefreshExpire time.Time `json:"refresh_expire"`
}

type attempt struct {
	captcha  string
	failNum  int
	lastTime time.Time
}

// MiddlewareInit initialize jwt configs.
func (mw *GinJWTMiddleware) MiddlewareInit() error {

	if mw.TokenLookup == "" {
		mw.TokenLookup = "header:Authorization"
	}

	if mw.SigningAlgorithm == "" {
		mw.SigningAlgorithm = "HS256"
	}

	if mw.AccessTimeout == 0 {
		mw.AccessTimeout = time.Hour
	}

	if mw.RefreshTimeout == 0 {
		mw.RefreshTimeout = 24 * time.Hour
	}

	if mw.TimeFunc == nil {
		mw.TimeFunc = time.Now
	}

	mw.TokenHeadName = strings.TrimSpace(mw.TokenHeadName)
	if len(mw.TokenHeadName) == 0 {
		mw.TokenHeadName = "Bearer"
	}

	if mw.Authorizator == nil {
		mw.Authorizator = func(userID string, c *gin.Context) (string, bool) {
			return "", true
		}
	}

	if mw.Unauthorized == nil {
		mw.Unauthorized = func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		}
	}

	if mw.IdentityHandler == nil {
		mw.IdentityHandler = func(claims jwt.MapClaims) string {
			return claims["id"].(string)
		}
	}

	if mw.AccessKey == nil {
		return errors.New("secret key is required")
	}

	return nil
}

// MiddlewareFunc makes GinJWTMiddleware implement the Middleware interface.
func (mw *GinJWTMiddleware) MiddlewareFunc() gin.HandlerFunc {
	if err := mw.MiddlewareInit(); err != nil {
		return func(c *gin.Context) {
			mw.unauthorized(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	return func(c *gin.Context) {
		mw.middlewareImpl(c)
		return
	}
}

func (mw *GinJWTMiddleware) middlewareImpl(c *gin.Context) {
	token, err := mw.parseToken(c)

	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, err.Error())
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	id := mw.IdentityHandler(claims)
	c.Set("JWT_PAYLOAD", claims)
	c.Set("user_id", id)
	if msg, ok := mw.Authorizator(id, c); !ok {
		mw.unauthorized(c, http.StatusForbidden, msg)
		return
	}

	//pass username of active user to map
	//online.Online.UpdateCreate(id, time.Now())

	c.Next()
}

// --Captcha--
func (mw *GinJWTMiddleware) sendCaptcha(c *gin.Context, userCache *attempt) {
	capt := captcha.New()

	capt.SetSize(128, 54)
	capt.SetDisturbance(captcha.MEDIUM)
	capt.SetFrontColor(color.RGBA{A: 191})
	capt.SetBkgColor(color.RGBA{R: 232, G: 232, B: 232, A: 255})

	img, str := capt.Create(6, captcha.NUM)
	userCache.captcha = str
	c.Writer.WriteHeader(429)
	c.Writer.Header().Set("Content-Type", "image/png")
	err := png.Encode(c.Writer, img)

	if err != nil {
		log.Println("Write image error:", err.Error())
		return
	}

	return
}

// LoginHandler - Авторизация пользователя
// @Summary Авторизация пользователя
// @ID authorize-user
// @Tags Регистрация и Авторизация
// @Produce json
// @Param id body Login true "Даные для авторизации"
// @Success 200 {object} ResponseStruct
// @Failure 400 {string} string "reason"
// @Failure 401 {string} string "reason"
// @Failure 402 {string} string "reason"
// @Router /visor/login [post]
func (mw *GinJWTMiddleware) LoginHandler(c *gin.Context) {
	_ = mw.MiddlewareInit()
	var loginVals Login

	if err := c.BindJSON(&loginVals); err != nil {
		mw.unauthorized(c, http.StatusBadRequest, "Missing Login or Password")
		return
	}

	log.Println("LoginHandler loginVals:", loginVals)

	if mw.Authenticator == nil {
		mw.unauthorized(c, http.StatusInternalServerError, "Missing define authenticator func")
		return
	}
	ipAddress := ""
	fwdAddress := c.GetHeader("X-Forwarded-For")
	if fwdAddress != "" {

		ipAddress = fwdAddress

		ips := strings.Split(fwdAddress, ", ")
		if len(ips) > 1 {
			ipAddress = ips[0]
		}
	}
	log.Println("clientIP:", ipAddress)
	log.Println("Header:", c.Request.Header)

	userID, ok := mw.Authenticator(loginVals.Login, loginVals.Password, c, ipAddress, loginVals.Otp)

	if !ok {
		mw.unauthorized(c, http.StatusUnauthorized, "Неверный логин или пароль")
		return
	}

	// Create the token
	atToken := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	atClaims := atToken.Claims.(jwt.MapClaims)

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(loginVals.Login) {
			atClaims[key] = value
		}
	}

	if userID == "" {
		userID = loginVals.Login
	}

	atExpire := mw.TimeFunc().Add(mw.AccessTimeout)
	atClaims["id"] = userID
	atClaims["t_exp"] = atExpire.Unix()
	atClaims["orig_iat"] = mw.TimeFunc().Unix()

	accessToken, err := atToken.SignedString(mw.AccessKey)

	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, "Create JWT Token failed")
		return
	}

	var curUser users.User
	if err := mw.DB.Raw("select * from users where login = ?", loginVals.Login).Scan(&curUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	rtToken := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	rtClaims := rtToken.Claims.(jwt.MapClaims)

	rtExpire := mw.TimeFunc().Add(mw.RefreshTimeout)

	rtClaims["t_exp"] = rtExpire.Unix()
	rtClaims["orig_iat"] = mw.TimeFunc().Unix()
	rtClaims["id"] = userID
	rtClaims["user_id"] = curUser.ID
	rtClaims["family_id"] = curUser.FamilyID
	rtClaims["role_id"] = curUser.RoleID
	refreshToken, err := rtToken.SignedString(mw.RefreshKey)

	if err != nil {
		pretty.Logln("error: can't create jwt token ")
		mw.unauthorized(c, http.StatusUnauthorized, "Create JWT Token failed")
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":        userID,
		"access_token":   accessToken,
		"refresh_token":  refreshToken,
		"access_expire":  atExpire.Format(time.RFC3339),
		"refresh_expire": rtExpire.Format(time.RFC3339),
	})
}

// RefreshToken refreshes token when necessary
func (mw *GinJWTMiddleware) RefreshToken(c *gin.Context) {
	rtTokenReq, err := mw.jwtFromCookie(c, "Refresh-Authorization")

	token, err := jwt.Parse(rtTokenReq, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(mw.SigningAlgorithm) != token.Method {
			return nil, errors.New("invalid signing algorithm")
		}

		return mw.RefreshKey, nil
	})

	if err != nil {
		pretty.Logln("error: signing algorithm")
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	//we dont retain refresh token
	expire := int64(claims["t_exp"].(float64))
	userID := int64(claims["user_id"].(int))

	if expire < mw.TimeFunc().Unix() {
		pretty.Logln("error: expired token")
		fmt.Println("token is expired1")
		mw.unauthorized(c, http.StatusUnauthorized, "token is expired")
		return
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		pretty.Logln("error: invalid token")
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	// Create the token
	atToken := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	atClaims := atToken.Claims.(jwt.MapClaims)

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(claims["id"].(string)) {
			atClaims[key] = value
		}
	}

	atExpire := mw.TimeFunc().Add(mw.AccessTimeout)
	atClaims["id"] = claims["id"].(string)
	atClaims["t_exp"] = atExpire.Unix()
	atClaims["orig_iat"] = mw.TimeFunc().Unix()

	accessToken, err := atToken.SignedString(mw.AccessKey)

	if err != nil {
		pretty.Logln("error: can't create jwt token ")
		mw.unauthorized(c, http.StatusUnauthorized, "Create JWT Token failed")
		return
	}

	rtToken := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	rtClaims := rtToken.Claims.(jwt.MapClaims)

	rtExpire := mw.TimeFunc().Add(mw.RefreshTimeout)

	rtClaims["t_exp"] = rtExpire.Unix()
	rtClaims["orig_iat"] = mw.TimeFunc().Unix()
	rtClaims["id"] = claims["id"].(int)
	rtClaims["family_id"] = claims["family_id"].(int)
	rtClaims["user_id"] = userID
	rtClaims["role_id"] = claims["role_id"].(int)

	refreshToken, err := rtToken.SignedString(mw.RefreshKey)

	if err != nil {
		pretty.Logln("error: can't create jwt token ")
		mw.unauthorized(c, http.StatusUnauthorized, "Create JWT Token failed")
	}

	c.JSON(http.StatusOK, gin.H{
		"userId":         userID,
		"access_token":   accessToken,
		"refresh_token":  refreshToken,
		"access_expire":  atExpire.Format(time.RFC3339),
		"refresh_expire": rtExpire.Format(time.RFC3339),
	})
}

// RefreshHandler can be used to refresh a token. The token still needs to be valid on refresh.
// Shall be put under an endpoint that is using the GinJWTMiddleware.
// Reply will be of the form {"token": "TOKEN"}.
func (mw *GinJWTMiddleware) RefreshHandler(c *gin.Context) {
	token, _ := mw.parseToken(c)
	claims := token.Claims.(jwt.MapClaims)

	origIat := int64(claims["orig_iat"].(float64))

	if origIat < mw.TimeFunc().Add(-mw.MaxRefresh).Unix() {
		mw.unauthorized(c, http.StatusUnauthorized, "Token is expired")
		return
	}

	// Create the token
	newToken := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	newClaims := newToken.Claims.(jwt.MapClaims)

	for key := range claims {
		newClaims[key] = claims[key]
	}

	expire := mw.TimeFunc().Add(mw.AccessTimeout)
	newClaims["id"] = claims["id"]
	newClaims["t_exp"] = expire.Unix()
	newClaims["orig_iat"] = origIat

	tokenString, err := newToken.SignedString(mw.AccessKey)

	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, "Create JWT Token failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  tokenString,
		"expire": expire.Format(time.RFC3339),
	})
}

// TokenGenerator handler that clients can use to get a jwt token.
func (mw *GinJWTMiddleware) TokenGenerator(userID string) string {
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(userID) {
			claims[key] = value
		}
	}

	claims["id"] = userID
	claims["t_exp"] = mw.TimeFunc().Add(mw.AccessTimeout).Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()

	tokenString, _ := token.SignedString(mw.AccessKey)

	return tokenString
}

func (mw *GinJWTMiddleware) jwtFromHeader(c *gin.Context, key string) (string, error) {
	authHeader := c.Request.Header.Get(key)

	if authHeader == "" {
		return "", errors.New("auth header empty")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == mw.TokenHeadName) {
		return "", errors.New("invalid auth header")
	}
	return parts[1], nil
}

func (mw *GinJWTMiddleware) jwtFromQuery(c *gin.Context, key string) (string, error) {
	token := c.Query(key)
	if token == "" {
		return "", errors.New("query token empty")
	}
	return token, nil
}

func (mw *GinJWTMiddleware) jwtFromCookie(c *gin.Context, key string) (string, error) {
	cookie, _ := c.Cookie(key)
	if cookie == "" {
		return "", errors.New("cookie token empty")
	}
	return cookie, nil
}

func (mw *GinJWTMiddleware) parseToken(c *gin.Context) (*jwt.Token, error) {
	var tokenStr string
	var err error
	parts := strings.Split(mw.TokenLookup, ":")
	switch parts[0] {
	case "header":
		tokenStr, err = mw.jwtFromHeader(c, parts[1])
	case "query":
		tokenStr, err = mw.jwtFromQuery(c, parts[1])
	case "cookie":
		tokenStr, err = mw.jwtFromCookie(c, parts[1])
	}
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(mw.SigningAlgorithm) != token.Method {
			return nil, errors.New("invalid signing algorithm")
		}
		return mw.AccessKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	expire := int64(claims["t_exp"].(float64))

	if expire < mw.TimeFunc().Unix() {
		return nil, errors.New("token is expired")
	}

	return token, err
}

func (mw *GinJWTMiddleware) unauthorized(c *gin.Context, code int, message string) {

	if mw.Realm == "" {
		mw.Realm = "gin jwt"
	}

	c.Header("WWW-Authenticate", "JWT realm="+mw.Realm)
	c.Abort()

	mw.Unauthorized(c, code, message)

	return
}

func Authenticator(login string, password string, c *gin.Context, ipadress string, otp string) (string, bool) {
	var user users.User

	if err := database.Postgres().Where("login = ? and active = true", login).Find(&user).Limit(1).Error; err != nil {
		log.Println("Authenticator func query error:", err.Error())
		return "", false
	}

	if checkPassword(user.Password, password) {
		return login, true
	}

	return "", false
}

func checkPassword(existing, provided string) bool {
	log.Println("pass - existing:", existing, " provided:", provided)
	err := bcrypt.CompareHashAndPassword([]byte(existing), []byte(provided))
	if err != nil {
		log.Println("checkPass err", err.Error())
		return false
	}
	return true
}

func Payload(login string) map[string]interface{} {
	var user users.User

	if err := database.Postgres().Where("login = ? ", login).Find(&user); err.Error != nil {
		return map[string]interface{}{
			"user_id":   0,
			"role":      "",
			"user_name": "",
		}
	}

	return map[string]interface{}{
		"user_id":   user.ID,
		"family_id": user.FamilyID,
		"role_id":   user.RoleID,
	}
}
