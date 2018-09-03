package upc

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Timothylock/inventory-management/config"
)

type Service struct {
	config config.Config
	client *http.Client
}

func NewService(c config.Config) Service {
	return Service{
		config: c,
		client: &http.Client{
			Timeout: 10000 * time.Millisecond,
		},
	}
}

type LookupAPIResponse struct {
	Product struct {
		Name     string `json:"name"`
		ImageURL string `json:"image_url"`
		Category struct {
			Name string `json:"name"`
		} `json:"category"`
	} `json:"product"`
}

type ItemDetail struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	PictureURL string `json:"pictureURL"`
}

type LookupBody struct {
	BarcodeNumber string `json:"barcode_number"`
	Token         string `json:"token"`
}

func (s *Service) LookupBarcode(barcode string) (ItemDetail, error) {
	payload := LookupBody{
		BarcodeNumber: barcode,
		Token:         s.config.UpcToken,
	}

	payloadStr, err := json.Marshal(payload)
	if err != nil {
		return ItemDetail{}, err
	}

	req, err := http.NewRequest("POST", s.config.UpcUrl, strings.NewReader(string(payloadStr)))
	if err != nil {
		return ItemDetail{}, err
	}

	req.Header.Add("user-agent", "Dalvik/2.1.0")
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Add("accept-encoding", "gzip")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	req.Header.Add("content-type", "application/json")

	res, err := s.client.Do(req)
	if err != nil {
		return ItemDetail{}, err
	}
	defer res.Body.Close()

	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(res.Body)
		defer reader.Close()
	default:
		reader = res.Body
	}
	defer reader.Close()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return ItemDetail{}, err
	}

	var apiReponse = new(LookupAPIResponse)
	if err = json.Unmarshal(body, &apiReponse); err != nil {
		return ItemDetail{}, err
	}

	return ItemDetail{
		ID:         barcode,
		Name:       apiReponse.Product.Name,
		Category:   apiReponse.Product.Category.Name,
		PictureURL: apiReponse.Product.ImageURL,
	}, nil
}
