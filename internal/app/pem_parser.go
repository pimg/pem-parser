package app

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log/slog"
	"strings"

	"github.com/smallstep/certinfo"
)

type PEMHandler struct {
	logger *slog.Logger
}

type PEMResponse struct {
	Result string
	Type   string
}

func NewPEMHandler(logger *slog.Logger) *PEMHandler {
	return &PEMHandler{logger: logger}
}

func (h *PEMHandler) Handle(pemRaw []byte) (*PEMResponse, error) {
	block, rest := pem.Decode(pemRaw)
	if block == nil || len(rest) > 0 {
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
		pemResponse.Result, err = certinfo.CertificateText(cert)
		if err != nil {
			h.logger.Error("failed to render certificate text", "error", err)
			return nil, err
		}
	case "CERTIFICATE REQUEST":
		csr, err := x509.ParseCertificateRequest(block.Bytes)
		if err != nil {
			h.logger.Error("failed to parse certificate request", "error", err)
			return nil, err
		}

		pemResponse.Result, err = certinfo.CertificateRequestText(csr)
		if err != nil {
			h.logger.Error("failed to render certificate request text", "error", err)
			return nil, err
		}
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
