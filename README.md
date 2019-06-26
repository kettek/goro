# Going Rogue / goRo
This library is for writing roguelikes in Go. Please note that it is highly experimental at the moment.

## Building
goRo can have the following tags supplied during building:

  * disableTCell
  * enableEbiten

If you wish to use the graphical interface only:

```
# go build -tags "disableTCell enableEbiten"
```

If you wish to use the terminal interface only:

```
# go build
```

If you wish to have support for both and choose one programmatically:

```
# go build -tags "enableEbiten"
```
