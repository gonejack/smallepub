# smallepub
This command line tool unpack .epub file and downgrade inside images and repack to make it small.

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gonejack/smallepub)
![Build](https://github.com/gonejack/smallepub/actions/workflows/go.yml/badge.svg)
[![GitHub license](https://img.shields.io/github/license/gonejack/smallepub.svg?color=blue)](LICENSE)

### Install
```shell
> go install github.com/gonejack/smallepub@latest
```

### Usage
```shell
> smallepub *.epub
```
```
Flags:
  -h, --help          Show context-sensitive help.
      --about         About.
  -v, --verbose       Verbose printing.
  -q, --quality=40    Picture compress rate/quality (1-100)
```
