# go-wunderkammer

## Important

This is work in progress.

## Tools

### wunderkammer-db

```
$> sqlite3 metmuseum.db < schema/sqlite/oembed.sqlite

$> /usr/local/go-metmuseum-openaccess/bin/emit \
	-oembed \
	-oembed-ensure-images \
	-with-images \
	-bucket-uri file:///usr/local/openaccess \
	-images-bucket-uri file:///usr/local/go-metmuseum-openaccess/data \

   | bin/wunderkammer-db
   	-database-dsn sql:///usr/local/go-wunderkammer/metmuseum.db

$> sqlite3 metmuseum.db 
SQLite version 3.32.1 2020-05-25 16:19:56
Enter ".help" for usage hints.
sqlite> SELECT COUNT(url) FROM oembed;
236288
```

## See also

* https://github.com/aaronland/ios-wunderkammer