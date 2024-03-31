FROM mongo:7.0
RUN openssl rand -base64 756 > /data/replica.key
RUN chmod 400 /data/replica.key
RUN chown mongodb:mongodb /data/replica.key