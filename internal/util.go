package jwt

import (
	"net/http"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

func GetGroupsFromAuthenticatedRequest(r *http.Request) []string {

	ctx := r.Context()
	if ctx == nil {
		logrus.Error("Could not get groups from request: no context available")
		return []string{}
	}

	claims, ok := ctx.Value("claims").(jwtgo.MapClaims)
	if !ok {
		logrus.Error("Could not get claims from context: cast failed")
		return []string{}
	}

	return castGroupsToString(claims["groups"])
}

func castGroupsToString(groups interface{}) []string {

	stringGroups := []string{}

	casted, ok := groups.([]interface{})
	if !ok {
		logrus.WithField("groups", groups).Error("Error casting groups to array")
		return []string{}
	}

	for _, v := range casted {
		s, ok := v.(string)
		if ok {
			stringGroups = append(stringGroups, string(s))
		}
	}

	return stringGroups
}

func hasAnyGroup(groups interface{}, required []string) bool {

	if len(required) == 0 {
		return true
	}

	stringGroups := castGroupsToString(groups)

	return contains(stringGroups, required)
}

func contains(have []string, required []string) bool {
	for _, v1 := range have {
		for _, v2 := range required {
			if v1 == v2 {
				return true
			}
		}
	}
	return false
}
