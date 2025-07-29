FROM docker.io/golang:1.24.5-bookworm AS builder

COPY . /app
WORKDIR /app
RUN go build .


FROM docker.io/debian:bookworm

RUN apt-get update \
  && apt-get install -y chromium \
  && rm -rf /var/lib/apt/lists/*

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/resume-server /app/resume-server
COPY ./assets /app/assets

RUN useradd -ms /bin/bash resumeserver
USER resumeserver:resumeserver

ENTRYPOINT ["/app/resume-server"]
