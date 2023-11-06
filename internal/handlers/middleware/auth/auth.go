package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/lammer90/shortener/internal/handlers"
	"github.com/lammer90/shortener/internal/logger"
	"github.com/lammer90/shortener/internal/userstorage"
	"net/http"
)

type Authenticator struct {
	shortener  handlers.ShortenerRestProviderWithContext
	repository userstorage.Repository
	privateKey string
}

func New(repository userstorage.Repository, privateKey string, shortener handlers.ShortenerRestProviderWithContext) handlers.ShortenerRestProvider {
	return Authenticator{
		shortener,
		repository,
		privateKey,
	}
}

func (c Authenticator) SaveShortURL(res http.ResponseWriter, req *http.Request) {
	requestContext := c.checkAuth(res, req)
	c.shortener.SaveShortURL(res, req, requestContext)
}

func (c Authenticator) FindByShortURL(res http.ResponseWriter, req *http.Request) {
	requestContext := c.checkAuth(res, req)
	c.shortener.FindByShortURL(res, req, requestContext)
}

func (c Authenticator) SaveShortURLApi(res http.ResponseWriter, req *http.Request) {
	requestContext := c.checkAuth(res, req)
	c.shortener.SaveShortURLApi(res, req, requestContext)
}

func (c Authenticator) SaveShortURLBatch(res http.ResponseWriter, req *http.Request) {
	requestContext := c.checkAuth(res, req)
	c.shortener.SaveShortURLBatch(res, req, requestContext)
}

func (c Authenticator) FindURLByUser(res http.ResponseWriter, req *http.Request) {
	userID := c.findAuth(req)
	if userID == "" {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	c.shortener.FindURLByUser(res, req, &handlers.RequestContext{
		UserID: userID,
	})
}

func (c Authenticator) Delete(res http.ResponseWriter, req *http.Request) {
	requestContext := c.checkAuth(res, req)
	c.shortener.Delete(res, req, requestContext)
}

func (c Authenticator) checkAuth(res http.ResponseWriter, req *http.Request) *handlers.RequestContext {
	userID := c.findAuth(req)
	if userID == "" {
		userID = uuid.New().String()
	}

	c.repository.Save(userID)
	token, _ := buildJWTString(c.privateKey, userID)
	http.SetCookie(res, &http.Cookie{
		Name:  "Authorization",
		Value: token,
		Path:  "/",
	})

	return &handlers.RequestContext{
		UserID: userID,
	}
}

func (c Authenticator) findAuth(req *http.Request) string {
	userID := ""
	for _, cookie := range req.Cookies() {
		if cookie.Name == "Authorization" {
			logger.Log.Info("Found Authorization token")
			userID = getUserID(cookie.Value, c.privateKey)
		}
	}
	return userID
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func getUserID(tokenString string, privateKey string) string {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(privateKey), nil
		})
	if err != nil {
		logger.Log.Info(err.Error())
		return ""
	}

	if !token.Valid {
		logger.Log.Info("Token is not valid")
		return ""
	}

	logger.Log.Info("Token is valid")
	return claims.UserID
}

func buildJWTString(privateKey string, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{},
		UserID:           userID,
	})

	tokenString, err := token.SignedString([]byte(privateKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
