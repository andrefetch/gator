# Gator

A command-line RSS feed aggregator built in Go that uses HTTP requests to fetch feeds and stores posts in PostgreSQL. Add feeds, manage your subscriptions, switch between local user profiles, and browse the latest posts from the feeds you follow. A collection command fetches feeds at a configurable interval and skips duplicate posts.

## Installation

Requires [Go 1.26.1 or later](https://go.dev/doc/install) and [PostgreSQL](https://www.postgresql.org/download/) installed and running.

```sh
go install github.com/andrefetch/gator@latest
```

Make sure your Go binary directory is on your `PATH`. With the default Go configuration:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

Complete the database setup below before running Gator commands.

### Install from source

```sh
git clone https://github.com/andrefetch/Gator.git
cd Gator
go build -o gator .
```

Use `./gator` instead of `gator` in the examples if you build from source.

### Database setup

Connect to PostgreSQL as an administrator:

```sh
psql -U postgres -d postgres
```

Replace the administrator role if your installation uses a different one. Create a database user, set its password at the prompt, and create the database:

```sql
CREATE ROLE gator LOGIN;
\password gator
CREATE DATABASE gator OWNER gator;
\q
```

Create `~/.gatorconfig.json` with your database connection details:

```json
{
  "db_url": "postgres://gator:YOUR_PASSWORD@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Replace `YOUR_PASSWORD` with the password you set, percent-encoding any reserved URL characters. This connection URL is for a local database. Gator updates `current_user_name` when you register or log in.

```sh
chmod 600 ~/.gatorconfig.json
```

Install [Goose](https://github.com/pressly/goose) to apply the database migrations:

```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
```

If you installed Gator with `go install`, clone the repository for its migration files:

```sh
git clone https://github.com/andrefetch/Gator.git
cd Gator
```

From the repository root, run the migrations using the same connection URL as your config:

```sh
export DB_URL='postgres://gator:YOUR_PASSWORD@localhost:5432/gator?sslmode=disable'
goose -dir sql/schema postgres "$DB_URL" up
```

Gator reads its connection details from `~/.gatorconfig.json`; `DB_URL` is only used by the migration command above.

## Usage & Quick Start

1. Run `gator register <name>` to create and select a local user.
2. Use `gator addfeed <name> <url>` to add an RSS feed and follow it.
3. Run `gator following` to see your subscriptions.
4. Start `gator agg 1m` to collect posts, leaving it running in your terminal.
5. In another terminal, run `gator browse` to read collected posts.

### Example

Run these commands in your terminal after completing setup:

```sh
gator register alice
gator addfeed "Hacker News" "https://news.ycombinator.com/rss"
gator following
gator agg 1m
```

While the aggregator is running, open another terminal to browse posts:

```sh
gator browse
gator browse 10
```

Adding a feed automatically follows it for the current user. The aggregator fetches one saved feed immediately, then one per interval across all users' feeds. Use a positive duration such as `30s`, `1m`, or `1h`. Press Ctrl+C to stop collection.

`browse` shows the two newest posts from your subscriptions by default. Pass a positive number to show more. If no posts appear, let the aggregator fetch your feeds first.

### Commands

| Command | Description |
| --- | --- |
| `register <name>` | Create a user and make it the current user. |
| `login <name>` | Switch to an existing user. |
| `users` | List users and mark the current user. |
| `addfeed <name> <url>` | Add a feed and follow it. Quote names containing spaces. |
| `follow <url>` | Follow a feed already in the database. |
| `following` | List your followed feeds. |
| `unfollow <url>` | Unfollow a feed. |
| `agg <interval>` | Continuously fetch feeds and save posts. |
| `browse [limit]` | Browse posts from your subscriptions. Defaults to 2. |
| `reset` | Delete all users, feeds, follows, and posts without confirmation. |

Prefix each command with `gator`. User profiles select whose subscriptions you manage; they do not use passwords for authentication.

Your users, feeds, subscriptions, and collected posts are stored in PostgreSQL and persist when you exit.
