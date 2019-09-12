FROM golang:latest

LABEL maintainer="Luis Fernando <128423@upf.br>"


COPY . $GOPATH/src/api_ceva/

ENV GIN_MODE=release

WORKDIR $GOPATH/src/api_ceva/
COPY Gopkg.lock Gopkg.toml $GOPATH/src/api_ceva/
RUN apk add --no-cache bash git openssh  && dep ensure --vendor-only && apk del bash git openssh

RUN go install

EXPOSE 8080

CMD ["api-ceva"]




