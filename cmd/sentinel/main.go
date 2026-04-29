package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"Galileu/internal/guardian" 
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n[GALILEU] Encerrando o proxy MITM...")
		os.Exit(0)
	}()

	fmt.Println("[GALILEU] Iniciando Proxy MITM para monitoramento de LLMs (macOS)...")
	fmt.Println("[INFO] Certifique-se de iniciar a IDE/OpenCode apontando para http://127.0.0.1:9000")

	guardian.StartGuardian()
}