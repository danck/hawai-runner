FROM golang

ADD . /go/src/HAWAI/repos/hawai-logginghub

RUN go install HAWAI/repos/hawai-logginghub

ENTRYPOINT /go/bin/hawai-logginghub 

EXPOSE 20000 

