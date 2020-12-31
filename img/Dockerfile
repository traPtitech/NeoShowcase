FROM ubuntu:20.10

RUN apt-get update &&  \
    apt-get install -y --no-install-recommends tzdata && \
    apt-get install -y --no-install-recommends \
    make git curl wget unzip \
    gcc \
    nodejs npm yarn \
    python3 python3-pip \
    php composer \
    openjdk-15-jdk &&\
    curl -sSf https://sh.rustup.rs | sh -s -- -y &&\
    wget https://golang.org/dl/go1.15.6.linux-amd64.tar.gz && \
    tar -C /usr/local -xvf go1.15.6.linux-amd64.tar.gz && \
    wget https://packages.microsoft.com/config/ubuntu/20.10/packages-microsoft-prod.deb -O packages-microsoft-prod.deb && \
    dpkg -i packages-microsoft-prod.deb &&\
    apt-get install -y --no-install-recommends apt-transport-https &&\
    apt-get update &&\
    apt-get install -y --no-install-recommends dotnet-sdk-5.0 aspnetcore-runtime-5.0 &&\
    wget https://services.gradle.org/distributions/gradle-6.7.1-bin.zip &&\
    mkdir /opt/gradle &&\
    unzip -d /opt/gradle gradle-6.7.1-bin.zip &&\
    rm -rf /var/lib/apt/lists/* packages-microsoft-prod.deb gradle-6.7.1-bin.zip go1.15.6.linux-amd64.tar.gz

ENV JAVA_HOME=/usr/lib/jvm/java-15-openjdk-amd64
ENV PATH=$PATH:$JAVA_HOME/bin
ENV PATH=$PATH:/usr/local/go/bin
ENV PATH=$PATH:/opt/gradle/gradle-6.7.1/bin