package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/mail"
	"sync"
	"time"

	"github.com/emersion/go-smtp"
)

type Email struct {
	ID        int64     `json:"id"`
	From      string    `json:"from"`
	To        []string  `json:"to"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	Headers   string    `json:"headers"`
	CreatedAt time.Time `json:"created_at"`
}

type Storage struct {
	sync.RWMutex
	Emails []Email
	nextID int64
}

var store = &Storage{
	Emails: []Email{},
	nextID: 1,
}

func (s *Storage) Add(e *Email) {
	s.Lock()
	defer s.Unlock()
	e.ID = s.nextID
	s.nextID++
	s.Emails = append([]Email{*e}, s.Emails...)
	if len(s.Emails) > 100 {
		s.Emails = s.Emails[:100]
	}
}

func (s *Storage) GetAll() []Email {
	s.RLock()
	defer s.RUnlock()
	result := make([]Email, len(s.Emails))
	copy(result, s.Emails)
	return result
}

type Backend struct{}

func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{}, nil
}

type Session struct {
	From string
	To   []string
}

func (s *Session) AuthPlain(username, password string) error {
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.From = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.To = append(s.To, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	var subject, bodyStr string
	if err == nil {
		subject = msg.Header.Get("Subject")
		b, _ := io.ReadAll(msg.Body)
		bodyStr = string(b)
	} else {
		bodyStr = string(data)
		log.Printf("Failed to parse mail: %v", err)
	}

	email := &Email{
		From:      s.From,
		To:        s.To,
		Subject:   subject,
		Body:      bodyStr,
		CreatedAt: time.Now(),
		Headers:   string(data),
	}

	store.Add(email)
	log.Printf("Received mail from %s", s.From)
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func getEmails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store.GetAll())
}

func main() {
	be := &Backend{}
	s := smtp.NewServer(be)
	s.Addr = "127.0.0.1:1025"
	s.Domain = "localhost"
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	log.Println("Starting SMTP Server at :1025")
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/api/emails", getEmails)

	log.Println("Serving frontend static files at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
