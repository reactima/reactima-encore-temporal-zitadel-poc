FROM ubuntu:22.04

RUN apt-get update && \
    apt-get install -y golang git bash mc

RUN mkdir /app

ENV ZITADEL_EXTERNALSECURE=true
ENV ZITADEL_EXTERNALPORT=443
ENV ZITADEL_EXTERNALDOMAIN=login.reactima.com

COPY ./zitadel /app/
RUN chmod +x /app/zitadel

COPY ./zitadel-config /app/

# 1s run
#ENTRYPOINT ["/app/zitadel","start-from-init","--config","/app/prod-zitadel-config.yaml","--steps","/app/prod-zitadel-init-steps.yaml","--masterkey","MasterkeyNeedsToHave32Adeptima##","--tlsMode","external"]

# 2nd+ run
ENTRYPOINT ["/app/zitadel","start-from-setup","--config","/app/prod-zitadel-config.yaml","--steps","/app/prod-zitadel-init-steps.yaml","--masterkey","MasterkeyNeedsToHave32Adeptima##","--tlsMode","external"]

EXPOSE 8080

#CMD ["/bin/bash"]