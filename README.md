# callow

## Testing

`while true; do inotifywait -e close_write -r . 2>/dev/null ; make test; clear; ./test; done`
