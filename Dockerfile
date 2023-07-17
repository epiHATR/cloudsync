#####################################################################################################################
#####################################################################################################################
FROM golang:alpine3.18 as builder

#SET WORKING DIRECTORY FOR PROJECT AND GOROOT
WORKDIR $GOPATH/cloudsync

#DOWNLOAD MODULES
COPY go.mod go.sum ./
RUN go mod download

# PREPARE SOURCES
COPY . .

# HANDLER BUILD ARGUMENTS
ARG COMMIT=1
ARG BUILD=1
ARG RELEASE_DATE='2023-01-01 12:00:00'

ENV COMMIT=${COMMIT}
ENV BUILD=${BUILD}
ENV RELEASE_DATE=${RELEASE_DATE}

#BUILD
RUN GOOS=linux GOARCH=amd64 go build -ldflags "-X 'cloudsync/cmd.commit=$COMMIT' -X 'cloudsync/cmd.build=$BUILD' -X 'cloudsync/cmd.releaseDate=$RELEASE_DATE' -s -w" -a -o /tmp/builder/cloudsync

#TEST COMMANDS
RUN /tmp/builder/cloudsync version

######################################################################################################################
######################################################################################################################
FROM golang:alpine3.18
WORKDIR /usr/local/bin
COPY --from=builder /tmp/builder/cloudsync .
ENTRYPOINT [ "cloudsync" ]