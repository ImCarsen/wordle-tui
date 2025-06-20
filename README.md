<h1 align="center">World TUI</h1>

> A wordle TUI built in Go for fun

# Build

```sh
go build -ldflags="-s -w" -o wordle-tui .
```

# Usage

Play (word of the day from nytimes):
```sh
./wordle-tui
```

Args:

| Argument | Use |
|----------|-----|
| --reset  | Resets the save file, allowing you to replay the game |
| --word <word> | Allows you to set a custom word |