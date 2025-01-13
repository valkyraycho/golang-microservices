FROM postgres:latest

COPY up.sql /docker-entrypoint-initdb.d/

CMD [ "postgres" ]