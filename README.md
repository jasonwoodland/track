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
curl "https://raw.githubusercontent.com/jasonwoodland/track/main/completion/_track" > ~/.zsh/completion/_track
```

## Todo

- [x] show totals for tasks when start/stop/status (add all frames for a total)
- [x] 'add' command
- [x] add confirmations for delete
- [x] report `--csv` output
- [x] sql migrations
- [x] add `frame move project task frame new_project new_task` command
- [x] add `shift` command to alter the start time of the running task eg `t shift -5m` advances the start time by 5 minutes
- [x] refactor: normalize output/logging
- [x] add `task set --repeat project task`
- [ ] add `last` command which would allow us to adjust the last inserted frame (synonymous for: `t frame edit [command options] <last_project> <last_task> <last_frame>`)
- [ ] refactor: create convenience functions for printProject, printTask, printFrame
- [ ] fix timeline: if a frame spans over two dates, it is not included (just print based on the end_time)
