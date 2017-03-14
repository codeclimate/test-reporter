FROM golang:1.8

ENV TR=$GOPATH/src/github.com/codeclimate/test-reporter
RUN mkdir -p $TR

WORKDIR $TR
ADD . .

VOLUME /artifacts
RUN rm -rfv ./bin
RUN mkdir ./bin
RUN GOOS=darwin GOARCH=amd64 go build -v -o ./bin/cc-test-reporter-darwin-amd64
RUN GOOS=linux GOARCH=amd64 go build -v -o ./bin/cc-test-reporter-linux-amd64
CMD cp -rv ./bin/* /artifacts/
RUN ls -la /artifacts
