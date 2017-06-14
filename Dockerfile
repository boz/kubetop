FROM alpine

ADD ./kubetop-linux /kubetop

ENTRYPOINT ./kubetop
