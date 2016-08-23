package main

/*
 * join.go
 * Handle JOIN messages
 * By MagisterQuis
 * Created 20160821
 * Last Modified 20160821
 */

/* hails is a list of users to hail on joins */
var hails = map[string]struct{}{
	"mairielpis":   struct{}{},
	"imperatorpat": struct{}{},
	"h46u":         struct{}{},
}

/* handleJoin handles JOIN messages */
func HandleJoin(nick, op, channel string) error {
	/* Log the first and latest join */
	go LogFirstLast(nick, "JOIN", channel, channel)

	/* Hail users */
	if _, ok := hails[nick]; ok {
		if err := Privmsg(channel, "All hail "+nick+"!"); nil != err {
			return err
		}
	}

	/* Note a user as in the channel */

	return nil
}
