FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /binjuice

EXPOSE 8000

CMD [ "/binjuice" ]