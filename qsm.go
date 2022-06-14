package goqsm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang/glog"
)

// Client .
type Client struct {
	apiKey     string
	baseURL    string
	HTTPClient *http.Client
}

type AuthClient struct {
	Client
	accessToken  string
	refreshToken string
}

type AuthRes struct {
	AccessToken  string `json:"accessToken"`
	ExpireTime   int    `json:"expireTime"`
	RefreshToken string `json:"refreshToken"`
}

type EmptyData []interface{}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// NewClient creates new Facest.io client with given API key
func NewClient(ip string) *Client {
	return &Client{
		HTTPClient: &http.Client{},
		baseURL:    "http://" + ip,
	}
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
// func (c *Client) NewRequest(ctx context.Context, method, urlPath string, body interface{}) (*http.Request, error) {
func (c *Client) NewRequest(ctx context.Context, method, urlPath string, body url.Values) (*http.Request, error) {
	var (
		req *http.Request
		err error
	)

	urlStr := c.baseURL + urlPath
	glog.V(2).Infof("[NewRequest] %s url: %s\n", method, urlStr)
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	if body != nil {
		glog.V(3).Infof("[NewRequest] body: %v\n", body)
		req, err = http.NewRequest(method, u.String(), strings.NewReader(body.Encode()))
	} else {
		req, err = http.NewRequest(method, u.String(), nil)
	}

	if err != nil {
		return nil, err
	}

	// req.Header.Add("Content-Type", mediaType)
	// req.Header.Add("Accept", mediaType)
	// req.Header.Add("User-Agent", c.UserAgent)
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Set("Content-Type", "application/json; charset=utf-8")
	return req, nil
}

// Content-type and body should be already added to req
func (c *Client) sendRequest(ctx context.Context, req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if c.apiKey != "" {
		glog.V(4).Infof("[sendRequest] apiKey: %s\n", c.apiKey)
		req.Header.Set("Authorization", c.apiKey)
	}

	req = req.WithContext(ctx)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		glog.Errorf("[sendRequest] err: %v\n", err)
		return err
	}

	defer res.Body.Close()

	glog.V(2).Infof("[sendRequest] StatusCode: %d\n", res.StatusCode)
	if res.StatusCode != http.StatusOK {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(v); err != nil {
		return err
	}

	glog.V(4).Infof("[sendRequest] res: %+v\n", v)
	return nil
}

func (c *Client) login(ctx context.Context, user string, passwd string) (*AuthRes, error) {
	params := url.Values{}
	params.Add("user", user)
	params.Add("password", passwd)
	params.Add("offlineAccess", "true")

	req, err := c.NewRequest(ctx, http.MethodPost, "/auth/get", params)
	if err != nil {
		return nil, err
	}

	res := AuthRes{}
	if err := c.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetAuthClient(ctx context.Context, user string, passwd string) (*AuthClient, error) {
	res, err := c.login(ctx, user, passwd)
	if err != nil {
		return nil, fmt.Errorf("login failed: %v\n", err)
	}

	glog.V(3).Infof("AccessToken: %s\n", res.AccessToken)

	ret := &AuthClient{}
	ret.accessToken = res.AccessToken
	ret.refreshToken = res.RefreshToken
	ret.baseURL = c.baseURL
	ret.HTTPClient = c.HTTPClient
	ret.apiKey = res.AccessToken

	return ret, nil
}
