package main

/*
 * config.go
 * Read config updates
 * By MagisterQuis
 * Created 20170625
 * Last Modified 20170913
 */

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

/* config is the struct used for creating and unJSONing the config file */
type config struct {
	Greetings     map[string]string
	Responses     map[string]string
	Bots          []string
	Channels      []string
	GuestChannels []string
}

var (
	// GREETINGS holds nick -> hello
	GREETINGS = map[string]string{}
	// RESPONSES holds canned responses to "commands"
	RESPONSES = map[string]string{}
	// BOTS get no greetings
	BOTS = map[string]struct{}{}
	// CHANNELS is the set of channels to join
	CHANNELS = map[string]struct{}{}
	// GUESTCHANNELS is the channels in which quisbot is quiet
	GUESTCHANNELS = map[string]struct{}{}
	// LOCK locks the above
	LOCK = &sync.Mutex{}

	/* defaultconfig is what's made if there's no config */
	defaultconfig = config{
		Greetings: map[string]string{
			"magisterquis": "All hail the creator!",
			"mairielpis":   "All hail the wife of the creator!",
		},
		Responses: map[string]string{
			"!help": "I am quisbot.  Have a nice day.",
		},
		Bots: []string{
			"quisbot",
			"moobot",
		},
		Channels: []string{"magisterquis"},
		GuestChannels: []string{
			"imperatorpat",
			"mairielpis",
			"yolomcswagatron454",
			"datacrab",
		},
	}
)

/* startConfigPoll polls the config file named fn, making a new one if it
doesn't exist.  If there are changes, it updates the appropriate structs.
w is used to join and part channels.
*/
func startConfigPoll(
	fn string,
	pint time.Duration,
	ready *sync.WaitGroup,
	w io.Writer,
) {
	/* Wait until we're ready */
	ready.Wait()

	/* Last time file was modified */
	var lastModified time.Time

	/* Read or create file */
	if err := doConfigPoll(fn, &lastModified, w); nil != err {
		log.Fatalf(
			"Unable to create or read config file %v: %v",
			fn,
			err,
		)
	}

	/* Start poller */
	for {
		time.Sleep(pint)
		if err := doConfigPoll(fn, &lastModified, w); nil != err {
			log.Printf(
				"Error reading or creating config file %v: %v",
				fn,
				err,
			)
		}
	}
}

/* doConfigPoll polls the config file fn for changes, or makes one if it
doesn't exist */
func doConfigPoll(fn string, last *time.Time, w io.Writer) error {

	/* If the file doesn't exist, make a new one */
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		if err := createConfig(fn, w); nil != err {
			return err
		}
		log.Printf("Created empty config file %v", fn)
		return nil
	}

	/* Skip it if file's not been modified */
	fi, err := os.Stat(fn)
	if nil != err {
		return err
	}
	if !(*last).Before(fi.ModTime()) {
		return nil
	}

	/* Slurp file, unJSON it */
	fb, err := ioutil.ReadFile(fn)
	*last = time.Now()
	if nil != err {
		return err
	}
	if err := loadConfig(fb, w); nil != err {
		return err
	}
	log.Printf("Read %v greetings from %v", len(GREETINGS), fn)
	log.Printf("Read %v responses from %v", len(RESPONSES), fn)
	log.Printf("Read %v bots from %v", len(BOTS), fn)
	log.Printf("Read %v channels from %v", len(CHANNELS), fn)
	log.Printf("Read %v guest channels from %v", len(GUESTCHANNELS), fn)

	return nil
}

/* createConfig creates a new, defaults-filled config file named fn */
func createConfig(fn string, w io.Writer) error {
	/* Json blob to write */
	ec, err := json.Marshal(defaultconfig)
	if nil != err {
		return err
	}
	/* Make it human-usable */
	buf := &bytes.Buffer{}
	if err := json.Indent(buf, ec, "", "        "); nil != err {
		return err
	}

	/* Load the default settings */
	if err := loadConfig(buf.Bytes(), w); nil != err {
		return err
	}

	/* Write it to the file */
	if err := ioutil.WriteFile(fn, buf.Bytes(), 0600); nil != err {
		return err
	}
	return nil
}

/* loadConfig updates the config from the json blob b */
func loadConfig(b []byte, w io.Writer) error {
	/* Unpack JSON */
	var c config
	if err := json.Unmarshal(b, &c); nil != err {
		return err
	}

	/* Update the internal settings */
	LOCK.Lock()
	defer LOCK.Unlock()
	GREETINGS = c.Greetings
	RESPONSES = c.Responses
	loadList(BOTS, c.Bots)

	/* Work out which channels we'd like to occupy */
	part := loadList(CHANNELS, c.Channels)
	part = append(part, loadList(GUESTCHANNELS, c.GuestChannels)...)
	makeChannels(CHANNELS)
	makeChannels(GUESTCHANNELS)

	/* Make sure nicks are all lower-case */
	toRem := []string{}
	for k, v := range GREETINGS {
		l := strings.ToLower(k)
		if k != l {
			GREETINGS[l] = v
			toRem = append(toRem, k)
		}
	}
	for _, r := range toRem {
		delete(GREETINGS, r)
	}

	/* Make sure we're in the right channels */
	for ch := range CHANNELS {
		Join(w, ch)
	}
	for ch := range GUESTCHANNELS {
		Join(w, ch)
	}
	for _, ch := range part {
		Part(w, ch)
	}

	return nil
}

/* loadList loads the contents of s into m.  It returns the elements of m
which were not in s (and which, then, are no longer in m). */
func loadList(m map[string]struct{}, s []string) []string {

	/* List of elements to delete.  Assume everything. */
	d := map[string]struct{}{}
	for k := range m {
		d[k] = struct{}{}
	}

	/* Add s to m, remove it from the list to be deleted */
	for _, v := range s {
		m[v] = struct{}{}
		delete(d, v)
	}

	/* Clean out removed elements */
	for k := range d {
		delete(m, k)
	}

	/* Return list of removed elements */
	ret := []string{}
	for k := range d {
		ret = append(ret, k)
	}
	return ret
}

/* makeChannels make sure the keys in the map have #'s */
func makeChannels(m map[string]struct{}) {
	for k := range m {
		if strings.HasPrefix(k, "#") {
			continue
		}
		delete(m, k)
		m["#"+k] = struct{}{}
	}
}
