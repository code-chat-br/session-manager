FROM golang:1.22.3-bullseye

LABEL version="0.0.1" description="Worker HTTP para Gerenciamento de Sessões na API de WhatsApp" 
LABEL maintainer="jrCleber" git="https://github.com/jrCleber"
LABEL contact="suporte@codechat.dev"

ENV DOCKER_ENV=true

WORKDIR /worker

COPY go.mod .

RUN go mod tidy

COPY ./main .

RUN mkdir -p ./instances

CMD [ "./main" ]