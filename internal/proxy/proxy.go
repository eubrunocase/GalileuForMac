package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func StartProxy() {
	target, _ := url.Parse("https://api.openai.com")

	proxy := httputil.NewSingleHostReverseProxy(target)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		log.Printf("[SENTINEL] Interceptando chamada para: %s %s", req.Method, req.URL.Path)
		
		req.Host = target.Host
	}

	log.Println("Sentinela ativo na porta :8080...")
	log.Fatal(http.ListenAndServe(":8080", proxy))
}