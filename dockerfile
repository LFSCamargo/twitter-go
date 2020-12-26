# Go version that got used to build the project
FROM golang:1.15.6-alpine

RUN mkdir /app
ADD . /app/
WORKDIR /app

# Builds the go code to the /app folder
RUN go build -o server .

# Executes the compiled Go Code
CMD ["/app/server"]