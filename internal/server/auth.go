package server

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/foto-leistenschneider/admin-panel/internal/config"
	"github.com/workos/workos-go/v4/pkg/auditlogs"
	"github.com/workos/workos-go/v4/pkg/directorysync"
	"github.com/workos/workos-go/v4/pkg/organizations"
	"github.com/workos/workos-go/v4/pkg/passwordless"
	"github.com/workos/workos-go/v4/pkg/portal"
	"github.com/workos/workos-go/v4/pkg/sso"
	"github.com/workos/workos-go/v4/pkg/usermanagement"
)

func init() {
	sso.Configure(config.WorkosApiKey, config.WorkosClientId)
	organizations.SetAPIKey(config.WorkosApiKey)
	passwordless.SetAPIKey(config.WorkosApiKey)
	directorysync.SetAPIKey(config.WorkosApiKey)
	auditlogs.SetAPIKey(config.WorkosApiKey)
	usermanagement.SetAPIKey(config.WorkosApiKey)
	portal.SetAPIKey(config.WorkosApiKey)
}

func getLoginRedirectURl(r *http.Request) string {
	scheme := r.URL.Scheme
	if len(scheme) == 0 {
		if r.TLS == nil {
			scheme = "http"
		} else {
			scheme = "https"
		}
	}
	return scheme + "://" + r.Host + "/api/login_callback"
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	codeVerifierBytes := make([]byte, 32)
	if _, err := rand.Read(codeVerifierBytes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	codeVerifier := hex.EncodeToString(codeVerifierBytes)

	sha256Hash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(sha256Hash[:])

	http.SetCookie(w, &http.Cookie{
		Name:     "code_verifier",
		Value:    codeVerifier,
		Path:     "/",
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})

	authUrl, err := usermanagement.GetAuthorizationURL(
		usermanagement.GetAuthorizationURLOpts{
			Provider:            "authkit",
			ClientID:            config.WorkosClientId,
			RedirectURI:         getLoginRedirectURl(r),
			CodeChallenge:       codeChallenge,
			CodeChallengeMethod: "S256",
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, authUrl.String(), http.StatusFound)
}

func loginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "code is empty", http.StatusBadRequest)
		return
	}

	codeVerifierCookie, err := r.Cookie("code_verifier")
	if err != nil {
		http.Error(w, "code_verifier cookie not found", http.StatusBadRequest)
		return
	}

	authenticateResponse, err := usermanagement.AuthenticateWithCode(
		r.Context(),
		usermanagement.AuthenticateWithCodeOpts{
			Code:         code,
			ClientID:     config.WorkosClientId,
			CodeVerifier: codeVerifierCookie.Value,
			IPAddress:    r.RemoteAddr,
			UserAgent:    r.UserAgent(),
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "code_verifier",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   r.TLS != nil,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    authenticateResponse.AccessToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    authenticateResponse.User.ID,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	http.SetCookie(w, &http.Cookie{
		Name:   "code_verifier",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "access_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "user_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func getAuthenticatedUser(
	r *http.Request,
) (usermanagement.User, bool) {
	userIDCookie, err := r.Cookie("user_id")
	if err != nil || userIDCookie.Value == "" {
		return usermanagement.User{}, false
	}

	user, err := usermanagement.GetUser(
		r.Context(),
		usermanagement.GetUserOpts{
			User: userIDCookie.Value,
		},
	)
	if err != nil {
		return usermanagement.User{}, false
	}

	return user, true
}
