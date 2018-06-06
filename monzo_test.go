package monzo_test

import (
	"github.com/Lavoaster/monzo-go"
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestClient_NewRequest(t *testing.T) {
	monzoClient := monzo.New("randomKey")
	req, err := monzoClient.NewRequest("GET", "/toaster", nil)

	if err != nil {
		t.Error(err)
	}

	// Test header got set

	authorization := req.Header.Get("Authorization")

	assert.NotNil(t, authorization)
	assert.Equal(t, "Bearer randomKey", authorization)

	// Test endpoint is generated correctly

	assert.Equal(t, "https://api.monzo.com/toaster", req.URL.String())
}
