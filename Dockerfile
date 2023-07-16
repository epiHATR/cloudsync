FROM golang:alpine3.18 as builder
# Set destination for COPY
WORKDIR $GOPATH/cloudsync

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# https://docs.docker.com/engine/reference/builder/#copy
COPY . .
RUN rm -rf examples
RUN rm -rf .github

# Build
ARG COMMIT=1
ARG BUILD=1
ARG RELEASE_DATE='2023-01-01 12:00:00'
ENV COMMIT=${COMMIT}
ENV BUILD=${BUILD}
ENV RELEASE_DATE=${RELEASE_DATE}
RUN GOOS=linux GOARCH=amd64 go build -ldflags "-X 'cloudsync/cmd.commit=$COMMIT' -X 'cloudsync/cmd.build=$BUILD' -X 'cloudsync/cmd.releaseDate=$RELEASE_DATE' -s -w" -a -o /usr/local/bin/cloudsync

# Test
RUN /usr/local/bin/cloudsync version

FROM golang:alpine3.18
WORKDIR /usr/local/bin

COPY --from=builder /usr/local/bin/cloudsync /usr/local/bin/



ENTRYPOINT [ "cloudsync" ]