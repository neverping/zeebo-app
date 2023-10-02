FROM python:3.11-slim-bookworm as builder

ENV LANG=C.UTF-8
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV PIP_DISABLE_PIP_VERSION_CHECK=1
ENV VIRTUAL_ENV=/opt/virtualenv

RUN python3 -m venv $VIRTUAL_ENV
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

COPY ./protobuf /usr/local/src/protobuf
COPY ./requirements.txt /usr/local/src/python-requirements/
RUN pip install -r /usr/local/src/python-requirements/requirements.txt

## Gen protobuf
RUN mkdir -p /usr/local/src/python-protobuf-generated/
RUN python3 -m grpc_tools.protoc \
  -I/usr/local/src/protobuf \
  --python_out=/usr/local/src/python-protobuf-generated \
  --pyi_out=/usr/local/src/python-protobuf-generated \
  --grpc_python_out=/usr/local/src/python-protobuf-generated \
  /usr/local/src/protobuf/helloworld.proto

FROM python:3.11-slim-bookworm as app
LABEL maintainer="neverping@gmail.com"

ENV VIRTUAL_ENV=/opt/virtualenv
ENV PATH="$VIRTUAL_ENV/bin:$PATH"
ENV LANG=C.UTF-8
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV PIP_DISABLE_PIP_VERSION_CHECK=1

COPY ./python-service /usr/local/src/python-service/

# HINT: We will replace generated versioned protobuf files with the ones generated from above
COPY --from=builder /usr/local/src/python-protobuf-generated/* /usr/local/src/python-service/
COPY --from=builder /opt/virtualenv /opt/virtualenv

# The app listens on port 50051 by default using PORT environment variable.
ENV PORT=50051
EXPOSE 50051/tcp

CMD [ "python", "/usr/local/src/python-service/greeter_server.py" ]