FROM golang:1.20-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o /coinbit-wallet
EXPOSE 3000
CMD [ "/coinbit-wallet" ]