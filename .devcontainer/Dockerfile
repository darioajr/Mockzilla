# Base: Ubuntu 24.04
FROM ubuntu:24.04

# Variáveis
ENV DEBIAN_FRONTEND=noninteractive
ENV GO_VERSION=1.24.2
ENV GOPATH=/home/dev/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

# Instalar dependências
RUN apt-get update && apt-get install -y \
    curl \
    git \
    wrk \
    apache2-utils \
    jq \
    graphviz \
    make \
    gcc \
    bash-completion \
    ca-certificates \
    sudo \
    unzip \
    vim \
    && rm -rf /var/lib/apt/lists/*

# Baixar e instalar o Go
RUN curl -fsSL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz -o go.tar.gz && \
    tar -C /usr/local -xzf go.tar.gz && \
    rm go.tar.gz

# Criar usuário "dev"
RUN useradd -ms /bin/bash dev && \
    echo "dev ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

# Setar diretório de trabalho
USER dev
WORKDIR /home/dev/app

# Criar estrutura do Go
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" "$GOPATH/pkg"

# Entrypoint padrão
CMD ["go", "version"]
