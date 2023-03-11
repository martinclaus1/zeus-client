package pkg

import (
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
)

func InstallPlaywright() {
	log.Debugln("Installing playwright")
	installOptions := playwright.RunOptions{Browsers: []string{"chromium"}}
	err := playwright.Install(&installOptions)
	if err != nil {
		log.WithField("error", err).Fatalln("Could not install playwright.")
	}
}

func GetPage(silent *bool) playwright.Page {
	log.Debugf("Starting playwright in headless mode: %v", *silent)
	pw, err := playwright.Run()
	if err != nil {
		log.WithField("error", err).Fatalln("Could not start playwright.")
	}

	options := playwright.BrowserTypeLaunchOptions{Headless: playwright.Bool(*silent)}

	browser, err := pw.Chromium.Launch(options)
	if err != nil {
		log.WithField("error", err).Fatalln("Could not launch browser.")
	}

	page, err := browser.NewPage()
	page.SetDefaultTimeout(60000)
	if err != nil {
		log.WithField("error", err).Fatalln("Could not create the page.")
	}

	if err = page.SetViewportSize(1920, 1080); err != nil {
		log.WithField("error", err).Fatalln("Could not change the viewport size.")
	}

	return page
}
