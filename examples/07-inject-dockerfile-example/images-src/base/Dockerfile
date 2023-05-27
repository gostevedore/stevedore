ARG image_from_name
ARG image_from_tag

FROM ${image_from_name}:${image_from_tag}

# Create a new user
RUN echo "anonymous:x:10001:10001:,,,:/app:/bin/sh" >> /etc/passwd && \
    echo "anonymous:x:10001:" >> /etc/group && \
    mkdir -p /app && \
    chown 10001:10001 /app

# Set the user as the default user
USER anonymous
WORKDIR /app
