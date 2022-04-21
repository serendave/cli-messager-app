package apps

import (
	"fmt"
	"os"
)

func Discord() {
	token := os.Getenv("DISCORD_TOKEN")
	fmt.Println(token)
}
