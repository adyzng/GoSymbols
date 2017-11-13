package restful

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"

	log "gopkg.in/clog.v1"
)

var (
	bufferSize4K = 4096
	bufferPool4K = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, bufferSize4K))
		},
	}

	httpHeaders = map[string]string{
		"User-Agent":  "AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.91 Safari/537.36",
		"ContentType": "application/json", //"text/html; charset=utf-8",
		"Connection":  "keep-alive",
	}
	httpClient = &http.Client{
		Timeout: time.Minute * 1,
	}
)

// HttpGet wrap of http.Get
//
func HttpGet(uri string, fnCallback func(*http.Request)) (*bytes.Buffer, error) {
	var (
		err    error
		buffer *bytes.Buffer
		resp   *http.Response
		req    *http.Request
	)
	if req, err = http.NewRequest("GET", uri, nil); err != nil {
		log.Warn("[GET] New request of %s failed: %s.", uri, err)
		return nil, err
	}
	req.Header.Set("User-Agent", httpHeaders["User-Agent"])

	if fnCallback != nil {
		fnCallback(req)
	}

	if resp, err = httpClient.Do(req); err != nil {
		log.Warn("[GET] Http GET %s failed: %s.", uri, err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
	}

	defer resp.Body.Close()
	buffer = bufferPool4K.Get().(*bytes.Buffer)
	buffer.Reset()

	if _, err2 := io.Copy(buffer, resp.Body); err2 != nil {
		log.Warn("[GET] Read response failed: %s.", err2)
		bufferPool4K.Put(buffer)
		return nil, err2
	}
	return buffer, err
}

// HttpPost wrap of http.Post
//
func HttpPost(uri string, body io.Reader, fnCallback func(*http.Request)) (*bytes.Buffer, error) {
	var (
		err    error
		buffer *bytes.Buffer
		resp   *http.Response
		req    *http.Request
	)
	if req, err = http.NewRequest("POST", uri, body); err != nil {
		log.Warn("[POST] New request of %s failed: %s.", uri, err)
		return nil, err
	}
	req.Header.Set("User-Agent", httpHeaders["User-Agent"])
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // must

	if fnCallback != nil {
		fnCallback(req)
	}

	if resp, err = httpClient.Do(req); err != nil {
		log.Warn("[POST] Http POST %s failed: %s.", uri, err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
	}

	defer resp.Body.Close()
	buffer = bufferPool4K.Get().(*bytes.Buffer)
	buffer.Reset()

	if _, err2 := io.Copy(buffer, resp.Body); err2 != nil {
		log.Warn("[POST] Read response failed: %s.", err2)
		bufferPool4K.Put(buffer)
		return nil, err2
	}
	return buffer, err
}
