FROM alpine:latest

WORKDIR /
ADD /distr/DealHunter-linux-x64 /DealHunter-linux-x64

HEALTHCHECK --interval=5s --timeout=5s --retries=3 CMD wget localhost:8086/health -q -O - > /dev/null 2>&1
VOLUME ["/logs","/resources"]
ENTRYPOINT [ "./DealHunter-linux-x64" ]
