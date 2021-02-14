FROM debian:buster as builder

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
  apt-get install -y protobuf-compiler && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/*

COPY ./protobuf /usr/local/src/protobuf
COPY python-service /usr/local/src/python-service/

RUN mkdir /usr/local/src/python-protobuf-generated && \
  protoc --proto_path=/usr/local/src \
    --python_out=/usr/local/src/python-protobuf-generated \
    /usr/local/src/protobuf/helloworld.proto && \
   cp /usr/local/src/python-protobuf-generated/protobuf/* /usr/local/src/python-service/


FROM python:3.8-slim-buster as virtualenv

ENV LANG=C.UTF-8
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV PIP_DISABLE_PIP_VERSION_CHECK=1
ENV VIRTUAL_ENV=/opt/virtualenv

RUN python3 -m venv $VIRTUAL_ENV
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

COPY ./requirements.txt /usr/local/src/python-requirements/
RUN pip install -r /usr/local/src/python-requirements/requirements.txt


FROM python:3.8-slim-buster as app
LABEL maintainer="neverping@gmail.com"

ENV VIRTUAL_ENV=/opt/virtualenv
ENV PATH="$VIRTUAL_ENV/bin:$PATH"
ENV LANG=C.UTF-8
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV PIP_DISABLE_PIP_VERSION_CHECK=1

COPY --from=builder /usr/local/src/python-service/* /usr/local/src/python-service/
COPY --from=virtualenv /opt/virtualenv /opt/virtualenv

# The app listens on port 50051 by default using PORT environment variable.
ENV PORT=50051
EXPOSE 50051/tcp

CMD [ "python", "/usr/local/src/python-service/greeter_server.py" ]
