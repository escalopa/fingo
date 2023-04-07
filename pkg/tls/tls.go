package grpctls

import (
	"log"

	"github.com/lordvidex/errs"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func LoadServerTLS(enabled string, certFile, keyFile string) (credentials.TransportCredentials, error) {
	log.Println("creating grpc server with TLS:", enabled)
	if enabled == "true" {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			return insecure.NewCredentials(), errs.B(err).Msg("failed to load TLS certificates").Err()
		}
		log.Println("loaded TLS certificates")
		return creds, nil
	} else {
		return insecure.NewCredentials(), nil
	}
}

func LoadClientTLS(enabled string, certFile string) (credentials.TransportCredentials, error) {
	log.Println("connecting to grpc server with TLS:", enabled)
	if enabled == "true" {
		creds, err := credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			return insecure.NewCredentials(), errs.B(err).Msg("failed to load TLS certificates").Err()
		}
		log.Println("loaded TLS certificates")
		return creds, nil
	} else {
		return insecure.NewCredentials(), nil
	}
}
