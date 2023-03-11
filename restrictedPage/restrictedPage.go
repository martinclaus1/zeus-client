package restrictedPage

import (
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type RestrictedPage struct {
	page         *playwright.Page
	myZeusButton string
	logger       *logrus.Logger
	refreshIcon  string
}

func Instance(page *playwright.Page, logger *logrus.Logger) *RestrictedPage {
	return &RestrictedPage{
		page:         page,
		logger:       logger,
		myZeusButton: "//div[@title='Mein ZEUS']",
		refreshIcon:  "i#workspaceCallbackProcess",
	}
}

type MyZeusView struct {
	parent         *RestrictedPage
	presenceStatus string
	refreshButton  string
	terminalButton string
	presentButton  string
	absentButton   string
}

func instance(parent *RestrictedPage) *MyZeusView {
	return &MyZeusView{
		parent:         parent,
		presenceStatus: ".account-info-status-bar img",
		refreshButton:  "//i[@title='Widget Daten aktualisieren' and contains(@onclick,'AccountInfo')]",
		terminalButton: "//button[@title='Buchungen']",
		presentButton:  "button#TerminalButton0",
		absentButton:   "button#TerminalButton1",
	}
}

func (r *RestrictedPage) NavigateToMyZeusView() *MyZeusView {
	r.logger.Infoln("Navigating to Mein ZEUS view")
	r.wait()
	if err := (*r.page).Click(r.myZeusButton); err != nil {
		r.logger.Fatalf("Could not click on Mein ZEUS: %v", err)
	}

	return instance(r)
}

func (r *RestrictedPage) wait() {
	start := time.Now()
	r.logger.Debugf("Waiting for the async loading to complete")
	options := playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden}
	_, err := (*r.page).WaitForSelector(r.refreshIcon, options)
	if err != nil {
		r.logger.Fatalf("Could not wait for the page to load: %v", err)
	}

	attribute, err := (*r.page).GetAttribute(r.refreshIcon, "style")
	if err != nil {
		r.logger.Fatalf("Could not get the attribute: %v", err)
	}

	if attribute != "" {
		r.logger.Debugf("Current style: %s, ms: %d", attribute, time.Since(start).Milliseconds())
	}
}

func (m *MyZeusView) GetStatus() string {
	m.refresh()
	status, err := (*m.parent.page).GetAttribute(m.presenceStatus, "title")
	if err != nil {
		m.parent.logger.Fatalf("Could not get the presence status: %v", err)
	}

	status = strings.TrimSpace(strings.ToLower(status))
	if status != "anwesend" && status != "abwesend" {
		m.parent.logger.Fatalf("Status is not valid: %v", status)
	}

	return status
}

func (m *MyZeusView) ToggleStatus() {
	if err := (*m.parent.page).Click(m.terminalButton); err != nil {
		m.parent.logger.Fatalf("Could not click on Buchungen: %v", err)
	}

	m.parent.wait()

	m.parent.logger.Infoln("Toggling the presence status")
	status := m.GetStatus()
	if status == "anwesend" {
		if err := (*m.parent.page).Click(m.absentButton); err != nil {
			m.parent.logger.Fatalf("Could not click on Abwesend: %v", err)
		}
	} else if status == "abwesend" {
		if err := (*m.parent.page).Click(m.presentButton); err != nil {
			m.parent.logger.Fatalf("Could not click on Anwesend: %v", err)
		}
	}

	newStatus := m.GetStatus()
	if newStatus == status {
		m.parent.logger.Fatalf("Could not toggle the presence status")
	}
	m.parent.logger.Infof("Status toggled from %s to %s", status, newStatus)
}

func (m *MyZeusView) refresh() {
	m.parent.logger.Debugln("Refreshing the page")
	if err := (*m.parent.page).Click(m.refreshButton); err != nil {
		m.parent.logger.Fatalf("Could not click refresh button: %v", err)
	}

	m.parent.wait()
}
