FROM alpine:3.19

RUN apk add --no-cache bash jq coreutils util-linux

ENV HOME=/root
ENV DOTCLAUDE_REPO_DIR=/dotclaude
ENV CLAUDE_DIR=/root/.claude

WORKDIR /dotclaude

COPY base/ /dotclaude/base/
COPY profiles/ /dotclaude/profiles/

RUN mkdir -p /root/.claude && \
    chmod +x /dotclaude/base/scripts/dotclaude && \
    ln -s /dotclaude/base/scripts/dotclaude /usr/local/bin/dotclaude

CMD ["bash"]
