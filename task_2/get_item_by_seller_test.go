package task_2

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetSellerItem_Positive(t *testing.T) {
	validItem := ItemRequest{
		SellerID: GenerateID(),
		Name:     "телефон",
		Price:    700,
		Statistics: Statistics{
			Contacts:  1,
			Likes:     1,
			ViewCount: 1,
		},
	}

	resp := postItem(t, validItem)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	url := fmt.Sprintf("%s/api/1/%d/item", BaseUrl, validItem.SellerID)
	getResp, err := http.Get(url)
	assert.NoError(t, err)
	defer getResp.Body.Close()
	assert.Equal(t, 200, getResp.StatusCode)

	body, err := io.ReadAll(getResp.Body)
	assert.NoError(t, err)

	layout := "2006-01-02 15:04:05.999999 +0300 +0300"

	var getResult []ItemResponse
	err = json.Unmarshal(body, &getResult)
	assert.NoError(t, err)
	assert.NotEmpty(t, getResult)

	for _, item := range getResult {
		if item.SellerID == validItem.SellerID {
			_, err = time.Parse(layout, item.CreatedAt)
			assert.NoError(t, err)
			assert.Equal(t, validItem.Name, item.Name)
			assert.Equal(t, validItem.Price, item.Price)
			statsExpected := validItem.Statistics
			statsGet := item.Statistics
			if statsGet != nil {
				assert.Equal(t, statsExpected.Contacts, statsGet.Contacts)
				assert.Equal(t, statsExpected.Likes, statsGet.Likes)
				assert.Equal(t, statsExpected.ViewCount, statsGet.ViewCount)
			}
		}
	}

}

func TestGetSellerItem_Empty(t *testing.T) {
	noExistID := GenerateID()

	url := fmt.Sprintf("%s/api/1/%d/item", BaseUrl, noExistID)
	resp, err := http.Get(url)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var items []ItemResponse
	err = json.Unmarshal(body, &items)
	assert.NoError(t, err)
	assert.Empty(t, items)

}

func TestGetSellerItem_InvalidSellerID(t *testing.T) {
	tests := []struct {
		name         string
		sellerID     interface{}
		expectedCode int
	}{
		{
			name:         "negative id",
			sellerID:     -100,
			expectedCode: 400,
		},
		{
			name:         "zero id",
			sellerID:     0,
			expectedCode: 400,
		},
		{
			name:         "out of range id",
			sellerID:     999999999999,
			expectedCode: 400,
		},
		{
			name:         "string id",
			sellerID:     "123",
			expectedCode: 400,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Skip("Баг описан в BUGS.md")
			url := fmt.Sprintf("%s/api/1/%v/item", BaseUrl, test.sellerID)
			resp, err := http.Get(url)
			assert.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, test.expectedCode, resp.StatusCode)
		})
	}
}
