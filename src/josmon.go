package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"crypto/tls"
	"net/mail"
	"net/smtp"
	"time"
)

var (
	ConfFile      string = "josmon.conf"
	SitesDataFile string = "sites.data"
	CareerPages   string = "career_pages.cdf"
	AlertFile     string = "sendmail.text"
)

func errHandler(e error) {
	if e != nil {
		panic(e)
	}
}

func readFile(infile string) *[]string {
	var webTarget []string

	_, err := os.Stat(infile)
	if err == nil {
		file, err := os.Open(infile)
		errHandler(err)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			webTarget = append(webTarget, scanner.Text())
		}
	}
	return &webTarget
}

func readWebList(file string) *[]string {
	ptrURLs := readFile(file)
	return ptrURLs
}

func readSigCache(file string) *map[string]map[string]string {
	mapSig := make(map[string]map[string]string)

	ptrSigs := readFile(file)
	for _, line := range *ptrSigs {
		p := strings.Split(line, ",")
		mapSig[p[0]] = map[string]string{
			"sig":     p[1],
			"lastupd": p[2],
		}
	}
	return &mapSig
}

func getWebPage(url string) string {
	client := http.Client {
		Timeout: time.Second * 4,
	}
	response, err := client.Get(url)
	errHandler(err)
	defer response.Body.Close()

	html, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Warn: %s down.\n", url)
	}
	return string(html)
}

func getFocusContent(html string, scopeStart string, scopeEnd string) string {
	startPos := strings.Index(html, scopeStart)
	endPos := strings.Index(html, scopeEnd)
	if startPos >= 0 && endPos >= 0 {
		return html[startPos:endPos]
	} else {
		return ""
	}
}

func fingerprint(text string, lookfor string) string {
	// look for all occurences of keyword
	parts := strings.Split(text, lookfor)
	fingerprint := ""
	for i := 1; i < len(parts)-1; i++ {
		fingerprint += fmt.Sprintf("%d_", len(parts[i]))
	}
	return fingerprint
}

func getPageFingerprint(url string, scopeStart string, scopeEnd string, keyword string) string {
	webpage := getWebPage(url)
	focus := getFocusContent(webpage, scopeStart, scopeEnd)
	pageFP := fingerprint(focus, keyword)
	return pageFP
}

func writeFile(ofile string, lines *[]string) {
	file, err := os.OpenFile(ofile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	errHandler(err)
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range *lines {
		writer.WriteString(line + "\n")
	}
	writer.Flush()
}

func runCmp(siteFile string, cacheFile string) string {
	var sitesWithUpdate []string
	var allSitesStatus []string
	var alertBody string

	// read first from data cache
	ptrSigs := readSigCache(cacheFile)

	// read the site list
	ptrURLs := readFile(siteFile)
	for _, line := range *ptrURLs {
		cols := strings.Split(line, ",")
		newSig := getPageFingerprint(cols[0], cols[1], cols[2], "Engineer")

		// see if Sig is in cache
		cSigs := *ptrSigs
		if cMap, ok := cSigs[cols[0]]; ok {
			if newSig != cMap["sig"] {
				sitesWithUpdate = append(sitesWithUpdate, cols[0])
				alertBody += cols[0] + "\r\n"
			}
		}
		allSitesStatus = append(allSitesStatus, fmt.Sprintf("%s,%s,%d", cols[0], newSig, 0))
	}
	writeFile(AlertFile, &sitesWithUpdate)
	writeFile(cacheFile, &allSitesStatus)
	return alertBody
}

func readConf() *map[string]string {
	conf := make(map[string]string)
	ptrConf := readFile(ConfFile)
	for _, line := range *ptrConf {
		cols := strings.Split(line, " ")
		// skip lines starting with # or split returns blank on col#2
		if cols[1] == "" || string(line[0]) == "#" {
			continue
		}
		conf[cols[0]] = cols[1]
	}
	return &conf
}

func sendMail(conf map[string]string, body string) {
	from := mail.Address{Name: "Sender", Address: conf["email_from"]}
	to := mail.Address{Name: "Recipient", Address: conf["email_to"]}
	auth := smtp.PlainAuth("", conf["smtp_user"], conf["smtp_pass"], conf["smtp_host"])
	tlsConfig := &tls.Config{
		ServerName: conf["smtp_host"],
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", conf["smtp_host"], conf["smtp_port"]), tlsConfig)
	errHandler(err)

	client, err := smtp.NewClient(conn, conf["smtp_host"])
	errHandler(err)

	defer client.Close()

	err = client.Auth(auth)
	errHandler(err)
	err = client.Mail(from.Address)
	errHandler(err)
	err = client.Rcpt(to.Address)

	writer, err := client.Data()
	errHandler(err)
	_, err = fmt.Fprintf(writer, "From: %s\r\n", conf["email_from"])
	errHandler(err)
	_, err = fmt.Fprintf(writer, "To: %s\r\n", conf["email_to"])
	errHandler(err)
	_, err = fmt.Fprintf(writer, "Subject: Jobs monitor alert\r\n")
	errHandler(err)
	_, err = fmt.Fprintf(writer, "\r\n%s\r\n", body)
	errHandler(err)
	err = writer.Close()
}

func main() {
	// flag vars
	var pURL string
	var pBegin string
	var pEnd string
	var pFile string
	var pHelp bool

	// retrieve from page, output focused text
	flag.StringVar(&pURL, "url", "", "Career page URL.")
	flag.StringVar(&pBegin, "begin", "", "Starting string of text of focus.")
	flag.StringVar(&pEnd, "end", "", "Ending string of text of focus.")

	// read focused text from file
	flag.StringVar(&pFile, "intext", "", "Read from file and get keyword positions.")

	// help
	flag.BoolVar(&pHelp, "help", false, "Show help.")

	flag.Parse()

	if len(pURL) > 0 { // URL is provided in parameters
		if pBegin == "" || pEnd == "" {
			fmt.Println("Usage: ", os.Args[0], "[[--url] [--begin] [--end]]")
		} else {
			html := getWebPage(pURL)
			fmt.Println(getFocusContent(html, pBegin, pEnd))
		}
	} else if len(pFile) > 0 { // focused input file is provided
		fdata, err := os.ReadFile(pFile)
		errHandler(err)
		// fingerprint keyword
		fmt.Println(fingerprint(string(fdata), "Engineer"))
	} else if pHelp { // help
		cmd := strings.Split(os.Args[0], "\\")
		bin := cmd[len(cmd)-1]
		fmt.Println("Usage:", bin, "(reads from career_pages.cdf)")
		fmt.Println("   To test your keywords and see if it will grab the correct portion ")
		fmt.Println("   of the webpage. Useful for piping to file and use with --intext.")
		fmt.Println("      ", bin, "[[--url <URL>] [--begin <keyword>] [--end <keyword>]]")
		fmt.Println("   To test the fingerprinting using the output from --url.")
		fmt.Println("      ", bin, "[--intext <file>]")
	} else { // no parameter specified, use comma-delimited file - normal operation
		cnf := readConf()
		alertText := runCmp(CareerPages, SitesDataFile)
		if alertText != "" {
			sendMail(*cnf, alertText)
		}
	}
}
