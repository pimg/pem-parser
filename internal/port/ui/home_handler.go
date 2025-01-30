package ui

import (
	"net/http"
)

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	err := s.Templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
