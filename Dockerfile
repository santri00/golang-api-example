FROM golang:1.15-alpine
## Set GOPATH env variable
ENV GOPATH /app
## Set the api directory as the Working directory
## to execute any further commands
WORKDIR /app/src/github.com/stripathy/demo-vehicle-api
## Copy code to the working directory from the root directory
COPY . /app/src/github.com/stripathy/demo-vehicle-api
EXPOSE 3001
RUN go install
CMD $GOPATH/bin/demo-vehicle-api


