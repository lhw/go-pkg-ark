# go-pkg-ark

`ark` is a Go package which allows interaction with the RCON Port of an ARK Survival Gameserver.
It provides all the currently available RCON commands as Go functions for ease of use

Tested only on Linux so far.

## Usage
In order to use the library:

```golang
go get github.com/lhw/go-pkg-ark/arkrcon
```

### Example

```go
package main

import (
  "github.com/lhw/go-pkg-ark/arkrcon"
  "fmt"
)


func main() {
  ark, err := arkrcon.NewARKRconConnection("127.0.0.1:27020", "adminPassword")
  if err != nil {
    fmt.Println(err)
    return
  }

  playerList, err := ark.ListPlayers()
  if err != nil {
    fmt.Println(err)
    return
  }
  if len(playerList) > 0 {
	  for _, pl := range playerList {
		 fmt.Println(pl.Username)
	  }
	  ark.Broadcast("'sup players")
  } else {
	  fmt.Println("No Players Online")
  }
}
```

## API Documentation

Available at https://godoc.org/github.com/lhw/go-pkg-ark/arkrcon

## License

See [LICENSE](../master/LICENSE)
