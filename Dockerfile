FROM alpine as builder

WORKDIR /app

COPY ./static .
COPY ./templates .
EXPOSE 8080

CMD ["./main.exe"]