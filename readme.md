# Galileu — Proxy de Segurança e Governança para LLMs
### Versão macOS (Apple Silicon & Intel)
 
> **Galileu** é uma ferramenta de segurança e governança de dados voltada para o monitoramento e sanitização de informações enviadas a provedores de Inteligência Artificial (LLMs). O projeto adota uma arquitetura de **Proxy Reverso MITM (Man-in-the-Middle)**, actuando como camada inteligente entre a sua ferramenta de desenvolvimento e os servidores das LLMs.
 
---
 
## Arquitectura do Sistema
 
```
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│   Cliente   │───▶  │  Galileu    │───▶  │   LLM       │
│  (OpenCode) │◀───  │  Proxy MITM │◀───  │  Provider   │
└─────────────┘      └─────────────┘      └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │  Analyzer   │
                    │ (Sanitização)│
                    └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │   Audit     │
                    │    Log      │
                    └─────────────┘
```
 
---
 
## Pré-requisitos
 
| Requisito | Detalhe |
|---|---|
| **Sistema Operativo** | macOS — compatível com Apple Silicon (M1/M2/M3) e arquitetura Intel |
| **Go** | Versão 1.23 ou superior (necessário apenas para compilação) |
| **Certificados** | Ficheiros `ca.pem` e `key.pem` gerados e colocados na raiz do projeto |
| **Privilégios** | Não requer `sudo` para execução na porta 9000 |
 
---
 
## Compilação
 
Abra o Terminal na raiz do projeto e execute o comando adequado à sua arquitectura:
 
**Apple Silicon (M1/M2/M3 — ARM64):**
```bash
GOOS=darwin GOARCH=arm64 go build -o galileu ./cmd/sentinel/main.go
```
 
**Intel (AMD64):**
```bash
GOOS=darwin GOARCH=amd64 go build -o galileu ./cmd/sentinel/main.go
```
 
---
 
## Estrutura de Ficheiros
 
```
Galileu/
├── galileu              # Executável principal (macOS)
├── ca.pem               # Certificado CA público (exportado do Acesso às Chaves)
├── key.pem              # Chave privada do CA (⚠️ NÃO submeter para o repositório)
├── start-opencode.sh    # Script shell para iniciar o OpenCode com proxy
└── galileu_audit.log    # Registo de auditoria (gerado automaticamente)
```
 
---
 
## Como Utilizar
 
### Passo 1 — Configurar o Certificado CA no macOS
 
Para que o proxy MITM funcione correctamente com HTTPS **sem emitir avisos de segurança**, o certificado `ca.pem` deve ser adicionado ao **Acesso às Chaves (Keychain Access)** do macOS e definido como confiável:
 
1. Abra o **Acesso às Chaves** (Keychain Access).
2. Arraste o ficheiro `ca.pem` para a lista de certificados, ou use **Ficheiro > Importar Itens**.
3. Clique duas vezes no certificado importado.
4. Expanda a secção **Confiar** (Trust).
5. Defina **Ao usar este certificado** como **Confiar Sempre** (Always Trust).
6. Feche a janela e introduza a sua palavra-passe de sistema quando solicitado.
### Passo 2 — Executar o Galileu
 
Abra o Terminal na raiz do projecto e execute:
 
```bash
./galileu
```
 
O programa irá:
 
- Carregar os certificados locais (`ca.pem` e `key.pem`).
- Iniciar o proxy na porta **9000**.
- Activar o registo (logging) de auditoria.
> Não são necessários privilégios `sudo` para a porta 9000.
 
### Passo 3 — Configurar o OpenCode
 
Num **novo Terminal**, dê permissão de execução ao script e execute-o:
 
```bash
chmod +x start-opencode.sh
./start-opencode.sh
```
 
Ou configure manualmente as variáveis de ambiente na sua sessão:
 
```bash
export HTTP_PROXY="http://127.0.0.1:9000"
export HTTPS_PROXY="http://127.0.0.1:9000"
export NODE_TLS_REJECT_UNAUTHORIZED=0
opencode
```
 
> **Nota:** Se o comando `opencode` não funcionar nativamente, substitua-o por `open -a "Visual Studio Code"` no seu script.
 
### Passo 4 — Utilizar o OpenCode normalmente
 
A partir deste momento, **todas as requisições do OpenCode** para os provedores de IA passarão pelo proxy Galileu, que irá:
 
- Detectar e remover dados sensíveis automaticamente.
- Registar cada requisição para auditoria.
---
 
## Hosts Monitorizados
 
O Galileu intercepta requisições para os seguintes provedores:
 
| Provedor | Host |
|---|---|
| OpenCode | `opencode.ai` |
| OpenAI | `api.openai.com` |
| Anthropic | `api.anthropic.com` |
| Google AI | `generativelanguage.googleapis.com` |
| Cohere | `api.cohere.ai` |
| Mistral | `api.mistral.ai` |
 
---
 
## Detecção de Dados Sensíveis
 
