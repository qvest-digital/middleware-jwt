package middleware-jwt

import (
	"context"
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

func JwtAuthAllowAll(jwtSecret string) func(http.Handler) http.Handler {
	middleware := jwtMiddleware{
		jwtSecret:     []byte(jwtSecret),
		allowedGroups: []string{},
	}
	return middleware.Handler
}

func JwtAuthAnyGroup(jwtSecret string, allowedGroups ...string) func(http.Handler) http.Handler {
	middleware := jwtMiddleware{
		jwtSecret:     []byte(jwtSecret),
		allowedGroups: allowedGroups,
	}
	return middleware.Handler
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

		log := logrus.
			WithField("remoteAddress", r.RemoteAddr).
			WithField("URI", r.RequestURI).
			WithField("method", r.Method)

		header := r.Header.Get("Authorization")
		if len(header) == 0 {
			log.Error("JWT Bearer token not present")
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
				Error("JWT parse error")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			log.Error("JWT is invalid")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Error("JWT has no claims")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		roles, ok := claims["groups"]
		if !ok {
			log.Error("Missing \"groups\" claim")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if !hasAnyGroup(roles, m.allowedGroups) {
			log.Error("JWT has none of the allowed groups")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Add claims to HTTP context
		ctx := context.WithValue(r.Context(), interface{}("claims"), claims)

		// Call next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
