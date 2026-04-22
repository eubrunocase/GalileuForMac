package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"Galileu/internal/guardian"
)

func main() {
	// 1. Configurar o redirecionamento ao iniciar
	if err := setupFirewall(); err != nil {
		fmt.Printf("Erro ao configurar firewall (requer sudo): %v\n", err)
		os.Exit(1)
	}

	// 2. Canal para capturar sinal de interrupção (Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\n[GALILEU] Encerrando e normalizando tráfego...")
		cleanupFirewall()
		os.Exit(0)
	}()

	// 3. Iniciar o Guardian na porta 9000
	guardian.StartGuardian()
}

func setupFirewall() error {
	fmt.Println("[GALILEU] Configurando regras de redirecionamento transparente...")
	
	// Comando para ativar PF e carregar regra específica
	// Nota: No mundo real, você pode escrever um arquivo temporário com as regras e carregá-lo
	rule := `rdr pass on lo0 proto tcp from 127.0.0.1 to any port 443 -> 127.0.0.1 port 9000`
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo \"%s\" | sudo pfctl -ef -", rule))
	return cmd.Run()
}

func cleanupFirewall() {
	// Desativa o PF ou limpa as regras
	exec.Command("sudo", "pfctl", "-d").Run()
	fmt.Println("[GALILEU] Tráfego normalizado.")
}