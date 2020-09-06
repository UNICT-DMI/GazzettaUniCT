FROM golang

RUN mkdir -p /app/
ADD . /go/src/GazzettaUniCT/

ENV GO111MODULE=on

RUN cd /go/src/GazzettaUniCT/ && go build -o GazzettaUniCT

WORKDIR /go/src/GazzettaUniCT/

ENTRYPOINT [ "/go/src/GazzettaUniCT/GazzettaUniCT" ]

