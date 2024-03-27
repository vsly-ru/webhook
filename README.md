# Go Webhook
A very simple webhook app written in go(golang). It listens for HTTP requests with a specified URL and when one occurs, runs a specified command (typically a bash script). `stdout` + `stderr` will be returned in the response.

## Arguments

- `-u` **REQUIRED** The **U**RL to trigger on. All other URLs will return 404. Starts with "/". Example: `-u /webhook/2ff80e9159b517704ce43f0f74e6e247`
- `-p` **P**ort to listen. Default is `7999`
- `-c` **C**ommand to execute when the webhook URL is invoked. Default is `./script.sh`. To provide arguments, put the command in single quotes. Example: `-c 'ls -a -l -h'`
- `-m` HTTP **m**ethod of the webhook URL. Default is `GET`.
- `-w` **W**ait duration in seconds between the command exuctions. Default is `10`.

## Security
- First, you need to specify a long and **unique** webhook URL (with `-u`), so it cannot be guessed.
- It's recommended to use HTTPS (to hide the exact URL)
- By default, only one command invocation every 10 seconds is allowed. This cooldown can be configured with the `-w` flag. If a new request arrives while the previous command is still running or within the 10-second cooldown period after its completion, the application will respond with the HTTP error "Too Many Requests" (status code 429).
- If you run this application as root, it will also run the command as root.

## Building

Build for current platform
```bash
go build .
```

Build for linux amd64
```bash
env GOOS=linux GOARCH=amd64 go build -o webhook_linux_x64
```