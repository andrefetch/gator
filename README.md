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
each feed URL and saves its posts to the database, fetching feeds that have never
been fetched first, then the least recently fetched feed. Duplicate post URLs are
ignored. Failed requests and insert errors are logged and the loop continues.
Leave it running in a separate terminal and press Ctrl+C to stop.

Browse posts from the feeds the current user follows, newest first:

```sh
go run . browse      # defaults to 2 posts
go run . browse 10   # show up to 10 posts
```

The optional limit must be a positive integer. Posts with missing or unrecognized
publication dates are saved without a date and appear after dated posts.
