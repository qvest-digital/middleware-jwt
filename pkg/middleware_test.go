package middleware-jwt

//go:generate mockgen -destination=mocks/http.go net/http Handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "github.com/tarent/middleware-jwt/internal/mocks"
)

type testCase struct {
	name           string
	auth           string
	status         int
	shouldCallNext bool
}

func TestHandlerAnyGroup(t *testing.T) {

	a := assert.New(t)

	// given: test subject
	subj := JwtAuthAnyGroup("mysecret", "groupB")

	testCases := []testCase{
		{
			name:           "happy path",
			auth:           "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZ3JvdXBzIjpbImdyb3VwQSIsImdyb3VwQiJdLCJpYXQiOjE1MTYyMzkwMjJ9.pPJGnFh4FUJnIcnReZlrrraG0Ep_bqEadYo6iH4KdHY",
			status:         http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:   "JWT is expired",
			auth:   "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZXhwIjoxNTE2MjM5MDIzLCJuYW1lIjoiSm9obiBEb2UiLCJncm91cHMiOlsiZ3JvdXBBIiwiZ3JvdXBCIl0sImlhdCI6MTUxNjIzOTAyMn0.pbsRS5wUW4Xkl6tcHU1H3hLkdTPWt9dWgESdy-PcFc4",
			status: http.StatusUnauthorized,
		},
		{
			name:   "Malformed JWT",
			auth:   "Bearer ",
			status: http.StatusUnauthorized,
		},
		{
			name:   "invalid signature",
			auth:   "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.W1ojVjJn5drM08g7Wm1QETFjc_VcqMYQxFP54KTAg-s",
			status: http.StatusUnauthorized,
		},
		{
			name:   "missing groups",
			auth:   "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.drt_po6bHhDOF_FJEHTrK-KD8OGjseJZpHwHIgsnoTM",
			status: http.StatusForbidden,
		},
		{
			name:   "malformed auth header",
			auth:   "Mal Formed",
			status: http.StatusUnauthorized,
		},
		{
			name:   "missing group",
			auth:   "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZ3JvdXBzIjpbImdyb3VwQSJdLCJpYXQiOjE1MTYyMzkwMjJ9.3np3-5cPj6iUHxgnzAc-bMA46jJIDmstiW5c5xKL7wg",
			status: http.StatusForbidden,
		},
	}

	// for each test case
	for _, v := range testCases {

		fmt.Println("Running test case: " + v.name)

		ctrl := gomock.NewController(t)

		// and: test request with JWT
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", v.auth)
		rr := httptest.NewRecorder()

		// and: http.Handler mock
		handlerMock := mocks.NewMockHandler(ctrl)
		if v.shouldCallNext {
			handlerMock.EXPECT().ServeHTTP(gomock.Any(), gomock.Any())
		}

		// when: handler is called
		subj(handlerMock).ServeHTTP(rr, req)

		// then: correct status code is returned
		a.Equal(v.status, rr.Code, "Wrong status code")

		ctrl.Finish()
	}

}

func TestHandlerAllowAll(t *testing.T) {

	a := assert.New(t)

	// given: test subject
	subj := JwtAuthAllowAll("mysecret")

	testCases := []testCase{
		{
			name:           "happy path",
			auth:           "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZ3JvdXBzIjpbImdyb3VwQSIsImdyb3VwQiJdLCJpYXQiOjE1MTYyMzkwMjJ9.pPJGnFh4FUJnIcnReZlrrraG0Ep_bqEadYo6iH4KdHY",
			status:         http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "no groups",
			auth:           "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZ3JvdXBzIjpbXSwiaWF0IjoxNTE2MjM5MDIyfQ.TLxPs_rZJbcTOfFD2XPYR2Lr6mkaJvRMTMi0usBd_B0",
			status:         http.StatusOK,
			shouldCallNext: true,
		},
	}

	// for each test case
	for _, v := range testCases {

		fmt.Println("Running test case: " + v.name)

		ctrl := gomock.NewController(t)

		// and: test request with JWT
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", v.auth)
		rr := httptest.NewRecorder()

		// and: http.Handler mock
		handlerMock := mocks.NewMockHandler(ctrl)
		if v.shouldCallNext {
			handlerMock.EXPECT().ServeHTTP(gomock.Any(), gomock.Any())
		}

		// when: handler is called
		subj(handlerMock).ServeHTTP(rr, req)

		// then: correct status code is returned
		a.Equal(v.status, rr.Code, "Wrong status code")

		ctrl.Finish()
	}

}

func TestHandlerMissingJWTCookie(t *testing.T) {
	a := assert.New(t)

	// given: test subject
	subj := JwtAuthAnyGroup("mysecret")

	// and: test request with empty JWT
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// when: handler is called
	subj(nil).ServeHTTP(rr, req)

	// then: 401
	a.Equal(http.StatusUnauthorized, rr.Code, "Wrong status code")
}
