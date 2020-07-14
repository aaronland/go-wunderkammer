# go-wunderkammer

## Important

This is work in progress.

## Tools

```
$> make cli
go build -mod vendor -o bin/wunderkammer-db cmd/wunderkammer-db/main.go
go build -mod vendor -o bin/append-dataurl cmd/append-dataurl/main.go
```

### append-dataurl

```
> ./bin/append-dataurl -h
Usage of ./bin/append-dataurl:
  -content-aware-height int
    	Content aware resizing to this height.
  -content-aware-resize
    	Enable content aware (seam carving) resizing.
  -content-aware-width int
    	Content aware resizing to this width.
  -dither
    	Dither (halftone) the final image.
  -format
    	Emit results as formatted JSON.
  -json
    	Emit results as a JSON array.
  -null
    	Emit to /dev/null
  -overwrite
    	Overwrite exisiting data_url properties
  -resize
    	Resize images to a maximum dimension (preserving aspect ratio).
  -resize-max-dimension int
    	Resize images to this maximum height or width (preserving aspect ratio).
  -stdout
    	Emit to STDOUT (default true)
```

For example:

```
$> /usr/local/go-smithsonian-openaccess/bin/emit \
	-oembed \
	-bucket-uri file:///Users/asc/code/OpenAccess metadata/objects/NASM \

   | bin/append-dataurl \
	-timings \
	-dither \

   | bin/wunderkammer-db \
	-database-dsn 'sql://sqlite3/usr/local/go-wunderkammer/nasm.db'

2020/07/14 09:04:40 Time to wait to process http://ids.si.edu/ids/deliveryService?id=NASM-A19670206000_PS01, 412ns
2020/07/14 09:04:40 Time to wait to process http://ids.si.edu/ids/deliveryService?id=NASM-NASM2011-00584, 254ns
2020/07/14 09:04:40 Time to wait to process https://ids.si.edu/ids/download?id=NASM-A19350058000-NASM2019-01744_screen, 112ns
2020/07/14 09:04:40 Time to wait to process https://ids.si.edu/ids/download?id=NASM-A19350058000-NASM2019-00617_screen, 162ns
2020/07/14 09:04:40 Time to wait to process https://ids.si.edu/ids/download?id=NASM-A19350058000-NASM2019-01743_screen, 111ns
2020/07/14 09:04:40 Time to wait to process https://ids.si.edu/ids/download?id=NASM-A19350058000-NASM2019-01745_screen, 90ns
2020/07/14 09:04:40 Time to wait to process https://ids.si.edu/ids/download?id=NASM-A19350058000-NASM2019-01747_screen, 100ns
2020/07/14 09:04:40 Time to wait to process https://ids.si.edu/ids/download?id=NASM-A19350058000-NASM2019-01748_screen, 109ns
2020/07/14 09:04:41 exif: failed to find exif intro marker
2020/07/14 09:04:42 Time to complete processing for http://ids.si.edu/ids/deliveryService?id=NASM-A19670206000_PS01, 1.446708455s
2020/07/14 09:04:42 Time to wait to process https://ids.si.edu/ids/download?id=NASM-A19350058000-NASM2019-01752_screen, 1.446648069s
2020/07/14 09:04:44 exif: failed to find exif intro marker
2020/07/14 09:04:44 exif: failed to find exif intro marker
2020/07/14 09:04:44 Time to complete processing for https://ids.si.edu/ids/download?id=NASM-A19350058000-NASM2019-01747_screen, 4.037780245s
2020/07/14 09:04:44 Time to wait to process https://ids.si.edu/ids/download?id=NASM-A19350058000-NASM2019-01758_screen, 2.591255228s
2020/07/14 09:04:44 Time to complete processing for https://ids.si.edu/ids/download?id=NASM-A19350058000-NASM2019-01744_screen, 4.184583514s
2020/07/14 09:04:44 Time to wait to process https://ids.si.edu/ids/download?id=NASM-A19350058000-NASM2019-01760_screen, 146.644379ms
...and so on
```

### wunderkammer-db

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
   	-database-dsn sql:///usr/local/go-wunderkammer/metmuseum.db

$> sqlite3 metmuseum.db 
SQLite version 3.32.1 2020-05-25 16:19:56
Enter ".help" for usage hints.
sqlite> SELECT COUNT(url) FROM oembed;
236288
```

## See also

* https://github.com/aaronland/ios-wunderkammer