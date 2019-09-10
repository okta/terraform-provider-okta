/*
 * Copyright 2018 - Present Okta, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package okta

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/okta/okta-sdk-golang/okta/cache"
)

type (
	RequestExecutor struct {
		httpClient *http.Client
		config     *Config
		BaseUrl    *url.URL
		cache      cache.Cache
	}

	Request struct {
		body func() (io.Reader, error)
		*http.Request
	}
)

var (
	Backoff = time.Sleep

	// Limit the size of body we read in when draining the body prior to retry as it will reuse the same connection
	respReadLimit = int64(4096)
)

func NewRequestExecutor(httpClient *http.Client, cache cache.Cache, config *Config) *RequestExecutor {
	re := RequestExecutor{}
	re.httpClient = httpClient
	re.config = config
	re.cache = cache

	if httpClient == nil {
		tr := &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			IdleConnTimeout: 30 * time.Second,
		}
		re.httpClient = &http.Client{Transport: tr}
	}

	return &re
}

func (re *RequestExecutor) NewRequest(method string, url string, body interface{}) (*http.Request, error) {
	var buff io.ReadWriter
	if body != nil {
		buff = new(bytes.Buffer)
		encoder := json.NewEncoder(buff)
		encoder.SetEscapeHTML(false)
		err := encoder.Encode(body)
		if err != nil {
			return nil, err
		}
	}
	url = re.config.Okta.Client.OrgUrl + url

	req, err := http.NewRequest(method, url, buff)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "SSWS "+re.config.Okta.Client.Token)
	req.Header.Add("User-Agent", NewUserAgent(re.config).String())
	req.Header.Add("Accept", "application/json")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func NewRequest(req *http.Request) *Request {
	if req.Body != nil {
		bodyBytes, _ := ioutil.ReadAll(req.Body)
		body := func() (io.Reader, error) {
			return bytes.NewReader(bodyBytes), nil
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		return &Request{body, req}
	}
	return &Request{nil, req}
}

func (re *RequestExecutor) Do(req *http.Request, v interface{}) (*Response, error) {
	cacheKey := cache.CreateCacheKey(req)
	if req.Method != http.MethodGet {
		re.cache.Delete(cacheKey)
	}
	inCache := re.cache.Has(cacheKey)

	if !inCache {
		retryableReq := NewRequest(req)
		resp, err := re.DoWithRetries(retryableReq, 0)

		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		origResp := ioutil.NopCloser(bytes.NewBuffer(respBody))
		resp.Body = origResp

		if req.Method == http.MethodGet && reflect.TypeOf(v).Kind() != reflect.Slice {
			re.cache.Set(cacheKey, resp)
		}

		return buildResponse(resp, &v)

	}

	resp := re.cache.Get(cacheKey)
	return buildResponse(resp, &v)

}

// DoWithRetries performs a request with configured retries and backup strategy. Exposed publicly for non JSON endpoints.
func (re *RequestExecutor) DoWithRetries(req *Request, retryCount int) (*http.Response, error) {
	// Always rewind the request body when non-nil.
	if req.Body != nil {
		body, err := req.body()
		if err != nil {
			return nil, err
		}

		if c, ok := body.(io.ReadCloser); ok {
			req.Body = c
		} else {
			req.Body = ioutil.NopCloser(body)
		}
	}

	resp, err := re.httpClient.Do(req.Request)
	maxRetries := int(re.config.MaxRetries)
	bo := re.config.BackoffEnabled

	if (err != nil || isTooMany(resp)) && retryCount < maxRetries {
		if resp != nil {
			// retrying so we must drain the body
			tryDrainBody(resp.Body)
		}

		if isTooMany(resp) {
			// Using an exponential back off method with no jitter for simplicity.
			if bo {
				Backoff(backoffDuration(retryCount, re.config.MinWait, re.config.MaxWait))
			}
		}
		retryCount++

		resp, err = re.DoWithRetries(req, retryCount)
	}

	return resp, err
}

func backoffDuration(attemptNum int, min, max time.Duration) time.Duration {
	mult := math.Pow(2, float64(attemptNum)) * float64(min)
	sleep := time.Duration(mult)
	if float64(sleep) != mult || sleep > max {
		sleep = max
	}
	return sleep
}

func isTooMany(resp *http.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusTooManyRequests
}

func tryDrainBody(body io.ReadCloser) {
	defer body.Close()
	io.Copy(ioutil.Discard, io.LimitReader(body, respReadLimit))
}

type Response struct {
	*http.Response
}

func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

func CheckResponseForError(resp *http.Response) error {
	statusCode := resp.StatusCode
	if statusCode >= http.StatusOK && statusCode < http.StatusBadRequest {
		return nil
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	e := new(Error)
	json.Unmarshal(bodyBytes, &e)
	return e

}

func buildResponse(resp *http.Response, v interface{}) (*Response, error) {
	response := newResponse(resp)

	err := CheckResponseForError(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		decodeError := json.NewDecoder(resp.Body).Decode(v)
		if decodeError == io.EOF {
			decodeError = nil
		}
		if decodeError != nil {
			err = decodeError
		}

	}
	return response, err
}
