# track

A time tracking CLI tool written in Go.

# Installation

```sh
go get github.com/jasonwoodland.com/track
cd $GOPATH/src/github.com/jasonwoodland.com/track
GOBIN="/usr/local/bin" go install

```

## ZSH Completion

For ZSH completion, you need to copy `completion/_track` into your `$FPATH`.
```sh
cp $GOPATH/src/github.com/jasonwoodland.com/track/completion/_track /usr/local/share/zsh/site-functions/_track
```

# TODO

- [x] show totals for tasks when start/stop/status (add all frames for a total)
- [ ] add 'add' command
- [x] add confirmations for delete
- [ ] add tags for tasks [closed]
- [ ] log `--csv` output
- [ ] sql migrations
- [ ] desktop notifications for running tasks
