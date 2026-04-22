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
	// 1. Carregar o par de chaves do certificado para interceptação MITM
	ca, err := tls.LoadX509KeyPair("ca.pem.cer", "key.pem")
	if err != nil {
		log.Fatal("Erro ao carregar certificados. Verifique os nomes dos arquivos:", err)
	}

	proxy := goproxy.NewProxyHttpServer()
	
	// Configura o proxy para usar o certificado para interceptação MITM
	goproxy.GoproxyCa = ca
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
					fmt.Println("[GALILEU] Interceptado: Chaves de API detectadas e mascaradas.")
					r.Body = io.NopCloser(bytes.NewBuffer(cleanPayload))
					r.ContentLength = int64(len(cleanPayload))
				} else {
					r.Body = io.NopCloser(bytes.NewBuffer(body))
				}
			}
			return r, nil
		})

	fmt.Println("[GALILEU] Ativo no macOS. Ouvindo na porta :9000...")
	log.Fatal(http.ListenAndServe(":9000", proxy))
}