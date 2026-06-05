FROM devopsfaith/krakend:2.4 AS builder

ARG ENV=dev

ENV TERM=xterm-256color
ENV COLORTERM=truecolor

WORKDIR /etc/krakend

COPY krakend-config.tmpl .
COPY config .
COPY jwk /jwk
RUN echo "CURRENT ENVIRONMENT: ${ENV} and $ENV"
COPY config/scripts /etc/krakend/scripts


RUN FC_ENABLE=1 \
    FC_OUT=/tmp/krakend.json \
    FC_PARTIALS="/etc/krakend/partials" \
    FC_SETTINGS="/etc/krakend/settings/$ENV" \
    FC_TEMPLATES="/etc/krakend/templates" \
    SERVICE_NAME="Gateway Service" \
    krakend check -d -t -c krakend-config.tmpl

FROM devopsfaith/krakend:2.4
# MOVE file from tmp/ to /etc/krakend
COPY --from=builder --chown=krakend:nogroup /tmp/krakend.json .
COPY --from=builder --chown=krakend:nogroup /jwk /jwk
COPY --from=builder --chown=krakend:nogroup /etc/krakend/scripts /etc/krakend/scripts