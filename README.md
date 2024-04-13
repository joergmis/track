# track

Timetracking application with focus on:

- [local-first](https://www.inkandswitch.com/local-first/)
- adding time entries
- visualisation of hours

## Setup and installation

Copy the configuration file and adjust the values.

```bash
cp track.example.yaml ~/.config/track.yaml
```

Generate the list of projects/clients/services which is required for the 
autocompletion to work.

```bash
go run cmd/track/main.go generate
```

Now you are ready to install the app:

```bash
go install cmd/track
```

Check if app has been installed; if the app is not found, check if `~/go/bin` 
is in $PATH.

```bash
track version
```

To make it easier to use, set up autocompletion:

```bash
track completion bash > /tmp/completion
source /tmp/completion
```

## References

Loosely inspired by:

- [timewarrior](https://github.com/GothenburgBitFactory/timewarrior)
- [goddd](https://github.com/marcusolsson/goddd)
