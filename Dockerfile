# Build Stage
FROM golang:1.19-buster as builder

# Add Maintainer info
LABEL maintainer="Sanat Tripathy <tripathy.sanat@gmail.com>"

# Set the current working directory inside the container
WORKDIR /app

# Copy the go mod and sum files
COPY go.mod go.sum ./

# Download dependecies
RUN go mod download

# Copy the source from current directory to the working directory inside the container
COPY . .

#Statically compile the binary by setting CGO_ENABLED=0.
#This means that the resulting binary will not be linked to any C libraries.
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

#RUN Stage
FROM gcr.io/distroless/static-debian11
WORKDIR /app
COPY --from=builder /app/api .
EXPOSE 3001

#RUN the output command when the container starts
CMD ["/app/api"]

