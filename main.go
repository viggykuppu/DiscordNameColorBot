package main

import (
          "fmt"
          "os"
          "os/signal"
          "strings"
          "syscall"

          "github.com/bwmarrin/discordgo"	
        )

func main(){
	discord, err := discordgo.New("Bot MzgyNjk5NzE4NzEyNjIzMTMz.DPZhbg.8SB4ux0S99p4FL9kiPdkWoUUdiU")

  if err != nil {
    fmt.Println("Error creating discord session: ", err)
  }

  discord.AddHandler(messageCreate)


  // Open the websocket and begin listening.
  err = discord.Open()

  if err != nil {
    fmt.Println("Error opening Discord session: ", err)
  }

  // Wait here until CTRL-C or other term signal is received.
  fmt.Println("Discord name color bot.  Press CTRL-C to exit.")
  sc := make(chan os.Signal, 1)
  signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
  <-sc

  // Cleanly close down the Discord session.
  fmt.Println("Closing bot... goodbye!")
  discord.Close()
}

func parseConfig(){

}

// Event to handle message creation on any channels that the bot has access to
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
  var msg = m.Message
  if msg.ChannelID == "360063821718487043" {
    if strings.HasPrefix(msg.Content, "*color ") {
      //var color = strings.Split(msg.Content, " ")
      fmt.Println(msg.Content)

      // TODO: There should be an easier way to do this. We should just be able to see which channels the bot is bound to somehow
      guildID := getGuildID(s, msg.ChannelID)

      createBasicRoleWithColor(s, 3447003)

      // _, err := s.channelmessagesend(msg.channelid, msg.content)
      // if err != nil {
      //   fmt.println("error sending messages to discord: ", err)
      // }
    }
  }
}

func getGuildID(s *discordgo.Session, channelID string) {
  channel, err := s.Channel(channelID)
  if err != nil {
    fmt.Println("error retrieving channel to get guild ID", err)
    return nil
  }
  return channel.GuildID
}

func createBasicRoleWithColor(s *discordgo.Session, guildID string, color int) {
  newRole, err := s.GuildRoleCreate(newRole)
  


  newRole := discordgo.Role{Name:"new color role", Managed: false, Mentionable: false, Hoist: false, Managed: false, Color: color, Position: 1, Permissions: 0}
  
}