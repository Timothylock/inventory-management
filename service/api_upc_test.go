package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Timothylock/inventory-management/config"
	"github.com/Timothylock/inventory-management/items"
	"github.com/Timothylock/inventory-management/upc"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLookupBarcodeMissingQuery(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()

	ip := items.NewMockPersister(mc)
	server := setupServerAuthenticated(ip, t)
	defer server.Close()

	resp, err := sendGet(server.URL + "/api/lookup")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestLookupBarcodeServerError(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()

	ip := items.NewMockPersister(mc)
	server := setupServerAuthenticated(ip, t)
	defer server.Close()

	resp, err := sendGet(server.URL + "/api/lookup?barcode=123")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestLookupBarcodeSuccess(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"product":{"name":"someName","image_url":"someURL","category":{"name":"categoryName"}}}`)
	}))
	defer ts.Close()

	cfg := config.Config{
		UpcUrl: ts.URL,
	}
	ip := items.NewMockPersister(mc)
	server := setupServerWithConfigAuthenticated(ip, cfg, t)
	defer server.Close()

	expected := upc.ItemDetail{
		ID:         "123",
		Name:       "someName",
		PictureURL: "someURL",
		Category:   "categoryName",
	}
	expectedStr, err := json.Marshal(expected)
	assert.NoError(t, err)

	resp, err := sendGet(server.URL + "/api/lookup?barcode=123")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.JSONEq(t, string(expectedStr), string(getBody(t, resp)))
}
