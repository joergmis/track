// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

const (
	Api_key_authScopes  = "api_key_auth.Scopes"
	Api_user_authScopes = "api_user_auth.Scopes"
)

// Customer defines model for Customer.
type Customer struct {
	Active          bool   `json:"active"`
	BillableDefault bool   `json:"billable_default"`
	Color           int    `json:"color"`
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Note            string `json:"note"`
}

// Customers defines model for Customers.
type Customers struct {
	Customers []Customer `json:"customers"`
	Paging    Pagination `json:"paging"`
}

// Pagination defines model for Pagination.
type Pagination struct {
	CountItems   int `json:"count_items"`
	CountPages   int `json:"count_pages"`
	CurrentPage  int `json:"current_page"`
	ItemsPerPage int `json:"items_per_page"`
}

// Project defines model for Project.
type Project struct {
	Active          bool   `json:"active"`
	BillableDefault bool   `json:"billable_default"`
	Completed       bool   `json:"completed"`
	CustomersId     int    `json:"customers_id"`
	Deadline        string `json:"deadline"`
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Number          string `json:"number"`
}

// Projects defines model for Projects.
type Projects struct {
	Paging   Pagination `json:"paging"`
	Projects []Project  `json:"projects"`
}

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetV2Customers request
	GetV2Customers(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetV2Projects request
	GetV2Projects(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetV2Customers(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetV2CustomersRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetV2Projects(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetV2ProjectsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetV2CustomersRequest generates requests for GetV2Customers
func NewGetV2CustomersRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v2/customers")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetV2ProjectsRequest generates requests for GetV2Projects
func NewGetV2ProjectsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v2/projects")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetV2CustomersWithResponse request
	GetV2CustomersWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetV2CustomersResponse, error)

	// GetV2ProjectsWithResponse request
	GetV2ProjectsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetV2ProjectsResponse, error)
}

type GetV2CustomersResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Customers
}

// Status returns HTTPResponse.Status
func (r GetV2CustomersResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetV2CustomersResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetV2ProjectsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Projects
}

// Status returns HTTPResponse.Status
func (r GetV2ProjectsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetV2ProjectsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetV2CustomersWithResponse request returning *GetV2CustomersResponse
func (c *ClientWithResponses) GetV2CustomersWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetV2CustomersResponse, error) {
	rsp, err := c.GetV2Customers(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetV2CustomersResponse(rsp)
}

// GetV2ProjectsWithResponse request returning *GetV2ProjectsResponse
func (c *ClientWithResponses) GetV2ProjectsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetV2ProjectsResponse, error) {
	rsp, err := c.GetV2Projects(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetV2ProjectsResponse(rsp)
}

// ParseGetV2CustomersResponse parses an HTTP response from a GetV2CustomersWithResponse call
func ParseGetV2CustomersResponse(rsp *http.Response) (*GetV2CustomersResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetV2CustomersResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Customers
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseGetV2ProjectsResponse parses an HTTP response from a GetV2ProjectsWithResponse call
func ParseGetV2ProjectsResponse(rsp *http.Response) (*GetV2ProjectsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetV2ProjectsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Projects
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List all customers
	// (GET /v2/customers)
	GetV2Customers(ctx echo.Context) error
	// List all projects
	// (GET /v2/projects)
	GetV2Projects(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetV2Customers converts echo context to params.
func (w *ServerInterfaceWrapper) GetV2Customers(ctx echo.Context) error {
	var err error

	ctx.Set(Api_user_authScopes, []string{})

	ctx.Set(Api_key_authScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetV2Customers(ctx)
	return err
}

// GetV2Projects converts echo context to params.
func (w *ServerInterfaceWrapper) GetV2Projects(ctx echo.Context) error {
	var err error

	ctx.Set(Api_user_authScopes, []string{})

	ctx.Set(Api_key_authScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetV2Projects(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/v2/customers", wrapper.GetV2Customers)
	router.GET(baseURL+"/v2/projects", wrapper.GetV2Projects)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7RVTY/jNgz9K6raoxtnpzfftlOgWLSHBRYtCgRBoMiMrYksqRSVrRH4vxeSP+IkznYG",
	"xdxkU3p8JN+TzlzaxlkDhjwvztzLGhqRls/Bk20A49qhdYCkIEWEJHWCuKLWAS/43loNwvAu43ultdhr",
	"2JVwEEHT8i5ptcVZSBmCCjCGVLn834hmntITKlOlgKWlQJdxhL+DQih5sYmwA0Y28l8gO6CNBLfZiGr3",
	"LyApphvb4u/7IuchRdCkxQ8IB17w7/NLp/OhzfnU427KJBBFG7+dqGIh/4HwOe4SpKy5K3kAyGa8lgqa",
	"IdxXZIOh3VTK/VT6DU5U8GhDQIRhy4OBR/SdA3y453aU1wduclxzyq5KWCwfbVq+i8obp4GgfBAex7J7",
	"pPkSRKmVWdb9240Smj3gK61yRW5yzgDxbQtNrOcd+EbrF5z0du1nEWFCe5X5xsnfee+Rj6YM98V0Gfcg",
	"Aypqv0T4QUJO7Y7Q7kSgOtEyvOA1iDL1sB8U/+vHZ23l0Zb2o1O/QcsvdPrvLktAwQO+AekPn0I3UJEo",
	"/EOARuhfrEwsS/ASletvAP5cgzwyG4hRDczbgBKYtGUcZ0Ad0xI5X+R5pagO+5W0Tf5iAatG+ZxQyGNS",
	"pznYBfCBH/MOZMYQ4n4o2QFtE/MpZNZE7bCPnz+x0srQgKE04hX7Gb4KhIxRrTxrVFUTOwG27CtozZSR",
	"OpTAANGiZxaZV43TLdtDjA06/C52RJGGOZeY6vTEM34C9D3P9Wq9+hDLsA6McIoX/KfVh9U6ikBQndqW",
	"n57yqxu/gnQPRB0nxp9KXvBfgf58urwZUVneWeN7gTyt1/01awgM9ZJxWsl0Pn/x/Z3cC/a1b4nvp3zd",
	"+C9BSvD+EDSb+PWqDU0jsOUF/115YkJrJmdsSVQ+emD2hMRTsfa53x6XPnn8HSufcvzPwt2F61j3xfLd",
	"3OO82JxvTbnZdtmt5TfbbhtP4SlpZHO+sVDTruQgw+SjKLV4Ykh/Hp19aX/MMfycuM3+xVRKgufdtvs3",
	"AAD//+7WNTbcCQAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
