# Used to release the binary
FROM articulate/articulate-golang:1.12

ENV SERVICE_ROOT /go/src/github.com/articulate/artcli
RUN mkdir -p $SERVICE_ROOT
WORKDIR $SERVICE_ROOT

RUN apt-get update && \
    apt-get install -y \
        python3 \
        python3-pip \
        python3-setuptools \
        groff \
        less \
        git \
    && pip3 install --upgrade pip \
    && apt-get clean

RUN pip3 --no-cache-dir install --upgrade awscli
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && chmod +x /usr/local/bin/dep

# Only used for dev build
RUN go get github.com/mitchellh/gox golang.org/x/tools/cmd/goimports github.com/josephspurrier/goversioninfo/cmd/goversioninfo

COPY Gopkg.lock Gopkg.toml Makefile ./
RUN dep ensure --vendor-only

COPY . ./

ENTRYPOINT [ "/entrypoint.sh" ]
