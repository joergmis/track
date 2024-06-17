# track

Timetracking application with focus on:

- [local-first](https://www.inkandswitch.com/local-first/)
- adding time entries
- visualisation of hours

The goal is to have a generic timetracking application, which can be extended
to sync timeentries or pull data like customers or projects from different 
services (starting with clockodo).

## Setup and installation

Copy the configuration file and adjust the values.

```bash
cp track.example.yaml ~/.config/track/track.yaml
```

Install the binary; if currently tries to retrieve customers, projects and 
services from clockodo which means you have to create a config with valid 
credentials (see above).

```bash
make install
```

Check if app has been installed; if the app is not found, check if `~/go/bin` 
is in `$PATH`.

```bash
track version
```

To make it easier to use, set up autocompletion:

```bash
track completion bash > /tmp/completion
source /tmp/completion
```

## Playground for different tools and approaches

- combinatory [approval testing](https://github.com/approvals/go-approval-tests)
- [mutation testing](https://github.com/avito-tech/go-mutesting)

## References

Loosely inspired by:

- [timewarrior](https://github.com/GothenburgBitFactory/timewarrior)
- [goddd](https://github.com/marcusolsson/goddd)
