package http

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

type HTTPOptions struct {
	URI     string
	Headers map[string]string

	Payload []byte
}

var (
	applicationJsonContentType = []byte("application/json")

	readTimeout, _         = time.ParseDuration("500ms")
	writeTimeout, _        = time.ParseDuration("500ms")
	maxIdleConnDuration, _ = time.ParseDuration("1h")

	// Fasthttp client configurations.
	fasthttpClient *fasthttp.Client = &fasthttp.Client{
		ReadTimeout:         readTimeout,
		WriteTimeout:        writeTimeout,
		MaxIdleConnDuration: maxIdleConnDuration,

		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,

		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}
)

// Get method (with fasthttp)
func Get(opts *HTTPOptions) (acquireResponse *fasthttp.Response, err error) {
	acquireRequest := fasthttp.AcquireRequest()
	acquireRequest.SetRequestURI(opts.URI)
	acquireRequest.Header.SetMethod(fasthttp.MethodGet)

	if len(opts.Headers) > 0 {
		for key, value := range opts.Headers {
			acquireRequest.Header.Set(key, value)
		}
	}

	acquireResponse = fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(acquireResponse)

	if err = func(request *fasthttp.Request, response *fasthttp.Response) error {
		defer fasthttp.ReleaseRequest(request)

		if err = fasthttpClient.Do(acquireRequest, acquireResponse); err != nil {
			log.Error().
				Err(err).
				Str("request_uri", opts.URI).
				Msg(ErrorDoRequest.Error())

			return ErrorDoRequest
		}

		return nil
	}(acquireRequest, acquireResponse); err != nil {
		return nil, err
	}

	return acquireResponse, nil
}

// Post method (with fasthttp)
func Post(opts *HTTPOptions) (acquireResponse *fasthttp.Response, err error) {
	acquireRequest := fasthttp.AcquireRequest()
	acquireRequest.SetRequestURI(opts.URI)
	acquireRequest.Header.SetMethod(fasthttp.MethodPost)
	acquireRequest.Header.SetContentTypeBytes(applicationJsonContentType)

	if opts.Payload != nil {
		acquireRequest.SetBodyRaw(opts.Payload)
	}

	if opts.Headers != nil {
		for key, value := range opts.Headers {
			acquireRequest.Header.Set(key, value)
		}
	}

	acquireResponse = fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(acquireResponse)

	if err = func(request *fasthttp.Request, response *fasthttp.Response) error {
		defer fasthttp.ReleaseRequest(request)

		if err = fasthttpClient.Do(acquireRequest, acquireResponse); err != nil {
			log.Error().
				Err(err).
				Str("request_uri", opts.URI).
				Msg(ErrorDoRequest.Error())

			return ErrorDoRequest
		}

		return nil
	}(acquireRequest, acquireResponse); err != nil {
		return nil, err
	}

	return acquireResponse, nil
}
