# go-wunderkammer

Go package for working with "[wunderkammer](https://github.com/aaronland/ios-wunderkammer)" databases.

## Tools

To build binary versions of these tools run the `cli` Makefile target. For example:

```
> make cli
go build -mod vendor -o bin/emit cmd/emit/main.go
go build -mod vendor -o bin/wunderkammer-db cmd/wunderkammer-db/main.go
```

### append-dataurl

The `append-dataurl` tool that used to be part of this package has been moved in to the [go-wunderkammer-image](https://github.com/aaronland/go-wunderkammer-image) package.

### emit

Emit the contents of a wunderkammer database as a stream of OEmbed records.

```
> ./bin/emit -h
Usage of ./bin/emit:
  -database-dsn string
    	A valid wunderkammer database DSN string. (default "sql://sqlite3/oembed.db")
  -format-json
    	Format JSON output for each record.
  -json
    	Emit a JSON list.
  -null
    	Emit to /dev/null
  -query value
    	One or more {PATH}={REGEXP} parameters for filtering records.
  -query-mode string
    	Specify how query filtering should be evaluated. Valid modes are: ALL, ANY (default "ALL")
  -stdout
    	Emit to STDOUT (default true)
```

For example, here's how you might "emit" a wunderkammer database produced by the `wunderkammer-db` tool described below:

```
$> ./bin/emit \
	-format-json \
	-database-dsn 'sql://usr/local/go-wunderkammer/hmsg.db'

{
  "version": "1.0",
  "type": "photo",
  "width": -1,
  "height": -1,
  "title": "Untitled (One Of Six Prints) (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)",
  "url": "https://ids.si.edu/ids/download?id=HMSG-HMSG_19865656_20150513_001_screen",
  "author_name": "Abraham Walkowitz, American, b. Tyumen, Russia, 1878–1965",
  "author_url": "https://hirshhorn.si.edu/search-results/search-result-details/?edan_search_value=hmsg_86.5656",
  "provider_name": "Hirshhorn Museum and Sculpture Garden",
  "provider_url": "https://hirshhorn.si.edu",
  "object_uri": "si://hmsg/o/86_5656",
  "data_url": "data:image/jpeg;base64,R0lGODlh9AFeAYcAAAAAAAAARAAAiAAAzABEAABERABE ...and so on
}
... and so on
```

#### JSON

By default records are emitted as a stream of line-separated JSON. If you want to emit a valid JSON list then you would pass in the `-json` flag. For example:

```
$> ./bin/emit -database-dsn 'sql://sqlite3/usr/local/go-wunderkammer/hmsg.db' \
	-json \

   | jq '.[]["author_name"]' \
   | sort \
   | uniq

"Abastenia St. Leger Eberle, American, b. Webster City, Iowa, 1878–1942"
"Abbott Thayer, American, 1849–1921"
"Abraham Walkowitz, American, b. Tyumen, Russia, 1878–1965"
"Albert Bierstadt, American, b. Solingen, Germany, 1830–1902"
"Albert Carrier-Belleuse, French, b. Anizy-le-Chateau, 1824–1887"
"Albrecht Dürer, German, b. Nuremberg, 1471–1528"
"Alexandre Falguiere, French, b. Toulouse, 1831–1900"
"Alfred Henry Maurer, American, b. New York City, 1868–1932"
"Antoine-Louis Barye, French, b. Paris, 1795–1875"
... and so on
```

#### Inline queries

You can also specify inline queries by passing a `-query` parameter which is a string in the format of:

```
{PATH}={REGULAR EXPRESSION}
```

Paths follow the dot notation syntax used by the [tidwall/gjson](https://github.com/tidwall/gjson) package and regular expressions are any valid [Go language regular expression](https://golang.org/pkg/regexp/). Successful path lookups will be treated as a list of candidates and each candidate's string value will be tested against the regular expression's [MatchString](https://golang.org/pkg/regexp/#Regexp.MatchString) method.

For example, all the records whose `author_name` contains the name "Abraham" output as a valid JSON list and sent to the `jq` tool (printing the object title):

```
$> ./bin/emit \
	-database-dsn 'sql://sqlite3/usr/local/go-wunderkammer/hmsg.db' \
	-query 'author_name=(?i)Abraham' \
	-json \

   | jq '.[]["title"]'

"Untitled (One Of Six Prints) (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Self-Portrait (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Abstraction (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, Gift of Joseph H. Hirshhorn, 1966)"
"Three Women (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Figure Sketch (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Isadora (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Head Of A Man (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"The City (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Under Two Trees (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"In The Street (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, Gift of Joseph H. Hirshhorn, 1966)"
"The Road, Paris (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, Gift of Joseph H. Hirshhorn, 1966)"
"Man And Woman (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"(Untitled) (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Street Scene, Building, Or Portrait (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Park With Figures (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Figure Sketch (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Street Scene, Building, Or Portrait (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Seated Woman (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, Gift of Joseph H. Hirshhorn, 1966)"
```

The default query mode is to ensure that all queries match but you can also specify that only one or more queries need to match by passing the `-query-mode ANY` flag:

```
$> ./bin/emit -database-dsn 'sql://sqlite3/usr/local/go-wunderkammer/hmsg.db' \
	-query 'author_name=(?i)abraham' \
	-query 'author_name=(?i)mary' \
	-query-mode ANY \
	-json \

   | jq '.[]["title"]'

"Untitled (One Of Six Prints) (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Self-Portrait (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Abstraction (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, Gift of Joseph H. Hirshhorn, 1966)"
"Three Women (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Figure Sketch (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Isadora (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Young Girl Reading (Jeune Fille Lisant) (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Head Of A Man (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"The City (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Under Two Trees (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Woman in Raspberry Costume Holding a Dog (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, Gift of Joseph H. Hirshhorn, 1972)"
"In The Street (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, Gift of Joseph H. Hirshhorn, 1966)"
"The Road, Paris (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, Gift of Joseph H. Hirshhorn, 1966)"
"Man And Woman (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"(Untitled) (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Street Scene, Building, Or Portrait (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Park With Figures (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Figure Sketch (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Street Scene, Building, Or Portrait (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, The Joseph H. Hirshhorn Bequest, 1981)"
"Baby Charles (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, Gift of the Joseph H. Hirshhorn Foundation, 1966)"
"Seated Woman (Hirshhorn Museum and Sculpture Garden, Smithsonian Institution, Washington, DC, Gift of Joseph H. Hirshhorn, 1966)"
```

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
