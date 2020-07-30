package jwt

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

type jwtMiddleware struct {
	jwtSecret     []byte
	allowedGroups []string
}

func NewJwtMiddleware(jwtSecret string, allowedGroups []string) *jwtMiddleware {
	return &jwtMiddleware{
		jwtSecret:     []byte(jwtSecret),
		allowedGroups: allowedGroups,
	}
}

func (m *jwtMiddleware) parseFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		logrus.WithField("alg", token.Header["alg"]).Error("JWT error: unexpected signing method")
		return nil, errors.New("Unexpected signing method")
	}
	return m.jwtSecret, nil
}

func (m *jwtMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log := logrus.WithField("remoteAddress", r.RemoteAddr)

		header := r.Header.Get("Authorization")
		if len(header) == 0 {
			log.Error("JWT Bearer token not found")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) < 2 || parts[0] != "Bearer" {
			log.Error("Malformed JWT Bearer token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(parts[1], m.parseFunc)
		if err != nil {
			log.WithField("err", err).
				Error("JWT error")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			log.Error("JWT has no claims")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		roles, ok := claims["groups"]
		if !ok {
			log.Error("Missing groups claim")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if !m.hasAnyGroup(roles, m.allowedGroups) {
			log.Error("JWT missing required group")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *jwtMiddleware) hasAnyGroup(groups interface{}, required []string) bool {
	stringGroups := []string{}

	casted, ok := groups.([]interface{})
	if !ok {
		logrus.WithField("groups", groups).
			Error("Error casting groups to array")
		return false
	}

	for _, v := range casted {
		s, ok := v.(string)
		if ok {
			stringGroups = append(stringGroups, string(s))
		}
	}

	return m.contains(stringGroups, required)
}

func (m *jwtMiddleware) contains(have []string, required []string) bool {
	for _, v1 := range have {
		for _, v2 := range required {
			if v1 == v2 {
				return true
			}
		}
	}
	return false
}
