// Package telemux is a flexible message router add-on for "go-telegram-bot-api".
//
// Make sure to check "go-telegram-bot-api" documentation first:
// https://github.com/go-telegram-bot-api/telegram-bot-api
package telemux

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Recoverer is a function that handles panics.
type Recoverer = func(*Update, error)

// Mux is a container for handlers.
type Mux struct {
	Handlers  []*Handler
	Recoverer Recoverer
}

// NewMux creates new Mux.
func NewMux() *Mux {
	return &Mux{}
}

// AddHandler adds handler to Mux.
func (m *Mux) AddHandler(h *Handler) *Mux {
	m.Handlers = append(m.Handlers, h)
	return m
}

// SetRecoverer registers a function to call when a panic occurs.
func (m *Mux) SetRecoverer(recoverer Recoverer) *Mux {
	m.Recoverer = recoverer
	return m
}

// Dispatch tells Mux to process the update.
func (m *Mux) Dispatch(u tgbotapi.Update) bool {
	uu := Update{u}

	defer func() {
		if err, ok := recover().(error); ok {
			if m.Recoverer != nil {
				m.Recoverer(&uu, error(err))
			} else {
				panic(err)
			}
		}
	}()

	for _, handler := range m.Handlers {
		if handler.Filter(&uu) {
			handler.Handle(&uu)
			return true
		}
	}
	return false
}
