FROM golang:1.22-alpine

RUN apk add --no-cache git curl screen sqlite-dev gcc musl-dev bash openssh-server sudo

# Setup SSH for gh codespace commands
RUN ssh-keygen -A && \
    echo "root:root" | chpasswd && \
    sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config && \
    sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config

WORKDIR /workspaces/wa-farm

# Pre-clone whatsmeow for faster setup
RUN git clone https://github.com/tulir/whatsmeow /workspaces/wa-farm/whatsmeow

EXPOSE 8080 22

CMD ["sleep", "infinity"]
