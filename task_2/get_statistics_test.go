package task_2

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetStatistics_Positive(t *testing.T) {
	validID := "ab0f3148-4d0e-479f-9f26-cebe00aa8047"

	getResp, err := http.Get(BaseUrl + "/api/1/statistic/" + validID)
	assert.NoError(t, err)
	defer getResp.Body.Close()
	assert.Equal(t, 200, getResp.StatusCode)

	body, err := io.ReadAll(getResp.Body)
	assert.NoError(t, err)

	var getResult []Statistics
	err = json.Unmarshal(body, &getResult)
	assert.NoError(t, err)
	assert.NotEmpty(t, getResult)

	for _, item := range getResult {
		assert.NotNil(t, item.Likes)
		assert.NotNil(t, item.ViewCount)
		assert.NotNil(t, item.Contacts)
	}
}

func TestGetStatistics_NotFound(t *testing.T) {
	noExistID := uuid.New().String()

	getResp, err := http.Get(BaseUrl + "/api/1/statistic/" + noExistID)
	assert.NoError(t, err)
	defer getResp.Body.Close()
	assert.Equal(t, 404, getResp.StatusCode)

	body, err := io.ReadAll(getResp.Body)
	assert.NoError(t, err)

	var errResp map[string]interface{}
	err = json.Unmarshal(body, &errResp)
	assert.NoError(t, err)

	assert.Contains(t, errResp, "status")
	assert.Contains(t, errResp, "result")
}

func TestGetStatistics_InvalidID(t *testing.T) {
	tests := []struct {
		name         string
		itemID       interface{}
		expectedCode int
	}{
		{
			name:         "integer id",
			itemID:       1223456,
			expectedCode: 400,
		},
		{
			name:         "string (not uuid)",
			itemID:       "1232-1232-1232",
			expectedCode: 400,
		},
		{
			name:         "too long id",
			itemID:       "debb6473-0cc5-4204-9f72-f31534d6f03f-aaaa",
			expectedCode: 400,
		},
		{
			name:         "empty id",
			itemID:       "",
			expectedCode: 400,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.name == "empty id" {
				t.Skip("баг API, см. BUGS.md")
			}
			id := fmt.Sprintf("%v", test.itemID)
			getResp, err := http.Get(BaseUrl + "/api/1/statistic/" + id)
			assert.NoError(t, err)
			defer getResp.Body.Close()
			assert.Equal(t, test.expectedCode, getResp.StatusCode)
		})
	}
}
