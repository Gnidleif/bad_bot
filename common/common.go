package common

import (
	"bad_bot/invoker"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	BotID     string
	ScriptDir string
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
	if len(split) < 2 || split[0] != "badbot" {
		return
	}
	switch script := split[1]; script {
	case "help":
		_, _ = s.ChannelMessageSend(m.ChannelID, helpMessage())

	case "calc":
		out, err := invoker.Invoke(ScriptDir, script, split[1:]...)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		_, _ = s.ChannelMessageSend(m.ChannelID, out)

	case "sverje_ven":
		args := strings.Join(split[2:], " ")
		out, err := invoker.Invoke(ScriptDir, script, args)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		_, _ = s.ChannelMessageSend(m.ChannelID, out)

	case "magmys":
		_, _ = s.ChannelMessageSend(m.ChannelID, magmysMessage())

	default:
		_, _ = s.ChannelMessageSend(m.ChannelID, helpMessage())
		return
	}
}

func helpMessage() string {
	return fmt.Sprintf(`
Commands:
calc <equation> - Calculates <equation>
sverje_ven <text> - Boosts your patriot-level
magmys - Receive important rules`)
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
