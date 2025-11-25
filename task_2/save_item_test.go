package task_2

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const BaseUrl = "https://qa-internship.avito.com"

func TestSaveItem_Positive(t *testing.T) {
	t.Skip("Баг 2 описан в BUGS.md")
	validItem := ItemRequest{
		SellerID: GenerateID(),
		Name:     "телефон",
		Price:    700,
		Statistics: Statistics{
			Contacts:  32,
			Likes:     1,
			ViewCount: 14,
		},
	}

	resp := postItem(t, validItem)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)

	var response ItemResponse

	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.ID)
	assert.Equal(t, validItem.Name, response.Name)
	assert.Equal(t, validItem.Price, response.Price)
	assert.Equal(t, validItem.SellerID, response.SellerID)

}

func TestSaveItem_InvalidPrice(t *testing.T) {
	tests := []struct {
		name         string
		item         map[string]interface{}
		expectedCode int
	}{
		{
			name: "negative price",
			item: map[string]interface{}{
				"sellerID":   GenerateID(),
				"name":       "телефон",
				"price":      -100,
				"statistics": map[string]int{"likes": 1, "viewCount": 1, "contacts": 1},
			},
			expectedCode: 400,
		},
		{
			name: "zero price",
			item: map[string]interface{}{
				"sellerID":   GenerateID(),
				"name":       "телефон",
				"price":      0,
				"statistics": map[string]int{"likes": 1, "viewCount": 1, "contacts": 1},
			},
			expectedCode: 400,
		},
		{
			name: "nil price",
			item: map[string]interface{}{
				"sellerID":   GenerateID(),
				"name":       "телефон",
				"price":      nil,
				"statistics": map[string]int{"likes": 1, "viewCount": 1, "contacts": 1},
			},
			expectedCode: 400,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.name == "negative price" {
				t.Skip("баг API, см. BUGS.md")
			}
			resp := postItem(t, test.item)
			defer resp.Body.Close()
			assert.Equal(t, test.expectedCode, resp.StatusCode)

		})
	}

}

func GenerateID() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return 111111 + r.Intn(888889)
}

func TestSaveItem_InvalidStatistics(t *testing.T) {
	tests := []struct {
		name         string
		item         map[string]interface{}
		expectedCode int
	}{
		{
			name: "negative likes",
			item: map[string]interface{}{
				"sellerID":   GenerateID(),
				"name":       "телефон",
				"price":      1,
				"statistics": map[string]int{"likes": -1, "viewCount": 1, "contacts": 1},
			},
			expectedCode: 400,
		},
		{
			name: "negative viewCount",
			item: map[string]interface{}{
				"sellerID":   GenerateID(),
				"name":       "телефон",
				"price":      1,
				"statistics": map[string]int{"likes": 1, "viewCount": -1, "contacts": 1},
			},
			expectedCode: 400,
		},
		{
			name: "negative contacts",
			item: map[string]interface{}{
				"sellerID":   GenerateID(),
				"name":       "телефон",
				"price":      1,
				"statistics": map[string]int{"likes": 1, "viewCount": 1, "contacts": -1},
			},
			expectedCode: 400,
		},
		{
			name: "nil likes",
			item: map[string]interface{}{
				"sellerID":   GenerateID(),
				"name":       "телефон",
				"price":      1,
				"statistics": map[string]interface{}{"likes": nil, "viewCount": 1, "contacts": 1},
			},
			expectedCode: 400,
		},
		{
			name: "nil statistics",
			item: map[string]interface{}{
				"sellerID":   GenerateID(),
				"name":       "телефон",
				"price":      1,
				"statistics": nil,
			},
			expectedCode: 400,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.name == "negative likes" || test.name == "negative viewCount" || test.name == "negative contacts" {
				t.Skip("Баг описан в BUGS.md")
			}

			resp := postItem(t, test.item)
			defer resp.Body.Close()
			assert.Equal(t, test.expectedCode, resp.StatusCode)
		})

	}
}

func TestSaveItem_MissingFields(t *testing.T) {
	tests := []struct {
		name         string
		item         map[string]interface{}
		expectedCode int
	}{
		{
			name: "missing name",
			item: map[string]interface{}{
				"sellerID":   GenerateID(),
				"price":      1,
				"statistics": map[string]int{"likes": 1, "viewCount": 1, "contacts": 1},
			},
			expectedCode: 400,
		},
		{
			name: "missing price",
			item: map[string]interface{}{
				"sellerID":   GenerateID(),
				"name":       "телеофн",
				"statistics": map[string]int{"likes": 1, "viewCount": 1, "contacts": 1},
			},
			expectedCode: 400,
		},
		{
			name: "missing sellerID",
			item: map[string]interface{}{
				"price":      1,
				"name":       "телефон",
				"statistics": map[string]int{"likes": 1, "viewCount": 1, "contacts": 1},
			},
			expectedCode: 400,
		},
		{
			name: "missing statistics",
			item: map[string]interface{}{
				"sellerID": GenerateID(),
				"price":    1,
				"name":     "телефон",
			},
			expectedCode: 400,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			resp := postItem(t, test.item)
			defer resp.Body.Close()
			assert.Equal(t, test.expectedCode, resp.StatusCode)
		})

	}
}

func postItem(t *testing.T, item interface{}) *http.Response {
	jsonData, err := json.Marshal(item)
	assert.NoError(t, err)

	resp, err := http.Post(BaseUrl+"/api/1/item", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)

	return resp
}
