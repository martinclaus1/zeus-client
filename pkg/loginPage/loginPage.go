package loginPage

//go:generate mockgen -destination=../mocks/mock_page.go -package=mocks github.com/playwright-community/playwright-go Page

import (
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
	"time"
	"zeus-client/pkg"
	"zeus-client/pkg/restrictedPage"
)

type LoginPage struct {
	url           string
	page          *playwright.Page
	usernameInput string
	passwordInput string
	loginButton   string
}

func Instance(page *playwright.Page) *LoginPage {
	return &LoginPage{
		url:           "https://saas.isgus.de/zeusxF/Environment/Account/LogOn.aspx?&AspxAutoDetectCookieSupport=0",
		page:          page,
		usernameInput: "//input[@name='uiUserName']",
		passwordInput: "//input[@name='uiPassword']",
		loginButton:   "div#uiLogOnButton",
	}
}

func (l *LoginPage) Login(username string, password string) *restrictedPage.RestrictedPage {
	defer pkg.Measure(time.Now(), "Login")
	log.Infoln("Navigating to ZEUS login page.")
	if _, err := (*l.page).Goto(l.url); err != nil {
		log.WithField("error", err).Fatalln("Could not navigate to ZEUS login page.")
	}

	log.Debugln("Type username")
	if err := (*l.page).Type(l.usernameInput, username); err != nil {
		log.WithField("error", err).Fatalln("Could not type user name.")
	}

	log.Debugln("Type password")
	if err := (*l.page).Type(l.passwordInput, password); err != nil {
		log.WithField("error", err).Fatalln("Could not type password.")
	}

	log.Infoln("Logging in")
	if err := (*l.page).Click(l.loginButton); err != nil {
		log.WithField("error", err).Fatalln("Could not login.")
	}

	return restrictedPage.Instance(l.page)
}
