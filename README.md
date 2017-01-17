# down
Multi goroutine downloader.

## To download:

```
go install github.com/vyasgiridhar/down
```

## Usage:

```
down "http://google.co.in" -g 8 -o page.html
```

```
-g : Number of simultaneous threads to use
-o : Output file name.
```
