FROM python:3.11-alpine

COPY ./reqirements.txt /app/reqirements.txt

WORKDIR /app

RUN pip install -r reqirements.txt

COPY . /app

EXPOSE 1234

ENTRYPOINT ["gunicorn"]

CMD ["-b", "0.0.0.0:1234", "wsgi:app"]

