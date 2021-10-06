package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/jung-kurt/gofpdf"
)

type Email struct {
	Date       time.Time       `json:"date"`
	From       []*mail.Address `json:"from"`
	To         []*mail.Address `json:"to"`
	Subject    string          `json:"subject"`
	Body       []string        `json:"body"`
	Attachment string          `json:"attachment"`
}

func main() {
	//readEmailSubject()
	emailData := readEmailMessage() //Read Message
	//readJsonFile()
	GeneratePdf(emailData)

}

func GeneratePdf(emailData []Email) {
	for i := 0; i < len(emailData); i++ {
		fmt.Println("emailData", emailData[i].Body)

		date := emailData[i].Date.String()
		from := emailData[i].From[0].String()
		subject := emailData[i].Subject
		body := emailData[i].Body[0]
		attach := emailData[i].Attachment

		filename := date + ".pdf"

		pdf := gofpdf.New("P", "mm", "A4", "")
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 16)

		// CellFormat(width, height, text, border, position after, align, fill, link, linkStr)
		pdf.CellFormat(190, 7, "Email from "+from, "0", 0, "CM", false, 0, "")
		pdf.Ln(12)

		pdf.CellFormat(190, 7, date, "0", 0, "CM", false, 0, "")
		pdf.Ln(12)

		pdf.CellFormat(190, 7, subject, "0", 0, "CM", false, 0, "")
		pdf.Ln(12)

		pdf.CellFormat(190, 15, body, "0", 0, "CM", false, 0, "")
		pdf.Ln(12)

		pdf.CellFormat(190, 20, attach, "0", 0, "CM", false, 0, "")

		// ImageOptions(src, x, y, width, height, flow, options, link, linkStr)
		pdf.ImageOptions(
			"email.jpeg",
			100, 120,
			20, 20,
			false,
			gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true},
			0,
			"",
		)

		pdf.OutputFileAndClose(filename)
	}
}

func readJsonFile() {
	// Open our jsonFile
	jsonFile, err := os.Open("email.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully opened email.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var results map[string]interface{}

	json.Unmarshal([]byte(byteValue), &results)
	fmt.Println(results)
	//	for i := 0; i < len(results); i++ {
	//		fmt.Println("date : ", results["date"])
	//	}

	//fmt.Println("email json : ", result[""]{0})

}

func readEmailMessage() []Email {
	allmail := make([]Email, 0)
	var em_date time.Time
	var em_from []*mail.Address
	var em_to []*mail.Address
	var em_subject string
	em_body := make([]string, 0)
	var em_file string

	fmt.Println("Connecting to server...")

	//c, err := client.DialTLS("mail.example.org:993", nil)
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login("test1day2020@gmail.com", "*****password*****"); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Logged in")

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		fmt.Println(err)
	}

	// Get the last message
	if mbox.Messages == 0 {
		fmt.Println("No message in mailbox")
	}
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(mbox.Messages)

	// Get the whole message body
	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 1)
	go func() {
		if err := c.Fetch(seqSet, items, messages); err != nil {
			fmt.Println(err)
		}
	}()

	msg := <-messages
	if msg == nil {
		fmt.Println("Server didn't returned message")
	}

	r := msg.GetBody(&section)
	if r == nil {
		fmt.Println("Server didn't returned message body")
	}

	// Create a new mail reader
	mr, err := mail.CreateReader(r)
	if err != nil {
		fmt.Println(err)
	}

	// Print some info about the message
	header := mr.Header
	if date, err := header.Date(); err == nil {
		em_date = date
		fmt.Println("Date:", date)
	}
	if from, err := header.AddressList("From"); err == nil {
		em_from = from
		fmt.Println("From:", from)
	}
	if to, err := header.AddressList("To"); err == nil {
		em_to = to
		fmt.Println("To:", to)
	}
	if subject, err := header.Subject(); err == nil {
		em_subject = subject
		fmt.Println("Subject:", subject)
	}

	// Process each message's part
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			// This is the message's text (can be plain-text or HTML)
			b, _ := ioutil.ReadAll(p.Body)
			fmt.Println("Got text: %v", string(b))
			em_body = append(em_body, string(b[:]))
		case *mail.AttachmentHeader:
			// This is an attachment
			fileName, _ := h.Filename()
			em_file = fileName
			fmt.Println("Got attachment: %v", fileName)
		}
	}
	email := Email{
		Date:       em_date,
		From:       em_from,
		To:         em_to,
		Subject:    em_subject,
		Attachment: em_file,
		Body:       em_body,
	}

	allmail = append(allmail, email)
	writeJSON(allmail)
	return allmail
}

func writeJSON(data []Email) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	_ = ioutil.WriteFile("email.json", file, 0644)
}

//Not using for now
func readEmailSubject() {
	fmt.Println("Connecting to server...")

	// Connect to server

	//c, err := client.DialTLS("mail.example.org:993", nil)
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login("test1day2020@gmail.com", "strongword"); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	fmt.Println("Mailboxes:")
	for m := range mailboxes {
		fmt.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		fmt.Println(err)
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Flags for INBOX:", mbox.Flags)

	// Get the last 2 messages
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 1 {
		// We're using unsigned integers here, only subtract if the result is > 0
		from = mbox.Messages - 1
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done = make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	fmt.Println("Last 2 messages:")
	for msg := range messages {
		fmt.Println("* "+msg.Envelope.Subject, msg.Envelope.From[0], msg.Envelope.To[0], msg.Envelope.Date)
		//literal := msg.GetBody("BODY[]")
		//fmt.Println(literal)
	}

	if err := <-done; err != nil {
		fmt.Println(err)
	}

	fmt.Println("Done!")
}
