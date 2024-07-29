FROM postgres:latest

# import data into container
COPY ./database/*.sql /docker-entrypoint-initdb.d/
