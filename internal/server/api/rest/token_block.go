package resthandler

import (
	"net/http"
)

// TokenBlock - block the user token
func (s Handler) TokenBlock(w http.ResponseWriter, r *http.Request) {

	accessToken := r.FormValue("access_token") // access_token will be "" if parameter is not set
	if accessToken == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	accessToken, err := s.token.Block(accessToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.log.Info(accessToken)
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}
