package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/proxy"
)

func proxyReq(c *gin.Context) {
	remoteURL, err := url.Parse(os.Getenv("REMOTE"))
	if err != nil {
		panic(err)
	}

	// Create a SOCKS5 proxy dialer
	dialer, err := proxy.SOCKS5("tcp", os.Getenv("PROXY"), nil, proxy.Direct)
	if err != nil {
		panic(err)
	}

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(remoteURL)
	proxy.Transport = &http.Transport{
		Dial: dialer.Dial,
	}

	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remoteURL.Host
		req.URL.Scheme = remoteURL.Scheme
		req.URL.Host = remoteURL.Host
		req.URL.Path = c.Param("proxyPath")
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
func main() {

	r := gin.Default()

	r.Any("/*proxyPath", proxyReq)

	fmt.Println("Running on localhost:8080")
	r.Run(":" + os.Getenv("PORT"))
}
