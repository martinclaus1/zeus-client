package loginPage

//go:generate mockgen -destination=../mocks/mock_page.go -package=mocks github.com/playwright-community/playwright-go Page

import (
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
)
import "zeus-client/restrictedPage"

type LoginPage struct {
	url           string
	page          *playwright.Page
	usernameInput string
	passwordInput string
	loginButton   string
	logger        *logrus.Logger
}

func Instance(page *playwright.Page, logger *logrus.Logger) *LoginPage {
	return &LoginPage{
		url:           "https://saas.isgus.de/zeusxF/Environment/Account/LogOn.aspx?&AspxAutoDetectCookieSupport=0",
		page:          page,
		usernameInput: "//input[@name='uiUserName']",
		passwordInput: "//input[@name='uiPassword']",
		loginButton:   "div#uiLogOnButton",
		logger:        logger,
	}
}

func (l *LoginPage) Login(username string, password string) *restrictedPage.RestrictedPage {
	l.logger.Infoln("Navigating to ZEUS login page")
	if _, err := (*l.page).Goto(l.url); err != nil {
		l.logger.Fatalf("Could not navigate to the Zeus service: %v", err)
	}

	l.logger.Debugln("Type username")
	if err := (*l.page).Type(l.usernameInput, username); err != nil {
		l.logger.Fatalf("Could not type user name: %v", err)
	}

	l.logger.Debugln("Type password")
	if err := (*l.page).Type(l.passwordInput, password); err != nil {
		l.logger.Fatalf("Could not type password: %v", err)
	}

	l.logger.Infoln("Logging in")
	if err := (*l.page).Click(l.loginButton); err != nil {
		l.logger.Fatalf("Could not login: %v", err)
	}

	return restrictedPage.Instance(l.page, l.logger)
}
