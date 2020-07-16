# go-wunderkammer

Go package for working with "[wunderkammer](https://github.com/aaronland/ios-wunderkammer)" databases.

## Tools

To build binary versions of these tools run the `cli` Makefile target. For example:

```
$> make cli
go build -mod vendor -o bin/wunderkammer-db cmd/wunderkammer-db/main.go
```

### append-dataurl

The `append-dataurl` tool that used to be part of this package has been moved in to the [go-wunderkammer-image](https://github.com/aaronland/go-wunderkammer-image) package.

### wunderkammer-db

Create a wunderkammer database from a stream of line-separated OEmbed JSON records.

```
> ./bin/wunderkammer-db -h
Usage of ./bin/wunderkammer-db:
  -database-dsn string
    	A valid wunderkammer database DSN string. (default "sql://sqlite3/oembed.db")
```

For example:

```
$> sqlite3 metmuseum.db < schema/sqlite/oembed.sqlite

# as in: https://github.com/aaronland/go-metmuseum-openaccess#emit

$> /usr/local/go-metmuseum-openaccess/bin/emit \
	-oembed \
	-oembed-ensure-images \
	-with-images \
	-bucket-uri file:///usr/local/openaccess \
	-images-bucket-uri file:///usr/local/go-metmuseum-openaccess/data \

   | bin/wunderkammer-db
   	-database-dsn sql://sqlite3/usr/local/go-wunderkammer/metmuseum.db

$> sqlite3 metmuseum.db 
SQLite version 3.32.1 2020-05-25 16:19:56
Enter ".help" for usage hints.
sqlite> SELECT COUNT(url) FROM oembed;
236288
```

## Database DSN strings

Database DSN strings are URIs in the form of `{DATABASE_CLASS}://{DATABASE_DRIVER}{DATABASE_PATH}`

For example: `sql://sqlite3/usr/local/oembed.db`

### Supported database classes

#### sql

Index records with any valid Go language `database/sql` driver assuming its been imported by your code. Currently the only default driver that is included is the [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) `sqlite3` driver.

## See also

* https://github.com/search?q=topic%3Awunderkammer+org%3Aaaronland&type=Repositories
