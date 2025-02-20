package app

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/smallstep/certinfo"
)

type PEMHandler struct {
	logger *slog.Logger
}

type PEMResponse struct {
	SerialNumber            string
	DistinguishedName       DistinguishedName
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
	Country            []string
	State              []string
	Locality           []string
	Organization       []string
	OrganizationalUnit []string
	EmailAddress       []string
}

func NewPEMHandler(logger *slog.Logger) *PEMHandler {
	return &PEMHandler{logger: logger}
}

func (h *PEMHandler) Handle(pemRaw []byte) (*PEMResponse, error) {
	block, rest := pem.Decode(pemRaw)
	if block == nil {
		return nil, errors.New("failed to decode PEM block, the submitted data does not seem to be PEM encoded")
	}

	if len(rest) > 0 {
		return nil, errors.New("PEM contains more than one PEM block, this is not yet supported")
	}

	pemResponse := &PEMResponse{}
	switch block.Type {
	case "CERTIFICATE":
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			h.logger.Error("failed to parse certificate", "error", err)
			return nil, err
		}
		pemResponse.Raw, err = certinfo.CertificateText(cert)
		if err != nil {
			h.logger.Error("failed to render certificate text", "error", err)
			return nil, err
		}

		pemResponse.SerialNumber = cert.SerialNumber.String()
		pemResponse.IssuerName = cert.Issuer.CommonName
		pemResponse.NotBefore = cert.NotBefore
		pemResponse.NotAfter = cert.NotAfter
		pemResponse.DistinguishedName = DistinguishedName{
			CommonName:         cert.Subject.CommonName,
			SerialNumber:       cert.Subject.SerialNumber,
			Country:            cert.Subject.Country,
			State:              cert.Subject.Province,
			Locality:           cert.Subject.Locality,
			Organization:       cert.Subject.Organization,
			OrganizationalUnit: cert.Subject.OrganizationalUnit,
		}
		pemResponse.SubjectAlternativeNames = cert.DNSNames
		pemResponse.KeyUsages = parseKeyUsage(cert.KeyUsage)
	case "CERTIFICATE REQUEST":
		csr, err := x509.ParseCertificateRequest(block.Bytes)
		if err != nil {
			h.logger.Error("failed to parse certificate request", "error", err)
			return nil, err
		}

		pemResponse.Raw, err = certinfo.CertificateRequestText(csr)
		if err != nil {
			h.logger.Error("failed to render certificate request text", "error", err)
			return nil, err
		}
		pemResponse.DistinguishedName = DistinguishedName{
			CommonName:         csr.Subject.CommonName,
			SerialNumber:       csr.Subject.SerialNumber,
			Country:            csr.Subject.Country,
			State:              csr.Subject.Province,
			Locality:           csr.Subject.Locality,
			Organization:       csr.Subject.Organization,
			OrganizationalUnit: csr.Subject.OrganizationalUnit,
		}
		pemResponse.SubjectAlternativeNames = csr.DNSNames
	default:
		h.logger.Info("unknown certificate type", "type", block.Type)
		if strings.Contains(block.Type, "PRIVATE") {
			return nil, errors.New("You have submitted a private key! \nEven though we do not store any PEM files, you should consider this private key compromised.")
		}
		return nil, errors.New("unsupported PEM type")
	}

	pemResponse.Type = block.Type

	return pemResponse, nil
}

func parseKeyUsage(usage x509.KeyUsage) []string {
	var keyUsages []string
	if usage&x509.KeyUsageDigitalSignature != 0 {
		keyUsages = append(keyUsages, "DigitalSignature")
	}
	if usage&x509.KeyUsageContentCommitment != 0 {
		keyUsages = append(keyUsages, "ContentCommitment")
	}
	if usage&x509.KeyUsageKeyEncipherment != 0 {
		keyUsages = append(keyUsages, "KeyEncipherment")
	}
	if usage&x509.KeyUsageDataEncipherment != 0 {
		keyUsages = append(keyUsages, "DataEncipherment")
	}
	if usage&x509.KeyUsageKeyAgreement != 0 {
		keyUsages = append(keyUsages, "KeyAgreement")
	}
	if usage&x509.KeyUsageCertSign != 0 {
		keyUsages = append(keyUsages, "CertSign")
	}
	if usage&x509.KeyUsageCRLSign != 0 {
		keyUsages = append(keyUsages, "CRLSign")
	}
	if usage&x509.KeyUsageEncipherOnly != 0 {
		keyUsages = append(keyUsages, "EncipherOnly")
	}
	if usage&x509.KeyUsageDecipherOnly != 0 {
		keyUsages = append(keyUsages, "DecipherOnly")
	}
	return keyUsages
}
