package main

type Contact struct {
	Name               string `xml:"contact:name"`
	TelephonePrimary   string `xml:"contact:telephonePrimary"`
	TelephoneSecondary string `xml:"contact:telephoneSecondary"`
	Fax                string `xml:"contact:fax"`
	Email              string `xml:"contact:e-mail"`
}

type Agency struct {
	Agency                string  `xml:"contact:agency"`
	PreferredAbbreviation string  `xml:"contact:preferredAbbreviation"`
	MailingAddress        string  `xml:"contact:mailingAddress"`
	PrimaryContact        Contact `xml:"contact:primaryContact"`
	SecondaryContact      Contact `xml:"contact:secondaryContact"`
	Notes                 string  `xml:"contact:notes"`
}
