package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	archeage "github.com/geeksbaek/archeage-go"
)

// Variables used for command line parameters
var (
	Token string
	BotID string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Get the account information.
	u, err := dg.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}

	// Store the account ID for later use.
	BotID = u.ID

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandlerOnce(cronTask)
	dg.AddHandler(auctionMessage)
	dg.AddHandler(charactorMessage)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func cronTask(s *discordgo.Session, m *discordgo.MessageCreate) {
	aa := archeage.New(&http.Client{})
	var oldNotices archeage.Notices
	ticker := time.Tick(time.Second * 10)
	for _ = range ticker {
		fmt.Print(".")
		newNotices, err := aa.FetchNotice()
		if err != nil || len(newNotices) == 0 {
			log.Println(err)
			continue
		}
		diffNotices := oldNotices.Diff(newNotices)
		if len(diffNotices) > 0 && len(oldNotices) > 0 {
			for _, notice := range diffNotices {
				msg := fmt.Sprintf("[%s] %s %s", notice.Category, notice.Title, notice.URL)
				s.ChannelMessageSend(m.ChannelID, msg)
				log.Println("\n" + msg)
			}
		}
		oldNotices = newNotices
	}
}

func auctionMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	if strings.HasPrefix(m.Content, "?경매장") {
		aa := archeage.New(&http.Client{})
		itemAndQuantity := strings.Split(m.Content, "*")
		var item string
		var quantity int
		if len(itemAndQuantity) == 2 {
			item = strings.TrimSpace(itemAndQuantity[0])
			quantity, _ = strconv.Atoi(strings.TrimSpace(itemAndQuantity[1]))
		} else {
			item = strings.TrimSpace(m.Content)
			quantity = 1
		}

		keyword := func() string {
			fields := strings.Fields(item)
			if len(fields) > 1 {
				return strings.Join(fields[1:], " ")
			}
			return ""
		}()

		result, err := aa.Auction("TOTAL", keyword, quantity)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			log.Println(err)
			return
		}
		lack, price := result.Price(quantity)
		var msg string
		if lack {
			msg = "경매장에 아이템이 모자랍니다."
		} else {
			msg = fmt.Sprintf("`%v` x `%v` = `%v`", result[0].Name, quantity, price.String())
		}
		s.ChannelMessageSend(m.ChannelID, msg)

	}
}

func charactorMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	if strings.HasPrefix(m.Content, "?캐릭터") {
		aa := archeage.New(&http.Client{})
		m.Content = strings.TrimLeft(m.Content, "?캐릭터 ")
		var cs archeage.Characters
		var err error
		if strings.Contains(m.Content, "@") {
			splited := strings.Split(m.Content, "@")
			cs, err = aa.SearchCharactor(archeage.ServerNameMap[strings.TrimSpace(splited[1])], strings.TrimSpace(splited[0]))
		} else {
			cs, err = aa.SearchCharactor("", strings.TrimSpace(m.Content))
		}
		if err != nil {
			return
		}
		for _, c := range cs {
			if c != nil {
				s.ChannelMessageSend(m.ChannelID, c.String())
			}
		}
	}
}
