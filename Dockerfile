FROM golang:bookworm AS nkfproxy

ARG IGNORECACHE=0

ADD ./nkfproxy /nkfproxy
RUN --mount=type=cache,target=/go \
    cd /nkfproxy \
    && echo "go get" \
    && go get -d \
    && echo "go build" \
    && GOCACHE=/go/.cache CGO_ENABLED=0 go build -ldflags='-s -w'

FROM golang:bookworm AS sedproxy

ARG IGNORECACHE=0

ADD ./sedproxy /sedproxy
RUN --mount=type=cache,target=/go \
    cd /sedproxy \
    && echo "go get" \
    && go get -d \
    && echo "go build" \
    && GOCACHE=/go/.cache CGO_ENABLED=0 go build -ldflags='-s -w'

FROM golang:bookworm AS awkproxy

ARG IGNORECACHE=0

ADD ./awkproxy /awkproxy
RUN --mount=type=cache,target=/go \
    cd /awkproxy \
    && echo "go get" \
    && go get -d \
    && echo "go build" \
    && GOCACHE=/go/.cache CGO_ENABLED=0 go build -ldflags='-s -w'

FROM golang:bookworm AS miniproxy

ARG IGNORECACHE=0

ADD ./miniproxy /miniproxy
RUN --mount=type=cache,target=/go \
    cd /miniproxy \
    && echo "go get" \
    && go get -d \
    && echo "go build" \
    && GOCACHE=/go/.cache CGO_ENABLED=0 go build -ldflags='-s -w'

FROM ateliersjp/openresty-lolhtml:bionic

RUN apt-get update && \
    apt-get install -y \
    wget \
    git \
    gcc \
    make \
    cargo \
    unzip && \
    luarocks install --server=https://luarocks.org/dev lolhtml

COPY --from=nkfproxy /nkfproxy/nkfproxy /bin/
COPY --from=sedproxy /sedproxy/sedproxy /bin/
COPY --from=awkproxy /awkproxy/awkproxy /bin/
COPY --from=miniproxy /miniproxy/miniproxy /bin/

COPY ./start.sh /bin/
COPY ./default.conf /etc/nginx/conf.d/
COPY ./*.lua /usr/local/openresty/nginx/

CMD [ "start.sh" ]

EXPOSE 8080
