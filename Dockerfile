FROM golang:1.22.3-bullseye

LABEL codechat.worker.version="0.0.1" 
LABEL codechat.worker.description="Worker HTTP para Gerenciamento de Sessões na API de WhatsApp" 
LABEL codechat.worker.maintainer="jrCleber" 
LABEL codechat.worker.git="https://github.com/jrCleber"
LABEL codechat.worker.contact="suporte@codechat.dev"

ENV DOCKER_ENV=true

WORKDIR /worker

COPY go.mod .

RUN go mod tidy

COPY ./main .

RUN mkdir -p ./instances

CMD [ "./main" ]