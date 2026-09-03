package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEphServiceClient(t *testing.T) {
	client, err := NewEphServiceClient("localhost:5678")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NoError(t, client.Close())
}
