FROM alpine as builder

WORKDIR /app

COPY ./static ./static
COPY ./templates ./templates
COPY ./strconv.js .
COPY ./bjfu .
RUN ["mkdir safeFiles"]
EXPOSE 8080

CMD ["./bjfu"]