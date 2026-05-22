package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type StatusRecord struct {
	ID         int      `json:"id"`
	ProjectName string  `json:"project_name"`
	ShortName   string  `json:"short_name"`
	Status      string  `json:"status"`
	Phase       *string `json:"phase,omitempty"`
	Summary     string  `json:"summary"`
	Reason      *string `json:"reason,omitempty"`
	Details     *string `json:"details,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Source      *string `json:"source,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type StatusRecordCreate struct {
	ProjectName string   `json:"project_name"`
	ShortName   string   `json:"short_name"`
	Status      string   `json:"status"`
	Phase       *string  `json:"phase,omitempty"`
	Summary     string   `json:"summary"`
	Reason      *string  `json:"reason,omitempty"`
	Details     *string  `json:"details,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Source      *string  `json:"source,omitempty"`
}

type StatusRecordUpdate struct {
	ProjectName *string  `json:"project_name,omitempty"`
	ShortName   *string  `json:"short_name,omitempty"`
	Status      *string  `json:"status,omitempty"`
	Phase       *string  `json:"phase,omitempty"`
	Summary     *string  `json:"summary,omitempty"`
	Reason      *string  `json:"reason,omitempty"`
	Details     *string  `json:"details,omitempty"`
	Tags        *[]string `json:"tags,omitempty"`
	Source      *string   `json:"source,omitempty"`
}

type ListResponse struct {
	Items     []StatusRecord `json:"items"`
	Total     int            `json:"total"`
	Page      int            `json:"page"`
	PerPage   int            `json:"per_page"`
	Pages     int            `json:"pages"`
}

type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Details string `json:"details,omitempty"`
	} `json:"error"`
}

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: strings.TrimSuffix(baseURL, "/"),
		HTTP:    http.DefaultClient,
	}
}

func (c *Client) GetRecord(id int) (*StatusRecord, error) {
	rawID := strconv.Itoa(id)
	resp, err := c.request("GET", "/api/"+rawID, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("status record not found: %d", id)
	}
	if resp.StatusCode != http.StatusOK {
		err := c.parseError(resp)
		return nil, err
	}

	var record StatusRecord
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &record, nil
}

func (c *Client) ListRecords(page, perPage int, status, phase string) (*ListResponse, error) {
	query := url.Values{}
	if page > 0 {
		query.Set("page", strconv.Itoa(page))
	}
	if perPage > 0 {
		query.Set("per_page", strconv.Itoa(perPage))
	}
	if status != "" {
		query.Set("status", status)
	}
	if phase != "" {
		query.Set("phase", phase)
	}

	endpoint := "/api"
	if query.Encode() != "" {
		endpoint += "?" + query.Encode()
	}

	resp, err := c.request("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err := c.parseError(resp)
		return nil, err
	}

	var response ListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &response, nil
}

func (c *Client) CreateRecord(record StatusRecordCreate) (*StatusRecord, error) {
	body, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	httpResp, err := c.request("POST", "/api", body)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode != http.StatusCreated {
		err := c.parseError(httpResp)
		return nil, err
	}

	var created StatusRecord
	if err := json.NewDecoder(httpResp.Body).Decode(&created); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &created, nil
}

func (c *Client) UpdateRecord(id int, record StatusRecordUpdate) (*StatusRecord, error) {
	rawID := strconv.Itoa(id)
	body, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	httpResp, err := c.request("PATCH", "/api/"+rawID, body)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode != http.StatusOK {
		err := c.parseError(httpResp)
		return nil, err
	}

	var updated StatusRecord
	if err := json.NewDecoder(httpResp.Body).Decode(&updated); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &updated, nil
}

func (c *Client) DeleteRecord(id int) error {
	rawID := strconv.Itoa(id)
	httpResp, err := c.request("DELETE", "/api/"+rawID, nil)
	if err != nil {
		return err
	}

	if httpResp.StatusCode != http.StatusNoContent && httpResp.StatusCode != http.StatusOK {
		err := c.parseError(httpResp)
		return err
	}
	return nil
}

func (c *Client) request(method, endpoint string, body []byte) (*http.Response, error) {
	reqURL := c.BaseURL + endpoint
	req, err := http.NewRequest(method, reqURL, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

func (c *Client) parseError(resp *http.Response) error {
	var errorResp ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server error (%d): %s", resp.StatusCode, string(body))
	}
	return fmt.Errorf("server error (%d): %s", errorResp.Error.Code, errorResp.Error.Message)
}

func ValidateURL(urlStr string) error {
	u, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	if u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("url must include scheme (http:// or https://) and host")
	}
	return nil
}
