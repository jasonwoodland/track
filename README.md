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

# Todo

- [x] show totals for tasks when start/stop/status (add all frames for a total)
- [x] add 'frame add' command
- [x] add confirmations for delete
- [ ] log `--csv` output
- [ ] sql migrations
- [ ] desktop notifications for running tasks
- [ ] add complete command
        add complete_at date column on task
        usage: t complete [--reset] project task
          sets the complete_at field for task
        task subcommands, start, will check if complete_at is set, and show a
        confirmation prompt before continuing.
        log|timeline -c 02 will show tasks completed in february.
