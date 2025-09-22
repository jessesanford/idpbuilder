FROM alpine:latest

# Simple hello world application
RUN apk add --no-cache bash

# Create a simple script
RUN echo '#!/bin/bash' > /hello.sh && \
    echo 'echo "Hello from idpbuilder OCI build!"' >> /hello.sh && \
    echo 'echo "Container built at: $(date)"' >> /hello.sh && \
    echo 'echo "Running on: $(hostname)"' >> /hello.sh && \
    chmod +x /hello.sh

# Add some test files to demonstrate build context
COPY hello.txt /app/hello.txt
COPY README.md /app/README.md

WORKDIR /app

CMD ["/hello.sh"]