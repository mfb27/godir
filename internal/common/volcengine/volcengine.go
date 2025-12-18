package volcengine

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Client struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Endpoint        string
	HTTPClient      *http.Client
}

func NewClient(accessKeyID, secretAccessKey, region, endpoint string) *Client {
	return &Client{
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		Region:          region,
		Endpoint:        endpoint,
		HTTPClient:      &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) signRequest(method, path string, query url.Values, body []byte) string {
	timestamp := time.Now().UTC().Format("20060102T150405Z")
	date := timestamp[:8]

	canonicalURI := path
	if canonicalURI == "" {
		canonicalURI = "/"
	}

	canonicalQueryString := ""
	if len(query) > 0 {
		keys := make([]string, 0, len(query))
		for k := range query {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			parts = append(parts, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(query.Get(k))))
		}
		canonicalQueryString = strings.Join(parts, "&")
	}

	canonicalHeaders := fmt.Sprintf("host:%s\nx-date:%s\n", c.Endpoint, timestamp)
	signedHeaders := "host;x-date"

	payloadHash := sha256Hash(body)
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		method, canonicalURI, canonicalQueryString, canonicalHeaders, signedHeaders, payloadHash)

	algorithm := "HMAC-SHA256"
	credentialScope := fmt.Sprintf("%s/%s/request", date, c.Region)
	stringToSign := fmt.Sprintf("%s\n%s\n%s\n%s",
		algorithm, timestamp, credentialScope, sha256Hash([]byte(canonicalRequest)))

	kDate := hmacSHA256([]byte("VOLC"+c.SecretAccessKey), date)
	kRegion := hmacSHA256(kDate, c.Region)
	kService := hmacSHA256(kRegion, "volcengine")
	kSigning := hmacSHA256(kService, "request")
	signature := hex.EncodeToString(hmacSHA256(kSigning, stringToSign))

	return fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		algorithm, c.AccessKeyID, credentialScope, signedHeaders, signature)
}

func (c *Client) Request(method, path string, query url.Values, body interface{}) (*http.Response, error) {
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	reqURL := fmt.Sprintf("https://%s%s", c.Endpoint, path)
	if len(query) > 0 {
		reqURL += "?" + query.Encode()
	}

	req, err := http.NewRequest(method, reqURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	timestamp := time.Now().UTC().Format("20060102T150405Z")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Date", timestamp)
	req.Header.Set("Host", c.Endpoint)

	authorization := c.signRequest(method, path, query, bodyBytes)
	req.Header.Set("Authorization", authorization)

	return c.HTTPClient.Do(req)
}

func sha256Hash(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func (c *Client) Call(action string, version string, body interface{}) (map[string]interface{}, error) {
	query := url.Values{}
	query.Set("Action", action)
	query.Set("Version", version)

	resp, err := c.Request("POST", "/", query, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w, body: %s", err, string(respBody))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API错误: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	return result, nil
}
