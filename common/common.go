package common

import (
	"bad_bot/invoker"
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	BotID     string
	ScriptDir string
	validCmd  = regexp.MustCompile(`^!(calc|sverjeven|proverb|argue|magmys|spellcheck|spongebob|help)$`)
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

	ScriptDir = dir
	return nil
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	split := strings.Split(m.Content, " ")
	if len(split) < 1 || len(validCmd.FindString(split[0])) == 0 {
		return
	}

	switch script := split[0][1:]; script {
	case "calc":
		err := sendScriptOutput(s, m, script[1:], split[1:]...)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	case "sverjeven":
		args := strings.Join(split[1:], " ")
		err := sendScriptOutput(s, m, "sverjeven", args)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	case "proverb":
		args := strings.Join(split[1:], " ")
		err := sendScriptOutput(s, m, script[1:], args)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	case "argue":
		err := sendScriptOutput(s, m, "argue", split[1:]...)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	case "magmys":
		s.ChannelMessageSend(m.ChannelID, magmysMessage())

	case "spellcheck":
		args := strings.Join(split[2:], " ")
		err := sendScriptOutput(s, m, "spellcheck", split[1], args)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	case "spongebob":
		args := strings.Join(split[1:], " ")
		err := sendScriptOutput(s, m, "spongebob", args)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	case "help":
		fallthrough
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
		s.ChannelMessageSend(m.ChannelID, helpMessage())
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
!spongebob <text> - tHe LeFt CaN't MeMe`)
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
