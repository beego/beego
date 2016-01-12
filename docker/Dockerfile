FROM golang
MAINTAINER Shoother shoother.li@gmail.com

# Install packages
RUN apt-get update
RUN apt-get install -y supervisor
RUN mkdir -p /var/log/supervisor

# Install beego and bee
RUN go get github.com/astaxie/beego
RUN go get github.com/beego/bee

# 对外暴露8080端口
EXPOSE     8080

# 设置容器默认工作目录为/go
WORKDIR    /go

RUN cd /go/src && bee new beego
RUN touch /etc/supervisor/conf.d/beego.conf
RUN echo "[supervisord]" >> /etc/supervisor/conf.d/beego.conf
RUN echo "nodaemon=true" >> /etc/supervisor/conf.d/beego.conf
RUN echo "[program:beego]" >> /etc/supervisor/conf.d/beego.conf
RUN echo "directory=/go/src/beego" >> /etc/supervisor/conf.d/beego.conf
RUN echo "command=bee run" >> /etc/supervisor/conf.d/beego.conf
RUN echo "user = root" >> /etc/supervisor/conf.d/beego.conf
RUN echo "autostart = true" >> /etc/supervisor/conf.d/beego.conf

CMD ["/usr/bin/supervisord"]
