package restrictedPage

import (
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
	"zeus-client/pkg"
)

type RestrictedPage struct {
	page         *playwright.Page
	myZeusButton string
	refreshIcon  string
}

func Instance(page *playwright.Page) *RestrictedPage {
	return &RestrictedPage{
		page:         page,
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
	defer pkg.Measure(time.Now(), "NavigateToMyZeusView")
	log.Debugln("Navigating to 'Mein ZEUS' view.")
	r.wait()
	if err := (*r.page).Click(r.myZeusButton); err != nil {
		log.WithField("error", err).Fatalln("Could not click on sidebar entry 'Mein ZEUS'.")
	}

	return instance(r)
}

func (r *RestrictedPage) wait() {
	defer pkg.Measure(time.Now(), "wait")
	log.Debugln("Waiting for the async loading to complete.")
	options := playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden}
	if _, err := (*r.page).WaitForSelector(r.refreshIcon, options); err != nil {
		log.WithField("error", err).Fatalln("Could not wait for the async loading to complete.")
	}

	attribute, err := (*r.page).GetAttribute(r.refreshIcon, "style")
	if err != nil {
		log.WithField("error", err).Fatalln("Could not get the style of the refresh icon.")
	}

	log.WithField("style", attribute).Debugln("Pressed refresh button successfully")
}

func (m *MyZeusView) GetStatus() PresenceStatus {
	defer pkg.Measure(time.Now(), "GetStatus")
	m.refresh()
	status, err := (*m.parent.page).GetAttribute(m.presenceStatus, "title")
	if err != nil {
		log.WithField("error", err).Fatalf("Could not get the presence status.")
	}

	status = strings.TrimSpace(strings.ToLower(status))

	var mappedStatus PresenceStatus
	switch status {
	case "anwesend":
		mappedStatus = Present
	case "abwesend":
		mappedStatus = Absent
	default:
		log.Fatalf("Status is not valid: %v", status)
	}

	return mappedStatus
}

func (m *MyZeusView) ToggleStatus() {
	defer pkg.Measure(time.Now(), "ToggleStatus")
	if err := (*m.parent.page).Click(m.terminalButton); err != nil {
		log.WithField("error", err).Fatalln("Could not click on terminal button 'Buchungen'.")
	}

	m.parent.wait()

	log.Debugln("Toggling the presence status")
	status := m.GetStatus()
	if status == Present {
		if err := (*m.parent.page).Click(m.absentButton); err != nil {
			log.WithField("error", err).Fatalln("Could not click on terminal button 'Geht'.")
		}
	} else if status == Absent {
		if err := (*m.parent.page).Click(m.presentButton); err != nil {
			log.WithField("error", err).Fatalln("Could not click on terminal button 'Kommt'.")
		}
	}

	newStatus := m.GetStatus()
	if newStatus == status {
		fields := log.Fields{"status": status, "newStatus": newStatus}
		log.WithFields(fields).Fatalln("Could not toggle the presence status.")
	}
	log.Infof("Status toggled from %s to %s", status, newStatus)
}

func (m *MyZeusView) refresh() {
	defer pkg.Measure(time.Now(), "refresh")
	log.Debugln("Refreshing the page")
	if err := (*m.parent.page).Click(m.refreshButton); err != nil {
		log.WithField("error", err).Fatalln("Could not click refresh button.")
	}

	m.parent.wait()
}

type PresenceStatus string

const (
	Absent  PresenceStatus = "absent"
	Present PresenceStatus = "present"
)
