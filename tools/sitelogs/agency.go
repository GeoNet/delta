package main

type Contact struct {
	Name               string `xml:"contact:name"`
	TelephonePrimary   string `xml:"contact:telephonePrimary"`
	TelephoneSecondary string `xml:"contact:telephoneSecondary"`
	Fax                string `xml:"contact:fax"`
	Email              string `xml:"contact:e-mail"`
}

type ContactInput struct {
	Name               string `xml:"name"`
	TelephonePrimary   string `xml:"telephonePrimary"`
	TelephoneSecondary string `xml:"telephoneSecondary"`
	Fax                string `xml:"fax"`
	Email              string `xml:"e-mail"`
}

type Agency struct {
	Agency                string  `xml:"contact:agency"`
	PreferredAbbreviation string  `xml:"contact:preferredAbbreviation"`
	MailingAddress        string  `xml:"contact:mailingAddress"`
	PrimaryContact        Contact `xml:"contact:primaryContact"`
	SecondaryContact      Contact `xml:"contact:secondaryContact"`
	Notes                 string  `xml:"contact:notes"`
}

type AgencyInput struct {
	Agency                string       `xml:"agency"`
	PreferredAbbreviation string       `xml:"preferredAbbreviation"`
	MailingAddress        string       `xml:"mailingAddress"`
	PrimaryContact        ContactInput `xml:"primaryContact"`
	SecondaryContact      ContactInput `xml:"secondaryContact"`
	Notes                 string       `xml:"notes"`
}
