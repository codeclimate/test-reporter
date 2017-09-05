FROM java:8

RUN curl -O https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz && \
    tar -xvf go1.8.linux-amd64.tar.gz && \
    mv go /usr/local

ENV PATH=$PATH:/usr/local/go/bin \
    GOPATH=/go
RUN mkdir $GOPATH

ENV PATH=$PATH:/go/bin \
    CCTR=$GOPATH/src/github.com/codeclimate/test-reporter
RUN mkdir -p $CCTR

WORKDIR $CCTR
COPY . .
RUN go install -v && \
    apt-get update && \
    apt-get install -y maven

RUN git clone https://github.com/ale7714/mavanagaiata.git
WORKDIR mavanagaiata
RUN mvn install -DskipTests=true -Dmaven.javadoc.skip=true -B

RUN echo "testing" > ignore.me && \
  git checkout 0.8 && \
  git config --global user.email "you@example.com" && \
  git config --global user.name "Your Name" && \
  git add ignore.me && \
  git commit -m "testing"

RUN mvn cobertura:cobertura

ENV CC_TEST_REPORTER_ID=c4881e09870b0fac1291c93339b36ffe36210a2645c1ad25e52d8fda3943fb4d
RUN test-reporter format-coverage -d -t cobertura target/site/cobertura/coverage.xml && \
    test-reporter upload-coverage -d
