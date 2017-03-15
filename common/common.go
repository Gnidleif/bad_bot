package common

import (
	"bad_bot/invoker"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type TwitterBot struct {
	Name   string
	Script string
}

var (
	BotID     string
	ScriptDir string
	LastTweet time.Time
)

func Start(token, dir string) error {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	u, err := dg.User("@me")
	if err != nil {
		return err
	}

	BotID = u.ID

	dg.AddHandler(messageCreate)

	if err = dg.Open(); err != nil {
		return err
	}

	var bots [2]*TwitterBot
	bots[0] = &TwitterBot{
		Name:   "Charlie Chill",
		Script: "pinkmonkey",
	}
	bots[1] = &TwitterBot{
		Name:   "Stefan Sund",
		Script: "sverjeven",
	}

	for _, bot := range bots {
		go bot.AutoTweet()
	}

	ScriptDir = dir
	return nil
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	split := strings.Split(m.Content, " ")
	if len(split) < 1 || len(regexp.MustCompile("^!.*").FindString(split[0])) == 0 {
		return
	}

	switch script := split[0]; script {
	case "!help":
		s.ChannelMessageSend(m.ChannelID, helpMessage())

	case "!calc":
		err := sendScriptOutput(s, m, script[1:], split[1:]...)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	case "!sverjeven":
		args := strings.Join(split[1:], " ")
		err := sendScriptOutput(s, m, "sverje_ven", args)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	case "!proverb":
		args := strings.Join(split[1:], " ")
		err := sendScriptOutput(s, m, script[1:], args)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	case "!argue":
		err := sendScriptOutput(s, m, "argue", split[1:]...)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	case "!magmys":
		s.ChannelMessageSend(m.ChannelID, magmysMessage())

	case "!spellcheck":
		args := strings.Join(split[2:], " ")
		err := sendScriptOutput(s, m, "spellcheck", split[1], args)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	case "!tweet":
		if !LastTweet.IsZero() {
			if time.Since(LastTweet).Minutes() < 10.0 {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("`Time since last tweet: %s. Cooldown 10 minutes.`", time.Since(LastTweet)))
				return
			}
		}

		var name string
		var script string
		switch split[1] {
		case "pinkmonkey":
			name = "Charlie Chill"
			script = "pinkmonkey"
		case "sverjeven":
			name = "Stefan Sund"
			script = "sverjeven"

		default:
			s.ChannelMessage(m.ChannelID, fmt.Sprintf("`Invalid personality: %s.`", split[1]))
			return
		}

		out, err := invoker.Invoke(ScriptDir, "argue", false, script, "1")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		LastTweet = time.Now()
		msg := fmt.Sprintf("@%s %s", split[2], out)
		reply := "-1"
		if len(split) > 3 {
			reply = split[3]
		}
		err = sendScriptOutput(s, m, "bad_tweet", name, msg, reply)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	default:
		s.ChannelMessageSend(m.ChannelID, helpMessage())
		return
	}
}

func sendScriptOutput(s *discordgo.Session, m *discordgo.MessageCreate, script string, args ...string) error {
	out, err := invoker.Invoke(ScriptDir, script, true, args...)
	if len(out) > 0 {
		fmt.Println(out)
	}
	if err != nil {
		return err
	}
	s.ChannelMessageSend(m.ChannelID, out)
	return nil
}

func helpMessage() string {
	return fmt.Sprintf(`
Commands:
!calc <equation> - Calculates <equation>
!sverjeven <text> - Boosts your patriot-level
!magmys - Receive important rules
!proverb <amount> - Receive wisdom
!argue <personality> <amount> - Receive DID-like powers
!spellcheck <percent> <text> - Stop being dyslexic
!tweet <handle> - Tweet something smart to <handle>`)
}

func magmysMessage() string {
	return fmt.Sprintf(`
Här i gruppen finns några regler som ska och måste efterföljas.
1. Admin har förtur
2. Magmys/magbox är tillåtet vid admins godkännande
3. Snopp på mage är tillåtet
4. Trevliga mot varandra
5. Blockering av admin/mods ej tillåtet.
6. Respektera alla människor
7. Man får inte vara kräsen!
8. Inget bagbang ( 2st tjejer och en man ) Här inte INTE OKEJ!
9. Ingen piercing i naveln
10. Fördriv inte tid utan bjud på er MAGMYS är MÅLET med nya gruppen :)`)
}

func (bot *TwitterBot) AutoTweet() {
	for {
		out, err := invoker.Invoke(ScriptDir, "argue", false, bot.Script, "1")
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if len(out) <= 5 {
			continue
		}

		if bot.Name == "Stefan Sund" {
			out = fmt.Sprintf("@MrHenko123 %s #svpol", out)
		} else {
			out = fmt.Sprintf("%s #svpol", out)
		}

		_, err = invoker.Invoke(ScriptDir, "bad_tweet", false, bot.Name, out, "-1")
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		time.Sleep(20 * time.Minute)
	}
}
