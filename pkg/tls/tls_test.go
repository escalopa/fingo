package tls

import (
	"reflect"
	"testing"

	"github.com/madflojo/testcerts"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func TestLoadClientTLS(t *testing.T) {
	cert, key, err := testcerts.GenerateCertsToTempFile("/tmp")
	require.NoError(t, err)

	tests := []struct {
		name      string
		enabled   string
		transport credentials.TransportCredentials
	}{
		{
			name:      "success on enabled",
			enabled:   "true",
			transport: getCreds(t, cert, key),
		},
		{
			name:      "success on disabled",
			enabled:   "false",
			transport: insecure.NewCredentials(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr1, err := LoadClientTLS(tt.enabled, cert)
			require.NoError(t, err)
			// validate that the cert and key are valid
			tr2, err := credentials.NewClientTLSFromFile(cert, key)
			require.NoError(t, err)
			if reflect.DeepEqual(tr1, tr2) {
				t.Errorf("LoadClientTLS() = %v, want %v", tr1, tr2)
			}
		})
	}
}

func getCreds(t *testing.T, certFile, keyFile string) credentials.TransportCredentials {
	cred, err := credentials.NewClientTLSFromFile(certFile, keyFile)
	require.NoError(t, err)
	return cred
}
