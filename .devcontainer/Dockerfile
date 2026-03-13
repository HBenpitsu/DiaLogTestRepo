# syntax=docker/dockerfile:1.7
FROM mcr.microsoft.com/devcontainers/typescript-node:4-24-trixie

ARG ACT_VERSION=v0.2.80
ARG GEMINI_CLI_VERSION=latest

RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
	--mount=type=cache,target=/var/lib/apt/lists,sharing=locked \
	set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends \
		ca-certificates \
		curl \
		gnupg \
		golang-go \
		lsb-release \
		unzip; \
	curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg \
		| gpg --dearmor -o /usr/share/keyrings/cloud.google.gpg; \
	echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" \
		> /etc/apt/sources.list.d/google-cloud-sdk.list; \
	apt-get update; \
	apt-get install -y --no-install-recommends google-cloud-cli

# Install act (GitHub Actions local runner).
RUN set -eux; \
	arch="$(dpkg --print-architecture)"; \
	case "${arch}" in \
	  amd64) act_arch="x86_64" ;; \
	  arm64) act_arch="arm64" ;; \
	  *) echo "Unsupported architecture: ${arch}"; exit 1 ;; \
	esac; \
	curl -fsSL -o /tmp/act.tar.gz "https://github.com/nektos/act/releases/download/${ACT_VERSION}/act_Linux_${act_arch}.tar.gz"; \
	tar -xzf /tmp/act.tar.gz -C /usr/local/bin act; \
	chmod +x /usr/local/bin/act; \
	rm -f /tmp/act.tar.gz

# Install Gemini CLI globally and reuse npm cache across rebuilds.
RUN --mount=type=cache,target=/root/.npm,sharing=locked \
	npm install -g "@google/gemini-cli@${GEMINI_CLI_VERSION}"

COPY .devcontainer/scripts/copy-host-ssh.sh /usr/local/bin/copy-host-ssh.sh
RUN chmod +x /usr/local/bin/copy-host-ssh.sh
