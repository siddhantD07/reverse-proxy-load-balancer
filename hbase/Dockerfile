#  Author: Hari Sekhon

FROM alpine:latest

ENV PATH $PATH:/hbase/bin

ARG HBASE_VERSION=1.4.11

ENV JAVA_HOME=/usr

LABEL Description="HBase Dev", \
      "HBase Version"="$HBASE_VERSION"

WORKDIR /

RUN \
    apk add --no-cache bash openjdk8-jre-base wget tar openssh && \
    wget --progress=bar:force:noscroll -t 100 --retry-connrefused -O "hbase-$HBASE_VERSION-bin.tar.gz" "https://archive.apache.org/dist/hbase/hbase-1.4.11/hbase-1.4.11-bin.tar.gz" && \
    mkdir hbase-$HBASE_VERSION && \
    tar zxf hbase-$HBASE_VERSION-bin.tar.gz -C hbase-$HBASE_VERSION --strip 1 && \
    ln -sv hbase-$HBASE_VERSION hbase && \
    rm -fv hbase-$HBASE_VERSION-bin.tar.gz && \
    { rm -rf hbase/{docs,src}; : ; } && \
    apk del wget tar

COPY entrypoint.sh /
RUN chmod +x entrypoint.sh
COPY conf/hbase-site.xml /hbase/conf/
COPY profile.d/java.sh /etc/profile.d/
COPY hbase_init.txt /

# Stargate  8080  / 8085
# Thrift    9090  / 9095
# HMaster   16000 / 16010
# RS        16201 / 16301
EXPOSE 2181 8080 8085 9090 9095 16000 16010 16201 16301

CMD "/entrypoint.sh"
