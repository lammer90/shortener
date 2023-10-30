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
	userId := c.findAuth(req)
	if userId == "" {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	c.shortener.FindURLByUser(res, req, &handlers.RequestContext{
		UserId: userId,
	})
}

type RequestWithId struct {
	*http.Request
	UserId int
}

func (c Authenticator) checkAuth(res http.ResponseWriter, req *http.Request) *handlers.RequestContext {
	userId := c.findAuth(req)
	if userId == "" {
		userId = uuid.New().String()
	}

	c.repository.Save(userId)
	token, _ := buildJWTString(c.privateKey, userId)
	http.SetCookie(res, &http.Cookie{
		Name:  "Authorization",
		Value: token,
		Path:  "/",
	})

	return &handlers.RequestContext{
		UserId: userId,
	}
}

func (c Authenticator) findAuth(req *http.Request) string {
	userId := ""
	for _, cookie := range req.Cookies() {
		if cookie.Name == "Authorization" {
			userId = getUserId(cookie.Value, c.privateKey, c.repository)
		}
	}
	return userId
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func getUserId(tokenString string, privateKey string, repository userstorage.Repository) string {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(privateKey), nil
		})
	if err != nil {
		return ""
	}

	if !token.Valid {
		logger.Log.Info("Token is not valid")
		return ""
	}

	if _, ok, _ := repository.Find(claims.UserID); !ok {
		logger.Log.Info("Token didn't find")
		return ""
	}

	logger.Log.Info("Token is valid")
	return claims.UserID
}

func buildJWTString(privateKey string, userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{},
		UserID:           userId,
	})

	tokenString, err := token.SignedString([]byte(privateKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
