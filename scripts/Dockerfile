FROM --platform=linux/amd64 golang:1.22.5-alpine3.20 as builder

WORKDIR /app

COPY ./scripts/ ./

RUN go mod download
RUN go build -o dotfile /app/main.go

FROM --platform=linux/amd64 debian

USER root
WORKDIR /root/

RUN apt update && apt install -y curl git

RUN mkdir -p /root/projects/

COPY --from=builder /app/dotfile /root/dotfile
COPY ./ /root/projects/dotfiles/

CMD ["/root/dotfile", "standard"]
