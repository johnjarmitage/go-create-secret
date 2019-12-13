FROM debian

WORKDIR app

COPY . /app/.

CMD ./secret-maker