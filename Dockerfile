
FROM golang:1.22.2-alpine as build-base

WORKDIR /app


COPY . .


RUN go mod download


RUN go build -o main .


ENV ADMIN_USERNAME=adminTax
ENV ADMIN_PASSWORD=admin!
ENV DATABASE_URL=postgres://mhxzvtem:Fu3zMMOsZuiRERnb5s7gPcEHwMxukDDV@rain.db.elephantsql.com
ENV PORT=8080



CMD ["./main"]