// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package utils

import (
	"crypto/tls"
	"io"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	headerUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"
	client          = &fasthttp.Client{
		ReadTimeout:                   10 * time.Second,
		WriteTimeout:                  10 * time.Second,
		MaxIdleConnDuration:           1 * time.Hour,
		NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
		DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
		DisablePathNormalizing:        true,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		// increase DNS cache time to an hour instead of default minute
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}
)

type FastResponse struct {
	Body       []byte
	StatusCode int
}

type ReqSetter func(req *fasthttp.Request)

// WithHeaders sets the headers of the request.
//
// Parameters:
//   - headers: a map of string keys to string values representing the headers to be set.
//
// Returns:
//   - ReqSetter: func(req *fasthttp.Request)
func WithHeaders(headers map[string]string) ReqSetter {
	return func(req *fasthttp.Request) {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
}

// WithBodyStream sets the body of the fasthttp request to the given reader and size.
//
// Parameters:
//   - r: The reader to use as the body of the request.
//   - size: The size of the body.
//
// Returns:
//   - ReqSetter: func(req *fasthttp.Request)
func WithBodyStream(r io.Reader, size int) ReqSetter {
	return func(req *fasthttp.Request) {
		req.SetBodyStream(r, size)
	}
}

// FastGet retrieves a response from the specified URL using the GET method.
//
// Parameters:
//   - u: The URL to send the GET request to.
//   - reqSetters: Optional request setters to modify the request before sending.
//
// Returns:
//   - *FastResponse: A pointer to the FastResponse struct containing the response body and status code.
//   - error: An error object if the request fails.
func FastGet(u string, reqSetters ...ReqSetter) (*FastResponse, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(u)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetUserAgent(headerUserAgent)

	for _, setter := range reqSetters {
		setter(req)
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := client.Do(req, resp); err != nil {
		return nil, err
	}

	fastResp := &FastResponse{
		Body:       resp.Body(),
		StatusCode: resp.StatusCode(),
	}

	return fastResp, nil
}

// FastPut sends a fast HTTP PUT request to the specified URL with optional request setters.
//
// Parameters:
//   - u: The URL to send the request to.
//   - reqSetters: Optional request setters to customize the request before sending.
//
// Returns:
//   - fastResp: The response from the request containing the response body and status code.
//   - error: An error if the request fails.
func FastPut(u string, reqSetters ...ReqSetter) (*FastResponse, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(u)
	req.Header.SetMethod(fasthttp.MethodPut)
	req.Header.SetUserAgent(headerUserAgent)

	for _, setter := range reqSetters {
		setter(req)
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := client.Do(req, resp); err != nil {
		return nil, err
	}

	fastResp := &FastResponse{
		Body:       resp.Body(),
		StatusCode: resp.StatusCode(),
	}

	return fastResp, nil
}

// FastStatusCode returns fasthttp status code
func FastStatusCode(resp *FastResponse) int {
	if resp == nil {
		return 0
	}

	return resp.StatusCode
}
