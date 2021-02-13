FROM debian:buster as builder

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
  apt-get install -y protobuf-compiler && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/*

COPY ./protobuf /usr/src/protobuf
COPY python-service /usr/src/python-service/

RUN mkdir /usr/src/python-protobuf-generated && \
  protoc --proto_path=/usr/src \
    --python_out=/usr/src/python-protobuf-generated \
    /usr/src/protobuf/helloworld.proto && \
   cp /usr/src/python-protobuf-generated/protobuf/* /usr/src/python-service/


FROM python:3.8-slim-buster as virtualenv

ENV LANG=C.UTF-8
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV PIP_DISABLE_PIP_VERSION_CHECK=1
ENV VIRTUAL_ENV=/opt/virtualenv

RUN python3 -m venv $VIRTUAL_ENV
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

COPY ./requirements.txt /usr/src/python-requirements/
RUN pip install -r /usr/src/python-requirements/requirements.txt


FROM python:3.8-slim-buster as app
LABEL maintainer="neverping@gmail.com"

ENV VIRTUAL_ENV=/opt/virtualenv
ENV PATH="$VIRTUAL_ENV/bin:$PATH"
ENV LANG=C.UTF-8
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV PIP_DISABLE_PIP_VERSION_CHECK=1

COPY --from=builder /usr/src/python-service/* /usr/src/python-service/
COPY --from=virtualenv /opt/virtualenv /opt/virtualenv

# The app listens on port 50051 by default using PORT environment variable.
ENV PORT=50051
EXPOSE 50051/udp

CMD [ "python", "/usr/src/python-service/greeter_server.py" ]