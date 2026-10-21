# Gator

Gator is a command-line RSS feed aggregator written in Go. It stores feeds and
posts in PostgreSQL, supports local user profiles, and lets each user browse posts
from the feeds they follow.

## Requirements

- Go 1.26.1 or later to build and install the CLI, as specified in `go.mod`.
- PostgreSQL installed and running, with the `psql` client available.
- Git to clone the repository and access its database migrations.
- [Goose](https://github.com/pressly/goose#install) to apply database migrations.

Install Go and PostgreSQL using the instructions for your operating system:
[Go installation](https://go.dev/doc/install) and
[PostgreSQL downloads](https://www.postgresql.org/download/).

Go is required during installation. Once compiled, the `gator` executable runs
without the Go toolchain or source code. It still needs its configuration file,
access to PostgreSQL, and network access to fetch RSS feeds.

## Install

Install the CLI using the lowercase module path declared in `go.mod`:

```sh
go install github.com/andrefetch/gator@latest
```

Go installs executables into `GOBIN` when configured, otherwise into the `bin`
directory under `GOPATH`. With the default Go configuration on macOS or Linux,
add that directory to your shell's `PATH`:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

Add the same export to `~/.zshrc` or `~/.bashrc` to persist it. If you configured
`GOBIN`, add that directory instead. See the
[Go installation command reference](https://go.dev/ref/mod#go-install).

Install the migration tool and clone the repository for its SQL files:

```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
git clone https://github.com/andrefetch/Gator.git
cd Gator
```

To install from this checkout instead, run `go install .` in the repository root.
The following examples use a POSIX shell such as Bash or Zsh.

## Set up PostgreSQL

Start your PostgreSQL service. Connect as a database administrator using `psql`:

```sh
psql -U postgres -d postgres
```

If your installation uses a different administrator role, replace `postgres`
after `-U` with that role. Some Linux installations require
`sudo -u postgres psql` for local administrator access.

Create a dedicated application role and database, then set the role's password
at the prompt:

```sql
CREATE ROLE gator LOGIN;
\password gator
CREATE DATABASE gator OWNER gator;
\q
```

## Configure Gator

Create `.gatorconfig.json` in your home directory, not the repository root:

```json
{
  "db_url": "postgres://gator:YOUR_PASSWORD@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Replace `YOUR_PASSWORD` with the password you set. Percent-encode any reserved
characters in the password when placing it in the URL. Adjust the host, port,
role, and database name if your PostgreSQL setup differs. The example disables
TLS for a local connection; use your database provider's TLS settings for a
remote database.

Restrict access to the file because it contains database credentials:

```sh
chmod 600 ~/.gatorconfig.json
```

Leave `current_user_name` empty initially. The `register` and `login` commands
update it automatically. These commands select local profiles and do not provide
password authentication between Gator users.

## Apply the database schema

From the repository root, set `DB_URL` to the same connection URL used in your
config file, then apply all migrations:

```sh
export DB_URL='postgres://gator:YOUR_PASSWORD@localhost:5432/gator?sslmode=disable'
goose -dir sql/schema postgres "$DB_URL" up
goose -dir sql/schema postgres "$DB_URL" status
```

The migrations create the users, feeds, feed follows, and posts tables, including
feed fetch timestamps and the posts index. Run `up` again after pulling changes
that add migrations. Gator does not apply migrations automatically.

`DB_URL` is used by the Goose command above. Gator reads `db_url` from
`~/.gatorconfig.json`, not from this environment variable.

## Run

Register a user, add an RSS feed, and check your subscriptions:

```sh
gator register alice
gator addfeed "Hacker News" "https://news.ycombinator.com/rss"
gator following
```

Adding a feed also follows it for the current user. Start collecting posts:

```sh
gator agg 1m
```

The aggregator fetches one saved feed immediately, then one feed per interval
across all users' feeds. It prioritizes feeds that have never been fetched, then
the least recently fetched feed. The interval must be a positive Go duration,
such as `10s`, `1m`, or `1h`. Choose an interval that avoids overwhelming feed
servers. Duplicate post URLs are ignored. Fetch and insert errors are logged,
and collection continues.

Leave the aggregator running and use another terminal to browse collected posts:

```sh
gator browse
gator browse 10
```

`browse` defaults to two posts and accepts a positive integer limit. It shows
posts from the current user's followed feeds, newest first. Posts with missing
or unrecognized publication dates appear after dated posts. If no posts appear,
check your subscriptions and allow the aggregator to fetch those feeds.

Press Ctrl+C in the aggregator terminal to stop collection.

## Commands

| Command | Behavior |
| --- | --- |
| `gator register <name>` | Create a user and set it as the current user. |
| `gator login <name>` | Switch to an existing user. |
| `gator users` | List users and mark the current user. |
| `gator addfeed <name> <url>` | Add an RSS feed and follow it as the current user. Quote names containing spaces. |
| `gator follow <url>` | Follow a feed already stored in the database. |
| `gator following` | List the current user's followed feeds. |
| `gator unfollow <url>` | Remove a feed from the current user's subscriptions. |
| `gator agg <interval>` | Continuously fetch saved feeds and store their posts. |
| `gator browse [limit]` | Show saved posts from the current user's subscriptions. Defaults to 2. |
| `gator reset` | Delete all users and, through cascading deletes, all feeds, follows, and posts. No confirmation is requested. |

Register or log in before using `addfeed`, `follow`, `following`, `unfollow`, or
`browse`. Use `reset` only when you intend to erase all application data in the
configured database. It leaves the schema and local configuration file intact.

## Development

From the repository root:

```sh
go test ./...
go build -o gator .
./gator users
```

`go build` creates a standalone executable for the build target. `go install .`
builds and installs it into your Go binary directory. Use `gator` for normal
operation. `go run .` compiles and runs the source and is only for development.

SQL migrations live in `sql/schema`, query definitions in `sql/queries`, and
generated database access code in `internal/database`. Command handlers are in
`internal/commands`, and configuration handling is in `internal/config`.
