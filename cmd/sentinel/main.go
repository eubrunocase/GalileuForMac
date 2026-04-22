package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"Galileu/internal/guardian" // Verifique se o path do import está correto conforme seu go.mod
)

func main() {
	// 1. Configura o redirecionamento apenas para tráfego local (127.0.0.1)
	if err := setupSafeFirewall(); err != nil {
		fmt.Printf("[ERRO] Falha ao configurar firewall: %v\n", err)
		fmt.Println("Certifique-se de rodar com sudo ou ter permissões no pfctl.")
		os.Exit(1)
	}

	// 2. Captura sinais de encerramento para normalizar a rede
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n[GALILEU] Encerrando e normalizando tráfego...")
		cleanupFirewall()
		os.Exit(0)
	}()

	fmt.Println("[GALILEU] Sistema de Proteção Ativo.")
	fmt.Println("[INFO] Monitorando apenas chamadas de aplicações locais para LLMs.")
	
	// 3. Inicia o servidor na porta 9000
	guardian.StartGuardian()
}

func setupSafeFirewall() error {
    fmt.Println("[GALILEU] Aplicando regras de firewall...")
    
    // Regra específica para redirecionar tráfego local para a porta 9000
    rule := "rdr pass on lo0 proto tcp from 127.0.0.1 to any port 443 -> 127.0.0.1 port 9000\n"
    
    // Cria um arquivo de configuração temporário
    tmpFile, err := os.CreateTemp("", "galileu_pf.conf")
    if err != nil {
        return err
    }
    defer os.Remove(tmpFile.Name())

    if _, err := tmpFile.WriteString(rule); err != nil {
        return err
    }
    tmpFile.Close()

    // -E ativa o PF, -f carrega o arquivo de regras
    cmd := exec.Command("sudo", "pfctl", "-E", "-f", tmpFile.Name())
    return cmd.Run()
}

func cleanupFirewall() {
	// Desativa as regras e limpa a tabela do PF
	exec.Command("sudo", "pfctl", "-d").Run()
	fmt.Println("[GALILEU] Firewall normalizado. Navegação segura restaurada.")
}