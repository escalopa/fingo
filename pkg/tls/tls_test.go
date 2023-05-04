package tls

import (
	"reflect"
	"testing"

	"github.com/madflojo/testcerts"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func TestLoadServerTLS(t *testing.T) {
	cert, key, err := testcerts.GenerateCertsToTempFile("/tmp")
	require.NoError(t, err)

	tests := []struct {
		name    string
		enabled bool
	}{
		{
			name:    "success on enabled",
			enabled: true,
		},
		{
			name:    "success on disabled",
			enabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cred1, err := LoadServerTLS(tt.enabled, cert, key)
			require.NoError(t, err)
			// validate that the cert and key are valid
			cred2, err := credentials.NewServerTLSFromFile(cert, key)
			require.NoError(t, err)
			if reflect.DeepEqual(cred1, insecure.NewCredentials()) && tt.enabled {
				t.Errorf("LoadServerTLS() = %v, want %v", cred1, cred2)
			}
		})
	}

}

func TestLoadClientTLS(t *testing.T) {
	cert, _, err := testcerts.GenerateCertsToTempFile("/tmp")
	require.NoError(t, err)

	tests := []struct {
		name    string
		enabled bool
	}{
		{
			name:    "success on enabled",
			enabled: true,
		},
		{
			name:    "success on disabled",
			enabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cred1, err := LoadClientTLS(tt.enabled, cert)
			require.NoError(t, err)
			// validate that the cert and key are valid
			cred2, err := credentials.NewClientTLSFromFile(cert, "")
			require.NoError(t, err)

			if reflect.DeepEqual(cred1, insecure.NewCredentials()) && tt.enabled {
				t.Errorf("LoadClientTLS() = %v, want %v", cred1, cred2)
			}
		})
	}
}
