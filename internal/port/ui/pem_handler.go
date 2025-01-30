package ui

import (
	"fmt"
	"net/http"
	"strings"
)

type PEMParserPage struct {
	SuccessMessage string
	ErrorMessage   string
	Result         string
}

func (s *Server) pemParserHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("received pem parser request")

	page := &PEMParserPage{}

	pem, err := parseForm(r)
	if err != nil {
		page.ErrorMessage = "failed to parse PEM"
	}

	out, err := s.App.PEMHandler.Handle(pem)
	if err != nil {
		page.ErrorMessage = err.Error()
	}

	if out != nil {
		page.SuccessMessage = fmt.Sprintf("Successfully parsed PEM %s file", out.Type)
		page.Result = out.Result
	}

	err = s.Templates.ExecuteTemplate(w, "result-block", page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func parseForm(r *http.Request) ([]byte, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	return []byte(strings.TrimSpace(r.PostFormValue("pem"))), nil
}
