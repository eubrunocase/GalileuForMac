package guardian

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

func StartGuardian() {
	ca, err := tls.LoadX509KeyPair("ca.pem", "key.pem")
	if err != nil {
		log.Fatal("[ERRO] Falha ao carregar certificados. Verifique se ca.pem e key.pem estão na raiz: ", err)
	}

	goproxy.GoproxyCa = ca

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true 

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	analyzer := NewAnalyzer()

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			if r.Body != nil {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					return r, nil
				}

				found, cleanPayload := analyzer.Analyze(body)

				if found {
					fmt.Printf("[GALILEU] INTERCEPTADO: Dados sensíveis removidos do destino: %s\n", r.Host)
					r.Body = io.NopCloser(bytes.NewBuffer(cleanPayload))
					r.ContentLength = int64(len(cleanPayload))
				} else {
					r.Body = io.NopCloser(bytes.NewBuffer(body))
				}
			}
			return r, nil
		})

	fmt.Println("[GALILEU] Proxy ativo e escutando em http://127.0.0.1:9000")
	log.Fatal(http.ListenAndServe(":9000", proxy))
}