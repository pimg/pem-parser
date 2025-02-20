package ui

import (
	"fmt"
	"net/http"
	"pem-parser/internal/app"
	"strings"
	"time"
)

type PEMParserPage struct {
	SuccessMessage string
	ErrorMessage   string
	Result         *PEMResponse
}

type PEMResponse struct {
	SerialNumber            string
	DN                      DistinguishedName
	IssuerName              string
	NotBefore               time.Time
	NotAfter                time.Time
	SubjectAlternativeNames []string
	KeyUsages               []string
	Raw                     string
	Type                    string
}

type DistinguishedName struct {
	CommonName         string
	SerialNumber       string
	Country            string
	State              string
	Locality           string
	Organization       string
	OrganizationalUnit string
	EmailAddress       string
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
		page.Result = mapPEMResponse(out)
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

func mapPEMResponse(response *app.PEMResponse) *PEMResponse {
	return &PEMResponse{
		SerialNumber: response.SerialNumber,
		DN: DistinguishedName{
			CommonName:         response.DistinguishedName.CommonName,
			SerialNumber:       response.DistinguishedName.SerialNumber,
			Country:            strings.Join(response.DistinguishedName.Country, ","),
			State:              strings.Join(response.DistinguishedName.State, ","),
			Locality:           strings.Join(response.DistinguishedName.Locality, ","),
			Organization:       strings.Join(response.DistinguishedName.Organization, ","),
			OrganizationalUnit: strings.Join(response.DistinguishedName.OrganizationalUnit, ","),
			EmailAddress:       strings.Join(response.DistinguishedName.EmailAddress, ","),
		},
		IssuerName:              response.IssuerName,
		NotBefore:               response.NotBefore,
		NotAfter:                response.NotAfter,
		SubjectAlternativeNames: response.SubjectAlternativeNames,
		KeyUsages:               response.KeyUsages,
		Raw:                     response.Raw,
	}
}
