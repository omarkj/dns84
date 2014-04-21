package main

import (
	"log"
	"os/user"
	"os"
	"fmt"
	"net/http/cookiejar"
	"net/http"
	flag "github.com/bgentry/pflag"
)

var (
	apiURL = "https://umsjon.1984.is"
	configFile = "/tmp/1984rc"
	cookieJar, _ = cookiejar.New(nil)
	client = &http.Client{
		Jar: cookieJar,
	}
	zone string
)

type Command struct {
	Run   func(cmd *Command, args []string)
	Name  string
	Flag  flag.FlagSet
}

var commands = []*Command{
	cmdLogin,
	cmdListZones,
	cmdListRecords,
	cmdDeleteRecord,
	cmdAddRecord,
}

func main() {
	usr, _ := user.Current()
	configFile = fmt.Sprintf("%s/.1984rc", usr.HomeDir)
	log.SetFlags(0)
	err := tryPrimeCookieJar()
	if err != nil {
		panic(err)
	}
	args := os.Args[1:]
	for _, cmd := range commands {
		if cmd.Name == args[0] {
			if err := cmd.Flag.Parse(args[1:]); err != nil {
				log.Fatalf(err.Error())
			}
			cmd.Run(cmd, cmd.Flag.Args())
		}
	}
}
