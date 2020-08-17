FROM devopsfaith/krakend:1.1.1
USER root
RUN apt-get update -qq && \
    apt-get install python3 -qq > /dev/null

COPY ./krakend /etc/krakend
COPY ./docker-entrypoint.sh .

# Overrides krakend's parent image entrypoint
ENTRYPOINT []
CMD sh docker-entrypoint.sh
