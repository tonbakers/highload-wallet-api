FROM ubuntu:20.04 as ton_builder

ENV DEBIAN_FRONTEND=noninteractive
ENV CC=/usr/bin/clang
ENV CXX=/usr/bin/clang++
ENV CCACHE_DISABLE=1

RUN apt-get update \
    && apt-get install --no-install-recommends -y \
        build-essential \
        ca-certificates \
        clang \
        cmake \
        git \
        libgsl-dev \
        libmicrohttpd-dev \
        libreadline-dev \
        libssl-dev \
        pkg-config \
        wget \
        zlib1g-dev \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# build fift and func
RUN git clone --recurse-submodules --shallow-submodules --depth 1 --no-tags https://github.com/ton-blockchain/ton.git \
    && cmake -DCMAKE_BUILD_TYPE=Release /ton \
    && make -j 4 fift func


# Build golang app
FROM golang:1.17 as go_builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build


# Combine binary with libs
FROM ubuntu:20.04

EXPOSE 8091

RUN apt-get update \
    && apt-get install --no-install-recommends -y \
        ca-certificates \
        curl \
        libssl-dev \
        libatomic1 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY --from=ton_builder /crypto /usr/bin/ton/crypto
COPY --from=ton_builder /ton/crypto /usr/src/ton/crypto

WORKDIR /app
COPY config.json ./
COPY contract/ contract/
COPY static/ static/
COPY --from=go_builder /app/highload-wallet-api /app/highload-wallet-api
RUN cd /app/contract/ && bash wallet.sh
RUN mkdir /app/temp && mkdir /app/temp/bocs && mkdir /app/temp/orders


CMD ./highload-wallet-api

