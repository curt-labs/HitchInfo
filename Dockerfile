FROM golang:1.4

RUN mkdir -p /home/deployer/gosrc/src/github.com/curt-labs/HitchInfo
ADD . /home/deployer/gosrc/src/github.com/curt-labs/HitchInfo
WORKDIR /home/deployer/gosrc/src/github.com/curt-labs/HitchInfo
RUN export GOPATH=/home/deployer/gosrc && go get
RUN export GOPATH=/home/deployer/gosrc && go build -o HitchInfo ./index.go

ENTRYPOINT /home/deployer/gosrc/src/github.com/curt-labs/HitchInfo/HitchInfo -http=:8087

EXPOSE 8087
