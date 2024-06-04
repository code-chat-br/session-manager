# Implementação de um Worker HTTP para Gerenciamento de Sessões na API de WhatsApp

## Worker HTTP para a API do WhatsApp

Este worker HTTP é um sistema que gerencia os arquivos de conexão das sessões do WhatsApp. Ele escuta requisições HTTP na porta **5656**, recebe e armazena arquivos de sessão, e também recupera esses arquivos quando solicitado. A segurança é garantida mantendo o worker na mesma rede que a API de WhatsApp, e a porta do worker não é exposta ao público. O worker pode ser executado em um ambiente Docker e está pronto para futuras expansões e refinamentos.


## Funcionalidades

1. **Escuta de Requisições HTTP:** O worker foi configurado para ouvir requisições HTTP na porta 5656. Qualquer requisição enviada para essa porta será tratada pelo worker.

2. **Recebimento de Arquivos de Sessão:** Quando uma sessão de WhatsApp é iniciada ou atualizada, o worker recebe os arquivos de conexão correspondentes através de requisições HTTP POST. Esses arquivos contêm informações necessárias para manter a conexão ativa e permitir a comunicação contínua com o WhatsApp.

3. **Armazenamento de Arquivos de Sessão:** Após receber os arquivos de sessão, o worker os armazena em um local seguro no servidor. O armazenamento é feito de maneira organizada para garantir que os arquivos possam ser facilmente recuperados e identificados.

4. **Recuperação de Arquivos de Sessão:** O worker também oferece uma funcionalidade de recuperação dos arquivos de sessão. Quando solicitado através de uma requisição HTTP GET, ele localiza e retorna os arquivos de sessão específicos, permitindo que a conexão com o WhatsApp seja restaurada ou mantida.

5. **Segurança e Confiabilidade:** A segurança é uma responsabilidade compartilhada. É essencial que o worker permaneça na mesma rede que a API para evitar acessos não autorizados. Apenas a API deve se comunicar diretamente com o worker, e a porta na qual o worker escuta as requisições não deve ser exposta ao público.</br>
Além disso, um token global deve ser definido nas variáveis de ambiente. O cliente deve enviar esse token em “headers.apikey” ao fazer uma requisição. Esse token global deve ser o mesmo configurado na API cliente para gestão das instâncias. Isso adiciona uma camada extra de segurança, garantindo que apenas componentes autorizados dentro da mesma rede possam acessar as funcionalidades do worker, minimizando riscos de segurança.

6. **Requisitos de Rede:** É recomendado, mas não obrigatório, que este worker permaneça na mesma rede que a aplicação principal para assegurar a comunicação eficiente e segura entre os componentes. Se estiver utilizando o Docker Swarm, o worker deve estar na mesma rede do Swarm para garantir o correto funcionamento e a integração dos serviços.

## Estrutura do projeto
```
/worker-session
├── /.vscode
│   └── settings.json
│
├── /api
│   ├── /handlers
│   │   ├── handler.go
│   │   └── sessio_handler.go 
│   ├── /middlewares
│   │   ├── auth_guard.go
│   │   └── body_size.go
│   └── /routers
│       ├── routers.go
│       └── session_routes.go
│
├── /cmd
│   └── /server
│       └── main.go
│
├── /instances # criasda internamente pelo serviço
│
├── /internal
│   ├─ /session
│   │  └── service.go
│   └── app.go
│
├── /pkg
│   ├── /config
│   │   └── config.go
│   └── /db
│       ├── /models
│       │   └── creds_model.go
│       └── conn.go
│
├── /tmp
│   ├── build-errors.log
│   └── main
│
├── .air.toml
├── .env
├── .gitignore
├── build.sh
├── docker-compose.yaml
├── Dockerfile
├── go.sum
├── LICENSE
├── main
└── README.md
```

---
## Ajuste de compatibilidade com SQLite se necessário

- `build-essential`
  ```sh
  sudo apt-get install build-essential
  ```

- `gcc`
  ```sh
  sudo apt-get install gcc
  ```
---

## Instalação

### Instalando o GO

1. **Download do binário**
```sh
wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz
```

2. **Extraindo binário**
```sh
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz
```

3. **Definindo o `PATH`**
```sh
export PATH=$PATH:/usr/local/go/bin
```

4. **Testando o `go`**
```sh
go version

# go version go1.22.3 linux/amd64
```

### Clonando o Repositório

```sh
git clone https://github.com/code-chat-br/session-manager.git
```

1. **Instalando as dependências**
```sh
go mod tidy
```

3. **Executando a aplicação**
```sh
go run cmd/server/main.go
```

### Variáveis de ambiente

> **[.env](./.env_dev)**

```sh
cp .env_dev .env
```

- `GLOBAL_AUTH_TOKEN` **(requerido)**: Token global definido na API cliente;
- `BODY_SIZE` **(padrão 5mb)**: comprimento máximo do corpo da requisição;
- `HTTP_LOGS` **(padrão `false`)**: Determina a exibição dos logs http.

### Build da aplicação

```sh
sh build.sh
```

Executando o build
```
./main
```

# Discurssões

As discurções sobre esse worker devem ser realizadas [aqui](https://github.com/code-chat-br/whatsapp-api/discussions/131).

- [Dockerfile](./Dockerfile)
- [docker-compose](./docker-compose.yaml)
- [codechat/session-manager](https://hub.docker.com/r/codechat/session-manager)
