package resthandler

import (
	"net/http"
)

// UserBlock - block the user and all his tokens
func (s Handler) UserBlock(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username") // username will be "" if parameter is not set
	if username == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID, err := s.user.Block(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = s.token.BlockAllTokenUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.log.Info(userID)
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}
