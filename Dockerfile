FROM golang:1.12.6-stretch

RUN apt-get update -y && DEBIAN_FRONTEND=noninteractive apt-get install -yq \
        build-essential \
        libffi-dev \
        curl \
        make \
        libxml2-dev \
        libxml2 \
        bzr \
        && apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY . .

RUN go build

ENTRYPOINT ./txdb_ejemplo
