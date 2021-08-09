# Track

Time tracking tool written in Go.

## Install

### Building from source

To build from source, you need a working Go environment installed.

You can use the `go` command to install the `track` binary into your `GOPATH`:

```sh
$ go get github.com/jasonwoodland/track/cmd/...
```

### ZSH Completion

For ZSH completion, you need to copy `completion/_track` somewhere into your `$fpath`.

```sh
curl "https://raw.githubusercontent.com/jasonwoodland/track/main/completion/_track" > /opt/homebrew/share/zsh/site-functions/_track
```

## Todo

- [x] show totals for tasks when start/stop/status (add all frames for a total)
- [x] add 'frame add' command
- [x] add confirmations for delete
- [x] report `--csv` output
- [x] sql migrations
- [ ] add `frame move project task frame new_project new_task` command
- [ ] add `add`/`sub` commands to alter the start time of the running task
- [ ] refactor: normalize output/logging, create convenience functions for printProject, printTask, printFrame
- [ ] [timeline] if a frame spans over two dates, it is not included
- [ ] t task merge acs-api feature/some-typo acs-api feature/23/correct-task
