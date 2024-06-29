FROM ubuntu:latest

RUN apt update -y && apt install -y openssh-server
RUN mkdir -p /run/sshd

RUN useradd user -m

COPY ./authorized_keys /home/user/.ssh/authorized_keys
RUN chown user:user /home/user/.ssh/authorized_keys
RUN chmod 600 /home/user/.ssh/authorized_keys

CMD /bin/bash -c "/usr/sbin/sshd; sleep infinity"
