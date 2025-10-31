# Software Factory 2.0 - Docker Environment
# Provides all required dependencies in a containerized environment

FROM ubuntu:22.04

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=UTC

# Install base packages
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    git \
    jq \
    make \
    tree \
    sudo \
    locales \
    ca-certificates \
    gnupg \
    lsb-release \
    software-properties-common \
    && rm -rf /var/lib/apt/lists/*

# Set locale
RUN locale-gen en_US.UTF-8
ENV LANG=en_US.UTF-8
ENV LANGUAGE=en_US:en
ENV LC_ALL=en_US.UTF-8

# Install yq
RUN wget https://github.com/mikefarah/yq/releases/download/v4.35.2/yq_linux_amd64 -O /usr/local/bin/yq && \
    chmod +x /usr/local/bin/yq

# Install GitHub CLI
RUN curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg && \
    chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list > /dev/null && \
    apt-get update && \
    apt-get install -y gh && \
    rm -rf /var/lib/apt/lists/*

# Install ripgrep and fd
RUN apt-get update && apt-get install -y \
    ripgrep \
    fd-find \
    && rm -rf /var/lib/apt/lists/* && \
    ln -s /usr/bin/fdfind /usr/local/bin/fd

# Create efforts directory
RUN mkdir -p /efforts && chmod 777 /efforts

# Create working directory
WORKDIR /workspace

# Create a non-root user
RUN useradd -m -s /bin/bash sfuser && \
    echo "sfuser ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

# Copy the Software Factory files
COPY --chown=sfuser:sfuser . /workspace/

# Make scripts executable
RUN chmod +x /workspace/*.sh && \
    chmod +x /workspace/utilities/*.sh

# Switch to non-root user
USER sfuser

# Set up git config (will be overridden by user)
RUN git config --global user.name "Software Factory User" && \
    git config --global user.email "sf@example.com" && \
    git config --global init.defaultBranch main

# Default command
CMD ["/bin/bash"]