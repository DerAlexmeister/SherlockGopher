FROM golang:latest as build
WORKDIR /app 
ADD . . 
RUN go build -o main . 

FROM alpine
COPY --from=build /app/main /app/
CMD ["/app/main"]