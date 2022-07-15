package bot

import (
	"fmt"                    //to print errors
	"shitty-steambot/config" //importing our config package which we have created above
	"shitty-steambot/steam"  //importing our steam package which we have created above
	"strings"
	"time"

	"github.com/bwmarrin/discordgo" //discordgo package from the repo of bwmarrin .
)

var BotId string
var goBot *discordgo.Session

func Start() {

	//creating new bot session
	goBot, err := discordgo.New("Bot " + config.Token)

	//Handling error
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Making our bot a user using User function .
	u, err := goBot.User("@me")
	//Handlinf error
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Storing our id from u to BotId .
	BotId = u.ID

	// Adding handler function to handle our messages using AddHandler from discordgo package. We will declare messageHandler function later.
	goBot.AddHandler(shittySteamBot)

	err = goBot.Open()
	//Error handling
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//If every thing works fine we will be printing this.
	fmt.Println("Bot is running !")
}

func shittySteamBot(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotId {
		return
	}
	if m.Content == config.BotPrefix+"ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
		return
	}
	if m.Content == config.BotPrefix+"help" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "	**Commands :**\n\n"+
			config.BotPrefix+"help : Shows this message.\n"+
			config.BotPrefix+"ping : Pong.\n"+
			config.BotPrefix+"steam : Gives random game on Steam.\n"+
			config.BotPrefix+"steam <game name> : Gives the games matching the given name on Steam.\n"+
			config.BotPrefix+"steam genre <genre> : Gives a random game matching the given genre Steam. (Not working cause it's reserved to partner, but I don't care)\n",
		)
		return
	}
	// Searching for a random game
	if m.Content == config.BotPrefix+"steam" {
		start := time.Now()
		message, length := steam.GetRandomGame()
		elapsed := time.Since(start)
		if length > 0 {
			_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I found this in a list of %d in %s : %s\n", length, elapsed, message))
		} else {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Sorry, there was a problem.\n")
		}
		return
	}
	// Searching for a game by genre
	if strings.HasPrefix(m.Content, config.BotPrefix+"steam genre ") {
		start := time.Now()
		arrayOfGenres := strings.Split(m.Content[len(config.BotPrefix+"steam genre "):], " ")
		fmt.Print(arrayOfGenres)
		message, length := steam.GetGameByGenre(arrayOfGenres)
		elapsed := time.Since(start)
		if strings.HasPrefix(message, "https://") {
			_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I found this in a list of %d in %s : %s\n", length, elapsed, message)+"\n")
		} else {
			_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Sorry, nothing was found in a list of %d.\n", length))
		}
		return
	}
	// Searching for a game by name
	if strings.HasPrefix(m.Content, config.BotPrefix+"steam ") {
		start := time.Now()
		message, length := steam.GetGameByName(m.Content[len(config.BotPrefix+"steam "):])
		elapsed := time.Since(start)
		if strings.HasPrefix(message, "https://") {
			_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I found this in a list of %d in %s : %s\n", length, elapsed, message)+"\n")
		} else {
			_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Sorry, nothing was found in a list of %d.\n", length))
		}
		return
	}
}
