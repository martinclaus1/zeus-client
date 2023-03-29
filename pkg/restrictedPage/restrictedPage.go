package restrictedPage

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/martinclaus1/zeus-client/pkg"
	"github.com/playwright-community/playwright-go"
	"github.com/rodaine/table"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
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

type String string

func (s String) trim() String {
	return String(strings.TrimSpace(string(s)))
}

func (s String) replace(old, new string) String {
	return String(strings.ReplaceAll(string(s), old, new))
}

func (s String) toLower() String {
	return String(strings.ToLower(string(s)))
}

func (s String) toFloat() (float64, error) {
	value, err := strconv.ParseFloat(string(s), 64)
	if err != nil {
		return -1, err
	}
	return value, nil
}

func (m *MyZeusView) GetOverview() (*Overview, error) {
	defer pkg.Measure(time.Now(), "PrintOverview")
	m.refresh()
	elements, err := (*m.parent.page).QuerySelectorAll(".account-list li.account-list-item")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not query the overview: %v", err))
	}

	result := make(map[string]float64)
	for _, element := range elements {
		text, err := element.InnerText()
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Could not get inner text: %v", err))
		}
		item := strings.Split(text, "\n")
		key := String(item[0]).trim().toLower()
		value, _ := String(item[1]).trim().replace(",", ".").toFloat()

		result[string(key)] = value
	}

	return &Overview{
		balance:         result["saldo"],
		netDiffDay:      result["differenz netto tag"],
		grossDiffDay:    result["brutto tag"],
		grossDiffWeek:   result["brutto wo"],
		grossDiffPeriod: result["brutto periode"],
		vacationTotal:   result["urlaub gesamt"],
		vacationUsed:    result["urlaub genommen"],
		vacationLeft:    result["urlaub unverplant"],
	}, nil
}

type Overview struct {
	balance         float64
	netDiffDay      float64
	grossDiffDay    float64
	grossDiffWeek   float64
	grossDiffPeriod float64
	vacationTotal   float64
	vacationUsed    float64
	vacationLeft    float64
}

func (o *Overview) Print() {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Position", "Value")
	tbl.WithPadding(6)
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	tbl.AddRow("Balance", fmt.Sprintf("%.2f", o.balance))
	tbl.AddRow("Net Difference Day", fmt.Sprintf("%.2f", o.netDiffDay))
	tbl.AddRow("Gross Difference Day", fmt.Sprintf("%.2f", o.grossDiffDay))
	tbl.AddRow("Gross Difference Week", fmt.Sprintf("%.2f", o.grossDiffWeek))
	tbl.AddRow("Gross Difference Period", fmt.Sprintf("%.2f", o.grossDiffPeriod))
	tbl.AddRow("Vacation Total", fmt.Sprintf("%.2f", o.vacationTotal))
	tbl.AddRow("Vacation Used", fmt.Sprintf("%.2f", o.vacationUsed))
	tbl.AddRow("Vacation Left", fmt.Sprintf("%.2f", o.vacationLeft))

	tbl.Print()
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
