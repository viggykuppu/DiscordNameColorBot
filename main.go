package main

import (
          "fmt"
          "io/ioutil"
          "os"
          "os/signal"
          "strings"
          "syscall"
          "strconv"

          "github.com/bwmarrin/discordgo"	
        )
var authToken string

func main() {
  parseConfig()

	discord, err := discordgo.New("Bot " + authToken)

  discord.AddHandler(messageCreate)

  // Open the websocket and begin listening.
  err = discord.Open()

  if err != nil {
    fmt.Println("Error opening Discord session: ", err)
  } else {
    err = discord.UpdateStatus(0,"Type *help")

    // Wait here until CTRL-C or other term signal is received.
    fmt.Println("Discord name color bot.  Press CTRL-C to exit.")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

    // Cleanly close down the Discord session.
    fmt.Println("Closing bot... goodbye!")
    discord.Close()
  }
}

func parseConfig(){
  defer handleInternalError()
  dat, err := ioutil.ReadFile("settings")
  check(err)

  var stringData = string(dat)
  var lines = strings.Split(stringData,"\n")

  for _, line := range lines {
    if !strings.HasPrefix(line, "#") {
      pieces := strings.Split(line,"=")
      if len(pieces) == 2 {
        key := pieces[0]
        value := pieces[1]
        switch key {
          case "auth_token":
            authToken = value
        }
      } 
    }
  }
}

func check(err error) {
  if err != nil {
    panic(err)
  }
}

func handleInternalError(){
  if r := recover(); r != nil {
    fmt.Println("Recovered in f", r)
  }
}

// Event to handle message creation on any channels that the bot has access to
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
  var msg = m.Message
  var pieces = strings.Split(msg.Content, " ")
  var command = pieces[0]
  switch command {
    case "*color":
      processColorCommand(pieces, s, msg)
    case "*help":
      processHelpCommand(s,msg)
  }
}

func processColorCommand(pieces []string, s *discordgo.Session, msg *discordgo.Message) {
  defer handleError()
  if len(pieces) == 2 {
    color := parseColor(pieces[1])
    // TODO: There should be an easier way to get the guild id
    guildID := getGuildID(s, msg.ChannelID)
    userID := msg.Author.ID

    setNameColor(s, guildID, userID, color)

    // directChannel, _ := s.UserChannelCreate(userID)
    // var messageContent = msg.Author.Username + ", your color has now been set to " + pieces[1]
    // s.ChannelMessageSend(directChannel.ID, messageContent)
  } else {
    fmt.Println("failed to parse color")
  }
}

func processHelpCommand(s *discordgo.Session, msg *discordgo.Message) {
  defer handleError()
  directChannel, err := s.UserChannelCreate(msg.Author.ID)
  check(err)
  var messageContent = "Please enter a command of the following format:\n\t***color <color_value>**\n\tWhere color value is of the format *#<hex_value>*, *0x<hex_value>*, or *<decimal_value>*\n\tYou can find the corresponding hex values for colors here: https://www.w3schools.com/colors/colors_picker.asp"
  s.ChannelMessageSend(directChannel.ID, messageContent)
}

func getGuildID(s *discordgo.Session, channelID string) string {
  channel, err := s.Channel(channelID)
  check(err)
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
    newRole, err := s.GuildRoleCreate(guildID)
    check(err)
    s.GuildRoleEdit(guildID, newRole.ID, userID + "'s color role", color, false, 0, false)
    s.GuildMemberRoleAdd(guildID, userID, newRole.ID)
  }
}

func parseColor(color string) int {
  var decimalValue int
  if strings.HasPrefix(color, "#") {
    decimalValue64, err := strconv.ParseInt(color[1:len(color)], 16, 32)
    decimalValue = int(decimalValue64)
    check(err)
  } else if strings.HasPrefix(color, "0x") {
    decimalValue64, err := strconv.ParseInt(color, 0, 32)
    decimalValue = int(decimalValue64)
    check(err)
  } else {
    decimalValue64, err := strconv.ParseInt(color, 10, 32)
    decimalValue = int(decimalValue64)
    check(err)
  }
  
  if decimalValue >= 0 && decimalValue <= 16777215 {
    return decimalValue
  } else {
    panic(-1)
  }
}

func handleError(){
  if r := recover(); r != nil {
    fmt.Println("Error", r)
  }
}