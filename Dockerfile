FROM golang:alpine
RUN apk add git

ENV MODE=production
ENV GIN_MODE=release
ENV GOPRIVATE=github.com/tribalwarshelp

COPY ./.netrc /root/.netrc
RUN chmod 600 /root/.netrc

WORKDIR /go/src/app
COPY . .

RUN go build -o main .

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.2.1/wait ./wait
RUN chmod +x ./wait

CMD ./wait && ./main