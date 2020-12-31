FROM ubuntu:20.10

ENV TZ=Asia/Tokyo 

RUN apt-get update &&  \
    apt-get install -y --no-install-recommends tzdata && \
    apt-get install -y --no-install-recommends \
    make git curl zip tar wget \
    gcc \
    nodejs npm yarn \
    python3 python3-pip \
    php composer \
    openjdk-15-jdk &&\
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y &&\
    wget https://golang.org/dl/go1.15.6.linux-amd64.tar.gz && \
    tar -C /usr/local -xvf go1.15.6.linux-amd64.tar.gz && \
    rm go1.15.6.linux-amd64.tar.gz &&\
    wget https://packages.microsoft.com/config/ubuntu/20.10/packages-microsoft-prod.deb -O packages-microsoft-prod.deb && \
    dpkg -i packages-microsoft-prod.deb &&\
    apt-get install -y --no-install-recommends apt-transport-https &&\
    apt-get update &&\
    apt-get install -y --no-install-recommends dotnet-sdk-5.0 aspnetcore-runtime-5.0 &&\
    rm -rf /var/lib/apt/lists/*

ENV JAVA_HOME=/usr/lib/jvm/openjdk-15-jdk
ENV PATH=$PATH:$JAVA_HOME/bin
ENV PATH=$PATH:/usr/local/go/bin
