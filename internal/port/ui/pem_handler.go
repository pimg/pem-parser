package ui

import (
	"errors"
	"fmt"
	"net/http"
	"pem-parser/internal/app"
	"strings"
	"time"
)

type PEMParserPage struct {
	SuccessMessage string
	ErrorMessage   string
	Result         *Result
}

type Result struct {
	Leaf  *PEMResponse
	Chain []*PEMResponse
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
	Fingerprint             string
	PublicKey               *PublicKey
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
	Short              string
}

type PublicKey struct {
	Fingerprint string
	Type        string
}

func (s *Server) pemParserHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("received pem parser request")

	page := &PEMParserPage{}

	pem, err := parseForm(r)
	var maxErr *http.MaxBytesError
	if errors.As(err, &maxErr) {
		page.ErrorMessage = fmt.Sprintf("request is too large, max supported limit is %d MB", maxErr.Limit/(1000*1000))
		s.RenderResultPage(w, page)
		return
	}

	if err != nil {
		page.ErrorMessage = "failed to parse PEM"
		s.RenderResultPage(w, page)
		return
	}

	out, err := s.App.PEMHandler.Handle(pem)
	if err != nil {
		page.ErrorMessage = err.Error()
		s.RenderResultPage(w, page)
		return
	}

	if len(out) > 0 {
		page.SuccessMessage = fmt.Sprintf("Successfully parsed PEM %s file", out[0].Type) // TODO consider changing this
		chain := mapResponse(out)
		page.Result = &Result{
			Leaf:  chain[len(chain)-1],
			Chain: chain,
		}
	}

	s.RenderResultPage(w, page)
}

func (s *Server) RenderResultPage(w http.ResponseWriter, page *PEMParserPage) {
	err := s.Templates.ExecuteTemplate(w, "result-block", page)
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

func mapResponse(handlerOut []*app.PEMResponse) []*PEMResponse {
	response := make([]*PEMResponse, len(handlerOut))
	for i, block := range handlerOut {
		response[len(handlerOut)-i-1] = mapPEMResponse(block)
	}
	return response
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
			Short:              response.DistinguishedName.Short,
		},
		IssuerName:              response.IssuerName,
		NotBefore:               response.NotBefore,
		NotAfter:                response.NotAfter,
		SubjectAlternativeNames: response.SubjectAlternativeNames,
		KeyUsages:               response.KeyUsages,
		Raw:                     response.Raw,
		Fingerprint:             response.Fingerprint,
		PublicKey: &PublicKey{
			Fingerprint: response.PublicKey.Fingerprint,
			Type:        response.PublicKey.Type,
		},
	}
}
