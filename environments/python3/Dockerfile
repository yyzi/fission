FROM alpine:3.4

RUN echo "@main http://dl-4.alpinelinux.org/alpine/v3.4/main" >> /etc/apk/repositories \
  && apk add --update \
              python3 \
	      python3-dev \
#	      build-base \
  && python3 -m ensurepip \
  && pip3 install --upgrade pip \
  && rm /var/cache/apk/*

ADD . /app

RUN pip3 install -r /app/requirements.txt

RUN cd /usr/bin \
  && ln -sf python3 python

ENTRYPOINT ["python", "./app/server.py"]
