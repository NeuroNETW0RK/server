FROM golang:latest
# RUN echo -e http://mirrors.ustc.edu.cn/alpine/v3.7/main/ > /etc/apk/repositories

ENV GO111MODULE=on
ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

# 安装必要的软件包和依赖包
USER root
RUN sed -i 's/deb.debian.org/mirrors.tuna.tsinghua.edu.cn/' /etc/apt/sources.list && \
    sed -i 's/security.debian.org/mirrors.tuna.tsinghua.edu.cn/' /etc/apt/sources.list && \
    sed -i 's/security-cdn.debian.org/mirrors.tuna.tsinghua.edu.cn/' /etc/apt/sources.list && \
    apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends \
    curl \
    zip \
    unzip \
    git \
    vim \
    openssh-server

ENV PATH $GOPATH/bin:$PATH

RUN apt-get install -y openssh-server && service ssh start && \
  ssh-keygen -f ~/.ssh/id_rsa -t rsa -N '' && cp ~/.ssh/id_rsa.pub ~/.ssh/authorized_keys \
  && echo "  StrictHostKeyChecking no" >> /etc/ssh/ssh_config \
  && echo "  Port 2233" >> /etc/ssh/ssh_config \
  && echo "Port 2233" >> /etc/ssh/sshd_config

# 清理垃圾
USER root
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
    rm /var/log/lastlog /var/log/faillog

EXPOSE 8196
EXPOSE 2233

WORKDIR /app

COPY . .
RUN go mod download

RUN  go get github.com/cosmtrek/air && go install github.com/cosmtrek/air



FROM ubuntu:18.04

RUN echo 'deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic main restricted universe multiverse\n\
deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-updates main restricted universe multiverse\n\
deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-backports main restricted universe multiverse\n\
deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-security main restricted universe multiverse' > /etc/apt/sources.list

RUN mkdir $HOME/.pip

RUN echo '[global] \n\
trusted-host=pypi.tuna.tsinghua.edu.cn \n\
index-url=https://pypi.tuna.tsinghua.edu.cn/simple' > $HOME/.pip/pip.conf

RUN apt clean && apt-get update && apt-get upgrade -y \
  && DEBIAN_FRONTEND=noninteractive apt-get install -y tzdata \
  && cp -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
  && dpkg-reconfigure --frontend noninteractive tzdata \
  && apt-get install -y wget


RUN wget -q http://11.1.203.3:8787/download/Miniconda3-py39_4.12.0-Linux-x86_64.sh && bash ./Miniconda3-py39_4.12.0-Linux-x86_64.sh -b -p $HOME/miniconda3 && rm ./Miniconda3-py39_4.12.0-Linux-x86_64.sh

RUN echo 'channels:\n\
 - defaults\n\
show_channel_urls: true\n\
default_channels:\n\
 - https://mirrors.tuna.tsinghua.edu.cn/anaconda/pkgs/main\n\
 - https://mirrors.tuna.tsinghua.edu.cn/anaconda/pkgs/r\n\
 - https://mirrors.tuna.tsinghua.edu.cn/anaconda/pkgs/msys2\n\
custom_channels:\n\
 conda-forge: https://mirrors.tuna.tsinghua.edu.cn/anaconda/cloud\n\
 msys2: https://mirrors.tuna.tsinghua.edu.cn/anaconda/cloud\n\
 bioconda: https://mirrors.tuna.tsinghua.edu.cn/anaconda/cloud\n\
 menpo: https://mirrors.tuna.tsinghua.edu.cn/anaconda/cloud\n\
 pytorch: https://mirrors.tuna.tsinghua.edu.cn/anaconda/cloud\n\
 pytorch-lts: https://mirrors.tuna.tsinghua.edu.cn/anaconda/cloud\n\
 simpleitk: https://mirrors.tuna.tsinghua.edu.cn/anaconda/cloud' > $HOME/.condarc

WORKDIR /server
COPY . .

RUN eval "$(/root/miniconda3/bin/conda shell.bash hook)" && conda init && pip config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple && \
    conda install protobuf==3.19.6 && \
    pip install -r requirements.txt  && \
    pip install redis && \
    pip install scipy
