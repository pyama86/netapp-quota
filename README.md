# netapp-quota

Apply NetApp quarter at regular intervals.

## Description
```
Usage of netapp-quota:
  -off-interval int
        quota off interval (default 300)
  -on-interval int
        quota on interval (default 10)
  -password string
        netapp api BasicAuthPassword
  -prefix string
        netapp volume prefix
  -svm string
        netapp svm server name
  -url string
        netapp api endpoint
  -user string
        netapp api BasicAuthUser
  -version
        Print version information and quit.
```


## Install

To install, use `go get`:

```bash
$ go get -d github.com/pyama86/netapp-quota
```

## Contribution

1. Fork ([https://github.com/pyama86/netapp-quota/fork](https://github.com/pyama86/netapp-quota/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[pyama86](https://github.com/pyama86)
