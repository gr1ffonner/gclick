FROM clickhouse/clickhouse-server:23.8

# import data into container
COPY ./database/*.sql /docker-entrypoint-initdb.d/
