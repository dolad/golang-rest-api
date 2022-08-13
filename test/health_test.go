//go:build e2e
// +build e2e

package test

import (
	"fmt"
	resty "github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	fmt.Println("Runnig E2E test for health check endpoint")

	client := resty.New()
	resp, err := client.R().Get("http://localhost:8080/api/health")
	if err != nil {
		t.Fail()
	}
	fmt.Println(resp.StatusCode())
	assert.Equal(t, 200, resp.StatusCode())
}
