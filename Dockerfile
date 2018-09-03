FROM golang:alpine

COPY /subscriptions /subscriptions
COPY /service/products/*.json /service/products/

CMD ["/subscriptions"]