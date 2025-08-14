package httpclient

import (
	"net/http"
	"time"
)

// NWS requiere un User-Agent descriptivo con contacto
func New(userAgent string) *http.Client {
	tr := &http.Transport{Proxy: http.ProxyFromEnvironment, MaxIdleConns: 100, IdleConnTimeout: 90 * time.Second}
	return &http.Client{Transport: roundTripper{base: tr, ua: userAgent}, Timeout: 12 * time.Second}
}

type roundTripper struct {
	base http.RoundTripper
	ua   string
}

func (rt roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := *req
	req2.Header = req.Header.Clone()
	if rt.ua != "" {
		req2.Header.Set("User-Agent", rt.ua)
	}
	req2.Header.Set("Accept", "application/ld+json")
	return rt.base.RoundTrip(&req2)
}
