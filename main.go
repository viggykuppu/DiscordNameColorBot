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
	discord, err := discordgo.New("MzgyNjk5NzE4NzEyNjIzMTMz.DPZhbg.8SB4ux0S99p4FL9kiPdkWoUUdiU")

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
      _, err := s.ChannelMessageSend(msg.ChannelID, msg.Content)
      if err != nil {
        fmt.Println("Error sending messages to discord: ", err)
      }
    }
  }
}