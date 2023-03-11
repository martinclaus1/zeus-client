package main

import (
	"flag"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/term"
	"io"
	"log"
	"os"
	"strings"
	"syscall"
	"zeus-client/loginPage"
)

var logger = logrus.New()

func main() {
	debugMode := flag.Bool("debug", true, "Enables debug mode")
	dryRun := flag.Bool("dry-run", false, "Does a dry run without toggling the presence state")
	silent := flag.Bool("silent", true, "Runs the selenium script in headless mode")
	username := flag.String("user", "", "Username for the ZEUS time tracking tool (mandatory)")
	password := flag.String("password", "", "Password for the ZEUS time tracking tool. If not provided, the script will prompt for the password.")

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "zeus-client is a script to toggle the presence status in ZEUS time tracking tool.\n\n")
		fmt.Fprintf(w, "options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(w, "\nexample usage:\n")
		fmt.Fprintf(w, "  zeus-client -user <username>\n")
	}

	flag.Parse()

	setupLogging(debugMode)

	flag.VisitAll(func(f *flag.Flag) {
		if f.Name != "password" {
			logger.Infof("%s: %s", f.Name, f.Value)
		}
	})

	if *username == "" {
		fmt.Println("Username is required")
		flag.Usage()
		os.Exit(1)
	}

	if *password == "" {
		fmt.Print("Password: ")
		bytepw, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			os.Exit(1)
		}
		p := strings.TrimSpace(string(bytepw))
		password = &p
		fmt.Println()
	}

	logger.Infoln("Starting the script")
	installPlaywright()

	page := getPage(silent)

	myZeusView := loginPage.Instance(&page, logger).Login(*username, *password).NavigateToMyZeusView()
	logger.Infof("Current status: %s", myZeusView.GetStatus())
	if *dryRun == false {
		myZeusView.ToggleStatus()
	}
}

func setupLogging(debugMode *bool) {
	log.SetOutput(logger.Writer())
	file, err := os.OpenFile("zeus-client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		mw := io.MultiWriter(os.Stdout, file)
		logger.SetOutput(mw)
	} else {
		logger.Info("Failed to logger to file, using default stderr")
	}
	if *debugMode == true {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
}

func installPlaywright() {
	logger.Infoln("Installing playwright")
	installOptions := playwright.RunOptions{Browsers: []string{"chromium"}}
	err := playwright.Install(&installOptions)
	if err != nil {
		logger.Fatalf("Could not install playwright: %v", err)
	}
}

func getPage(silent *bool) playwright.Page {
	logger.Debugf("Starting playwright in headless mode: %v", *silent)
	pw, err := playwright.Run()
	if err != nil {
		logger.Fatalf("Could not start playwright: %v", err)
	}

	options := playwright.BrowserTypeLaunchOptions{Headless: playwright.Bool(*silent)}

	browser, err := pw.Chromium.Launch(options)
	if err != nil {
		logger.Fatalf("Could not launch browser: %v", err)
	}

	page, err := browser.NewPage()
	if err != nil {
		logger.Fatalf("Could not create the page: %v", err)
	}

	if err = page.SetViewportSize(1920, 1080); err != nil {
		logger.Fatalf("Could not change the viewport size: %v", err)
	}

	return page
}
