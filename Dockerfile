FROM alpine as builder

WORKDIR /app

COPY ./static ./static
COPY ./templates ./templates
COPY ./strconv.js .
COPY ./main.exe .
CMD ["ls"]
EXPOSE 8080

CMD ["./main.exe"]