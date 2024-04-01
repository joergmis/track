# track

Timetracking application with primary focus on adding time-entries and 
visualising them.

## Setup

* copy `track.example.yaml` to `~/.config/track.yaml` and adjust the values
* install app by running `make install`
* ensure `~/go/bin` is in your `$PATH`
* optional; setup bash completion
  * `track completion bash > /tmp/completion`
  * `source /tmp/completion`

## References

Loosely inspired by:

- [gotimetrack](https://github.com/danielbatw/gotimetrack)
- [timewarrior](https://github.com/GothenburgBitFactory/timewarrior)
- [goddd](https://github.com/marcusolsson/goddd)
