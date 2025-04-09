package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"

	loginService "projects/auth-api-gateway/service/login"
	regService "projects/auth-api-gateway/service/register"
	"projects/auth-api-gateway/service/session"
)

var router = gin.Default()

func init() {
	router.POST("/login", login)
	router.POST("/register", register)
	router.GET("/session", startSession)
}

func Start(binding string) {
	router.Run(binding)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegistrationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// we will use the login-request-id to find the session creation event.
// the session creation even will tell us wether a prevvious login request was successful or not.
// if the login was successful, we will create a session cookie for the user and a jwt
// if the login is not done yet, we will return 404
// if the login was not successful, we will return 403
func startSession(c *gin.Context) {

	cookie, err := c.Request.Cookie("LOGIN-REQUEST-ID")
	if err != nil {
		println("Error getting cookie:", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Missing cookie",
			"message": "LOGIN-REQUEST-ID cookie is required",
		})
		return
	}
	requestId := cookie.Value

	status, sessionValue := session.FindRegisteredSession(requestId)
	session.ClearRegisteredSession(requestId)

	if status == session.WAITING {
		c.JSON(http.StatusNotFound, gin.H{})
	} else if status == session.VALID {
		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie("REQUEST-TOKEN", sessionValue.SessionKey, 3600, "", "", true, true)
		c.JSON(http.StatusOK, gin.H{})
	} else if status == session.INVALID {
		c.JSON(http.StatusUnauthorized, gin.H{})
	}
}

func login(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		println("Error binding JSON:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	if req.Username == "" || req.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Error",
			"message": "username and password are required",
		})
		return
	}

	// Log credentials to console
	println("Username:", req.Username)
	println("Password:", req.Password)

	subscriberId, requestId := getMetaHeaders(c)
	if subscriberId == "" {
		return
	}

	loginRequestId, err := uuid.GenerateUUID()
	if err != nil {
		println("Error generating UUID:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	loginService.Login(req.Username, req.Password, loginRequestId, subscriberId, requestId)

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("LOGIN-REQUEST-ID", loginRequestId, 3600, "", "", true, true)
	c.JSON(http.StatusAccepted, gin.H{"requestId": requestId})
}

func getMetaHeaders(c *gin.Context) (string, string) {
	subscriberID := c.GetHeader("x-subscriber-id")
	requestID := c.GetHeader("x-request-id")

	if subscriberID == "" || requestID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Missing Headers",
			"message": "x-request-id and x-subscriber-id is required",
		})
		return "", ""
	}
	return subscriberID, requestID
}

func register(c *gin.Context) {
	var req RegistrationRequest
	if err := c.BindJSON(&req); err != nil {
		println("Error binding JSON:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	if req.Username == "" || req.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Error",
			"message": "username and password are required",
		})
		return
	}

	subscriberId, requestId := getMetaHeaders(c)
	if subscriberId == "" {
		return
	}

	regService.Register(req.Username, req.Password, subscriberId, requestId)

	c.JSON(http.StatusAccepted, gin.H{"requestId": requestId})
}
