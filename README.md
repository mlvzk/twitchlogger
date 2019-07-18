# twitch\_logger

Lightweight and minimal Twitch chat logger.
It is able to log every live chat on twitch (at once) with very low CPU and memory usage.

## Commands

There are 2 commands, `twitch-log` and `twitch-channels`:
- `twitch-channels` is for fetching all live channels from the Twitch API, it outputs CSV data (user\_id,channel\_name,viewers). Check the usage with `twitch-channels -h`
- `twitch-log` is for logging the chats. The input must be a stream with every channel name separated by a newline. It outputs IRC messages.

You can connect `twitch-channels` and `twitch-log` by piping `twitch-channels` to `awk`, as in the examples below.

## Installation

Download the executables from this repo's [Releases page](https://github.com/mlvzk/twitchlogger/releases) and put them in your $PATH

OR

```bash
go get -u github.com/mlvzk/twitchlogger/...
```

## Examples

Pipe `twitch-channels` to `awk` to only get the channel's name and pipe it then to `twitch-log`:
```bash
twitch-channels \
| awk -F, '{ print $2 }' \
| twitch-log
```

Log only one channel:
```bash
twitch-log <<<"moonmoon_ow"
```

Fetch only english(locale id `en`) streams with no less than `50` viewers, and loop it to also fetch future streams:
```bash
twitch-channels -language en -min 50 -loop
```

Fetch, transform, log, filter to only PRIVMSG messages and beautify the output:
```bash
twitch-channels -language en -min 20 -loop \
| awk -F, '{ print $2 }' \
| twitch-log \
| awk 'match($0, /:(\w+)!.*PRIVMSG #(\w+) :(.*)/, m) { print strftime("%Y-%m-%d %H:%M:%S")" \033[34m"m[1]"\033[0m to \033[35m#"m[2]"\033[0m: "m[3] }' # only gnu version of awk, gawk
```
