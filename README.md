# Gator

Apply the database migrations before running the aggregator:

```sh
goose -dir sql/schema postgres "$DB_URL" up
```

Collect posts from saved feeds immediately, then fetch one feed each minute:

```sh
go run . agg 1m
```

The interval must be a positive duration, such as `10s`, `1m`, or `1h`.
Choose an interval that avoids overwhelming feed servers. The aggregator prints
each feed URL and its post titles, fetching feeds that have never been fetched
first, then the least recently fetched feed. Failed requests are logged and the
loop continues. Leave it running in a separate terminal and press Ctrl+C to stop.
