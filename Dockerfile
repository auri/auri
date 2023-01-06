FROM alpine

ENV AURI_POSTGRES_PORT=5432
ENV AURI_BUFFALO_PORT=3000
ENV AURI_GO_VERSION=1.19.4
ENV AURI_BUFFALO_VERSION=0.18.4
ENV AURI_ARCH=arm64
#or amd64

ENV PGHOST=db
ENV PGUSER=postgres
ENV PGPASSWORD=postgres

COPY support/dev-env/ /

RUN apk add -U gcompat wget git postgresql12-client

# go setup
RUN <<NUR
  wget --progress=dot:giga -O golang.tgz https://golang.org/dl/go${AURI_GO_VERSION}.linux-${AURI_ARCH}.tar.gz
  tar -C /usr/local -xzf golang.tgz
NUR
ENV PATH=$PATH:/usr/local/go/bin

# buffalo setup
RUN <<NUR
  apk add nodejs npm gcc musl-dev
  git clone https://github.com/gobuffalo/cli.git
  cd cli && git checkout v${AURI_BUFFALO_VERSION}
  go mod tidy && cd cmd/buffalo/ && go build -tags sqlite && mv buffalo /usr/local/go/bin/

  npm install -g yarn
  yarn config set yarn-offline-mirror /npm-packages-offline-cache
  yarn config set yarn-offline-mirror-pruning true

  npm install webpack webpack-cli

  git config --global --add safe.directory /data
NUR

# dependencies for node
RUN <<NUR
  apk add make g++
NUR

# buffalo dev should listen on 0.0.0.0 to get port forwarding working
ENV ADDR=0.0.0.0

RUN mkdir /data
WORKDIR /data
VOLUME /data
CMD docker-entrypoint.sh
