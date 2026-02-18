FROM grafana/k6:latest
COPY scripts /scripts
WORKDIR /scripts
CMD ["run", "/scripts/auth-load-test.js"]