O **Analyzer** detecta e sanitiza automaticamente os seguintes padrões:
 
| Tipo | Padrão | Exemplo |
|---|---|---|
| OpenAI API Key | `sk-...` | `sk-1234567890abcdef...` |
| OpenAI Project Key | `sk-proj-...` | `sk-proj-abc123...` |
| Anthropic API Key | `sk-ant-...` | `sk-ant-abc123...` |
| Google API Key | `AIzaSy...` | `AIzaSyABC123...` |
| GitHub Token | `ghp_...` | `ghp_abcdef123456...` |
| Slack / Discord | `xox[baprs]-...` | `xoxb-123456...` |
| AWS Access Key | `AKIA...` | `AKIAIOSFODNN7...` |
 
Todos os dados sensíveis detectados são substituídos por `[REDACTED_BY_GALILEU]`.
 
---
 
## Registos de Auditoria
 
O ficheiro `galileu_audit.log` contém um registo JSON de cada requisição interceptada:
 
```json
{"timestamp":"2026-04-29T10:00:00Z","host":"opencode.ai","path":"/v1/chat/completions","method":"POST","redacted":true,"pattern_type":"sensitive_data"}
{"timestamp":"2026-04-29T10:05:00Z","host":"api.openai.com","path":"/v1/chat/completions","method":"POST","redacted":false,"pattern_type":""}
```
 
---
 
## Comprovação de Testes
 
Os testes foram realizados com sucesso. As imagens de comprovação encontram-se na pasta `img/`:
 
| # | Teste | Descrição |
|---|---|---|
| 1 | **Ficheiro `.env` com Dados Sensíveis** | Ambiente simulado com múltiplas chaves de API |
| 2 | **Terminal do Galileu** | Registo de execução do proxy com interceptação das requisições |
| 3 | **Resposta do OpenCode** | O assistente recebe o payload com as chaves substituídas pela etiqueta de segurança |
 
### Resultado dos Testes
 
| Dados Enviados | Dados Redatados | Estado |
|---|---|---|
| `sk-...` (OpenAI) | ✅ `[REDACTED_BY_GALILEU]` | Detectado |
| `sk-proj-...` (OpenAI Project) | ✅ `[REDACTED_BY_GALILEU]` | Detectado |
| `sk-ant-...` (Anthropic) | ✅ `[REDACTED_BY_GALILEU]` | Detectado |
| `AIzaSy...` (Google) | ✅ `[REDACTED_BY_GALILEU]` | Detectado |
| `ghp_...` (GitHub) | ✅ `[REDACTED_BY_GALILEU]` | Detectado |
| `xoxb-...` (Slack/Discord) | ✅ `[REDACTED_BY_GALILEU]` | Detectado |
| `AKIA...` (AWS Key) | ✅ `[REDACTED_BY_GALILEU]` | Detectado |
 
---
 
## Performance
 
O Galileu foi optimizado para ambientes de desenvolvimento:
 
- **Goroutines** — Cada requisição é processada na sua própria goroutine.
- **Buffer Pooling** — Reutilização de memória com `sync.Pool`.
- **Regex Pré-compilado** — Padrões de detecção compilados na inicialização.
- **Processamento Assíncrono** — Logging não bloqueante.
---
 
## Resolução de Problemas
 
### "Falha ao carregar certificados"
 
Verifique se exportou correctamente a chave e o certificado do **Acesso às Chaves** para os ficheiros `ca.pem` e `key.pem`, e confirme que ambos se encontram na raiz do directório onde está a executar o comando.
 
### OpenCode não conecta ao proxy
 
Confirme que as variáveis de ambiente foram exportadas correctamente na sessão actual:
 
```bash
echo $HTTP_PROXY
```
 
O resultado deve ser: `http://127.0.0.1:9000`
 
### Erros de certificado SSL/TLS no cliente
 
O certificado CA (`ca.pem`) deve constar no **Acesso às Chaves (Keychain Access)** do macOS com a confiança definida como **Confiar Sempre** (Always Trust). Consulte o Passo 1 para instruções detalhadas.
 
---
 
## Arquitectura do Código
 
```
cmd/sentinel/main.go      # Ponto de entrada do sistema
internal/guardian/
  ├── guardian.go         # Configuração do proxy MITM e carregamento de PEM
  ├── analyzer.go         # Detecção e sanitização de dados sensíveis
  └── audit.go            # Sistema de registo de auditoria
```
 
---
 
## Segurança
 
- A chave privada (`key.pem`) deve ser mantida em estrita segurança na máquina local.
- **Nunca** efectue commit dos ficheiros `.pem` para o repositório — confirme que o `.gitignore` está actualizado.
- O proxy actua exclusivamente sobre as ferramentas que configurarem explicitamente a porta **9000**.
---
 
## Licença
 
Este projecto é para fins educacionais e de segurança interna.  
Todos os direitos são reservados ao programador **Bruno Dantas de Oliveira Cazé** — [github.com/eubrunocase/Galileu](https://github.com/eubrunocase/Galileu)