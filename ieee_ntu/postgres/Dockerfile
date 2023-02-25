FROM postgres:latest

ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD admin
ENV POSTGRES_DB dev

COPY init.sql /docker-entrypoint-initdb.d/

COPY accounts.sql /