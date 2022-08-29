FROM alpine as builder

WORKDIR /app

COPY ./static ./static
COPY ./templates ./templates
ADD ./strconv.js ./safeFiles
COPY ./bjfu .
EXPOSE 8080

CMD ["./bjfu"]