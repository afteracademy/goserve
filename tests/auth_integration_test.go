package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afteracademy/goserve/api/auth/model"
	"github.com/afteracademy/goserve/arch/network"
	"github.com/afteracademy/goserve/startup"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationAuthController_SignupSuccess(t *testing.T) {
	router, module, shutdown := startup.TestServer()
	defer shutdown()

	apikey, err := module.GetInstance().AuthService.CreateApiKey("test_key", 1, []model.Permission{"test"}, []string{"comment"})
	if err != nil {
		t.Fatalf("could not create apikey: %v", err)
	}

	body := `{"email":"test@abc.com","password":"123456","name":"test name"}`

	req, err := http.NewRequest("POST", "/auth/signup/basic", bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(network.ApiKeyHeader, apikey.Key)

	rr := httptest.NewRecorder()
	router.GetEngine().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"success"`)
	assert.Contains(t, rr.Body.String(), `"data"`)
	assert.Contains(t, rr.Body.String(), `"user"`)
	assert.Contains(t, rr.Body.String(), `"roles"`)
	assert.Contains(t, rr.Body.String(), `"tokens"`)

	_, err = module.GetInstance().AuthService.DeleteApiKey(apikey)
	if err != nil {
		t.Fatalf("could not delete apikey: %v", err)
	}

	_, err = module.GetInstance().UserService.RemoveUserByEmail("test@abc.com")
	if err != nil {
		t.Fatalf("could not delete user: %v", err)
	}
}
