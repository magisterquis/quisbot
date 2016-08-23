package main

/*
 * mode.go
 * Handle mode changes
 * By MagisterQuis
 * Created 20160821
 * Last Modified 20160821
 */

import (
	"fmt"
	"log"
	"strings"
)

/* HandleMode handles IRC mode changes.  Specifically it looks for ops */
func HandleMode(src, op, args string) error {
	/* Split the message apart */
	parts := strings.Fields(args)
	if 3 > len(parts) {
		return fmt.Errorf("not enough parts in MODE message %v", args)
	}
	ch := parts[0]
	mode := parts[1]
	nicks := parts[2:]

	/* Take action for each nick changed */
	for _, nick := range nicks {
		go LogFirstLast(nick, op, ch, fmt.Sprintf(
			"%v by %v",
			mode,
			src,
		))
		/* It seems to be only ops/deops */
		switch mode {
		case "-o": /* DeOp */
			if err := SetChanOp(nick, ch, false); nil != err {
				return err
			}
			log.Printf(
				"[MODE] %v deopped by %v in %v",
				nick,
				src,
				ch)
		case "+o": /* Op */
			if err := SetChanOp(nick, ch, true); nil != err {
				return err
			}
			log.Printf("[MODE] %v opped by %v in %v",
				nick,
				src,
				ch,
			)
		default:
			log.Printf(
				"[MODE] Unknown mode %v for %v by %v in %v",
				mode,
				nick,
				src,
				ch,
			)
		}
	}
	return nil
}
