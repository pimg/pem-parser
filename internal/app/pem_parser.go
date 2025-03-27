package app

import (
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/smallstep/certinfo"
)

type PEMHandler struct{}

type PEMResponse struct {
	SerialNumber            string
	DistinguishedName       *DistinguishedName
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
	Country            []string
	State              []string
	Locality           []string
	Organization       []string
	OrganizationalUnit []string
	EmailAddress       []string
	Short              string
}

type PublicKey struct {
	Fingerprint string
	Type        string
}

func NewPEMHandler() *PEMHandler {
	return &PEMHandler{}
}

func (h *PEMHandler) Handle(pemRaw []byte) ([]*PEMResponse, error) {
	chain := make([]*PEMResponse, 0)
	var certDERBlock *pem.Block
	for {
		certDERBlock, pemRaw = pem.Decode(pemRaw)
		if certDERBlock == nil && len(pemRaw) > 0 {
			return nil, errors.New("failed to decode PEM block, the submitted data does not seem to be PEM encoded")
		}
		if certDERBlock == nil {
			break
		}
		pemResponse, err := h.parsePEMBlock(certDERBlock)
		if err != nil {
			return nil, err
		}
		chain = append(chain, pemResponse)
	}

	if len(chain) < 1 {
		return nil, errors.New("failed to decode PEM block, the submitted data does not seem to be PEM encoded")
	}

	return chain, nil
}

func (h *PEMHandler) parsePEMBlock(block *pem.Block) (*PEMResponse, error) {
	pemResponse := &PEMResponse{}
	switch block.Type {
	case "CERTIFICATE":
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			slog.Error("failed to parse certificate", "error", err)
			return nil, err
		}
		pemResponse.Raw, err = certinfo.CertificateText(cert)
		if err != nil {
			slog.Error("failed to render certificate text", "error", err)
			return nil, err
		}

		pemResponse.SerialNumber = cert.SerialNumber.String()
		pemResponse.IssuerName = cert.Issuer.CommonName
		pemResponse.NotBefore = cert.NotBefore
		pemResponse.NotAfter = cert.NotAfter
		pemResponse.DistinguishedName = mapSubject(cert.Subject)
		pemResponse.SubjectAlternativeNames = cert.DNSNames
		pemResponse.KeyUsages = parseKeyUsage(cert.KeyUsage)
		pemResponse.Fingerprint = fingerprint(cert.Raw)
		pemResponse.PublicKey = &PublicKey{
			Fingerprint: fingerprint(cert.RawSubjectPublicKeyInfo),
			Type:        cert.PublicKeyAlgorithm.String(),
		}
	case "CERTIFICATE REQUEST":
		csr, err := x509.ParseCertificateRequest(block.Bytes)
		if err != nil {
			slog.Error("failed to parse certificate request", "error", err)
			return nil, err
		}

		pemResponse.Raw, err = certinfo.CertificateRequestText(csr)
		if err != nil {
			slog.Error("failed to render certificate request text", "error", err)
			return nil, err
		}

		pemResponse.DistinguishedName = mapSubject(csr.Subject)
		pemResponse.SubjectAlternativeNames = csr.DNSNames
		pemResponse.Fingerprint = fingerprint(csr.RawTBSCertificateRequest)
		pemResponse.PublicKey = &PublicKey{
			Fingerprint: fingerprint(csr.RawSubjectPublicKeyInfo),
			Type:        csr.PublicKeyAlgorithm.String(),
		}
	default:
		slog.Info("unknown certificate type", "type", block.Type)
		if strings.Contains(block.Type, "PRIVATE") {
			return nil, errors.New("you have submitted a private key! \neven though we do not store any PEM files, you should consider this private key compromised")
		}
		return nil, errors.New("unsupported PEM type")
	}

	pemResponse.Type = block.Type
	return pemResponse, nil
}

func mapSubject(sub pkix.Name) *DistinguishedName {
	return &DistinguishedName{
		CommonName:         sub.CommonName,
		SerialNumber:       sub.SerialNumber,
		Country:            sub.Country,
		State:              sub.Province,
		Locality:           sub.Locality,
		Organization:       sub.Organization,
		OrganizationalUnit: sub.OrganizationalUnit,
		Short:              sub.String(),
	}
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

func fingerprint(raw []byte) string {
	cs := sha256.Sum256(raw)
	return hex.EncodeToString(cs[:])
}
