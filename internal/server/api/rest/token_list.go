package resthandler

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type ViewDataToken struct {
	Tokens map[string]string
}

// TokenList -  the page that displays all tokens user
func (s Handler) TokenList(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	tmpl, err := template.ParseFiles(s.config.TemplatePathToken)
	if err != nil {
		s.log.Errorf(fmt.Errorf("parse failed: %w", err).Error())
		http.Error(w, "Error loading user list page", http.StatusInternalServerError)
		return
	}

	tokens := make(map[string]string)

	id, err := s.user.GetUserID(username)
	if err != nil {
		s.log.Error(fmt.Errorf("error in get UserId from DB for user %s: %w", username, err))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	tokenList, err := s.token.GetList(id)
	if err != nil {
		s.log.Error(fmt.Errorf("error in get token list from DB for user %s: %w", username, err))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	for _, token := range tokenList {
		now := time.Now()
		end := token.EndDateAt
		check := now.After(end)

		if check {
			tokens[token.AccessToken] = "Block"
		} else {
			tokens[token.AccessToken] = "Active"
		}
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	data := ViewDataToken{Tokens: tokens}
	err = tmpl.Execute(w, data)
	if err != nil {
		s.log.Errorf(fmt.Errorf("execution failed: %w", err).Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
