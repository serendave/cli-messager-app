package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to cli-messager-app. Please which bot you would like to use:")
	fmt.Println("1: slack")
	fmt.Println("2: discord")
	fmt.Println("q: exit")

	for {
		fmt.Print("Your choice: ")
		char, _, err := reader.ReadRune()

		if err != nil {
			panic(err)
		}

		if char == '1' {
			fmt.Println("Slack is selected")
			slack()
			break
		} else if char == '2' {
			fmt.Println("Discord is selected")
			discord()
			break
		} else if char == 'q' {
			os.Exit(1)
		} else {
			fmt.Println("Please enter a valid choice")
		}
	}
}
