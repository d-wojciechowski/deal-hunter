FROM alpine:latest

WORKDIR /
ADD resources /resources
ADD /distr/DealHunter-linux-x64 /DealHunter-linux-x64

# executable
#RUN addgroup -S appgroup && adduser -S nonroot -G appgroup
#USER nonroot
HEALTHCHECK --interval=5s --timeout=5s --retries=3 CMD wget localhost:8086/health -q -O - > /dev/null 2>&1
VOLUME ["/logs"]
ENTRYPOINT [ "./DealHunter-linux-x64" ]
