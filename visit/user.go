package visit

import (
	"fmt"
	"math/rand"
)

type User struct {
	Firstname string
	Lastname  string
	Email     string
	IP        string
	Language  string
	Useragent string
	Anonymous bool
}

func generateRandomUser(c *BotConfig, anonymous, mobile bool) *User {
	var user = new(User)
	user.Anonymous = anonymous
	if !anonymous {
		user.Firstname = c.Users.FirstNames[rand.Intn(len(c.Users.FirstNames))]
		user.Lastname = c.Users.LastNames[rand.Intn(len(c.Users.LastNames))]
		user.Email = fmt.Sprint(user.Firstname, ".", user.Lastname, c.Users.Emails[rand.Intn(len(c.Users.Emails))])
	}
	user.IP = c.Users.RandomIPs[rand.Intn(len(c.Users.RandomIPs))]
	user.Language = c.Users.Languages[rand.Intn(len(c.Users.Languages))]
	if mobile {
		user.Useragent = c.Users.MobileUserAgents[rand.Intn(len(c.Users.MobileUserAgents))]
	} else {
		agents := append(c.Users.UserAgents, c.Users.MobileUserAgents...)
		user.Useragent = agents[rand.Intn(len(agents))]
	}

	return user
}
