package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// HTTPClient 封装http.Client
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient 创建HTTPClient实例
func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// RequestParams 请求信息
type RequestParams struct {
	Method      string
	URL         string
	Headers     map[string]string
	Body        interface{}
	QueryParams map[string]string
}

// SendRequest 执行Http请求的发送
func (hc *HTTPClient) SendRequest(params RequestParams) ([]byte, error) {
	// 添加请求参数
	if len(params.QueryParams) > 0 {
		parsedURL, err := url.Parse(params.URL)
		if err != nil {
			return nil, fmt.Errorf("URL解析失败: %w", err)
		}
		query := parsedURL.Query()
		for key, value := range params.QueryParams {
			query.Add(key, value)
		}
		parsedURL.RawQuery = query.Encode()
		params.URL = parsedURL.String()
	}

	// 将请求体序列化
	var body []byte
	if params.Body != nil {
		var err error
		body, err = json.Marshal(params.Body)
		if err != nil {
			return nil, fmt.Errorf("请求体序列化失败: %w", err)
		}
	}

	// 创建HTTP请求
	req, err := http.NewRequest(params.Method, params.URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("HTTP请求创建失败: %w", err)
	}

	// 添加请求头
	for key, value := range params.Headers {
		req.Header.Set(key, value)
	}

	// 若请求体不为空，将Content-Type设为application/json
	if params.Body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	// 发送http请求
	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求发送失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP请求状态
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	// Read 读取请求结果
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("请求结果读取失败: %w", err)
	}

	return respBody, nil
}

// Get 发送Get请求
func (hc *HTTPClient) Get(url string, headers map[string]string, queryParams map[string]string) ([]byte, error) {
	return hc.SendRequest(RequestParams{
		Method:      http.MethodGet,
		URL:         url,
		Headers:     headers,
		QueryParams: queryParams,
	})
}

// Post 发送Post请求
func (hc *HTTPClient) Post(url string, headers map[string]string, body interface{}) ([]byte, error) {
	return hc.SendRequest(RequestParams{
		Method:  http.MethodPost,
		URL:     url,
		Headers: headers,
		Body:    body,
	})
}

// Put 发送Put请求
func (hc *HTTPClient) Put(url string, headers map[string]string, body interface{}) ([]byte, error) {
	return hc.SendRequest(RequestParams{
		Method:  http.MethodPut,
		URL:     url,
		Headers: headers,
		Body:    body,
	})
}

// Delete 发送Delete请求
func (hc *HTTPClient) Delete(url string, headers map[string]string, queryParams map[string]string) ([]byte, error) {
	return hc.SendRequest(RequestParams{
		Method:      http.MethodDelete,
		URL:         url,
		Headers:     headers,
		QueryParams: queryParams,
	})
}
