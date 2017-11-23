package main

import (
          "fmt"
          "os"
          "os/signal"
          "strings"
          "syscall"
          "strconv"

          "github.com/bwmarrin/discordgo"	
        )

func main(){
	discord, err := discordgo.New("Bot MzgyNjk5NzE4NzEyNjIzMTMz.DPePIQ.Tuk8zqYLBmtKFSqJyPk1Cv9mAYQ")

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
  var pieces = strings.Split(msg.Content, " ")
  var command = pieces[0]
  switch command {
    case "*color":
      processColorCommand(pieces, s, msg)
    default:
      {
        if strings.HasPrefix(command, "*") {
          processHelpCommand(s,msg)
        }
      }
  }
}

func processColorCommand(pieces []string, s *discordgo.Session, msg *discordgo.Message) {
  if len(pieces) == 2 {
    color, success := parseColor(pieces[1])
    if success {
      // TODO: There should be an easier way to get the guild id
      guildID := getGuildID(s, msg.ChannelID)
      userID := msg.Author.ID

      setNameColor(s, guildID, userID, color)
    } else {
      fmt.Println("failed to parse color")
    }
  } else {
    fmt.Println("failed to parse color")
  }
}

func processHelpCommand(s *discordgo.Session, msg *discordgo.Message) {
  var messageContent = "Please enter a command of the following format:\n\t***color <color_value>**\n\tWhere color value is of the format *#<hex_value>*, *0x<hex_value>*, or *<decimal_value>*\n\tYou can find the corresponding hex values for colors here: https://www.w3schools.com/colors/colors_picker.asp"
  s.ChannelMessageSend(msg.ChannelID, messageContent)
}

func getGuildID(s *discordgo.Session, channelID string) string {
  channel, err := s.Channel(channelID)
  if err != nil {
    fmt.Println("error retrieving channel to get guild ID", err)
    return ""
  }
  return channel.GuildID
}

func setNameColor(s *discordgo.Session, guildID string, userID string, color int) {
  existingRoles, _ := s.GuildRoles(guildID)
  var existingRole *discordgo.Role
  for _, element := range existingRoles {
    if element.Name == userID + "'s color role" {
      existingRole = element
    }
  }
  if existingRole != nil {
    s.GuildRoleEdit(guildID, existingRole.ID, existingRole.Name, color, existingRole.Hoist, existingRole.Permissions, existingRole.Mentionable)
  } else {
    newRole, _ := s.GuildRoleCreate(guildID)
    s.GuildRoleEdit(guildID, newRole.ID, userID + "'s color role", 3447003, false, 0, false)
    s.GuildMemberRoleAdd(guildID, userID, newRole.ID)
  }
}

func parseColor(color string) (int, bool) {
  var decimalValue int
  if strings.HasPrefix(color, "#") {
    decimalValue64, err := strconv.ParseInt(color[1:len(color)], 16, 32)
    decimalValue = int(decimalValue64)
    if err != nil {
      decimalValue = -1
    }
  } else if strings.HasPrefix(color, "0x") {
    decimalValue64, err := strconv.ParseInt(color, 0, 32)
    decimalValue = int(decimalValue64)
    if err != nil {
      decimalValue = -1
    }
  } else {
    decimalValue64, err := strconv.ParseInt(color, 10, 32)
    decimalValue = int(decimalValue64)
    if err != nil {
      decimalValue = -1
    }
  }
  
  if decimalValue >= 0 && decimalValue <= 16777215 {
    return decimalValue, true
  } else {
    return -1, false
  }
}