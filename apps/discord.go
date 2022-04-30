package apps

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func Discord() {
	token := os.Getenv("DISCORD_TOKEN")
	addBotUrl := os.Getenv("DISCORD_BOT_URL")

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.Identify.Intents = discordgo.IntentGuildMembers | discordgo.IntentsGuildMessages | discordgo.IntentsGuilds

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	fmt.Println("Add bot to your server via this link: " + addBotUrl)

	dg.AddHandler(bogJoined)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func bogJoined(s *discordgo.Session, m *discordgo.GuildCreate) {
	members, _ := s.GuildMembers(m.Guild.ID, "", 100)
	reader := bufio.NewReader(os.Stdin)

	var selectedMember discordgo.Member

	fmt.Println("Discord bot is added to your server. \nSelect which user you want to write")

	for i, member := range members {
		if member.User.ID != s.State.User.ID {
			fmt.Printf("%d: %s\n", i+1, member.User.Username)
		}
	}

	for {
		fmt.Print("Type the number: ")
		char, _, err := reader.ReadRune()

		if err != nil {
			panic(err)
		}

		number := int(char - '0')

		if number >= len(members) {
			fmt.Println("Please, provide a valid number")
		} else {
			selectedMember = *members[number-1]
			break
		}
	}

	channel, err := s.UserChannelCreate(selectedMember.User.ID)

	if err != nil {
		panic(err)
	}

	message, err := reader.ReadString('\n')

	for {
		fmt.Println("Type message you want to send to this user: ")
		message, err = reader.ReadString('\n')

		if err != nil {
			panic(err)
		}

		_, err = s.ChannelMessageSend(channel.ID, message)

		if err != nil {
			panic(err)
		} else {
			fmt.Println("Message sent successfully")
		}
	}
}
