package bete

import (
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/yi-jiayu/datamall/v3"
	"github.com/yi-jiayu/ted"
)

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (c RealClock) Now() time.Time {
	return time.Now()
}

type DataMall interface {
	GetBusArrival(busStopCode string, serviceNo string) (datamall.BusArrival, error)
}

type Telegram interface {
	Do(request ted.Request) (ted.Response, error)
}

type Bete struct {
	Clock    Clock
	BusStops BusStopRepository
	DataMall DataMall
	Telegram Telegram
}

func (b Bete) SendETAMessage(chatID int, stopID string, filter []string) error {
	text, err := b.etaMessageText(stopID, filter)
	if err != nil {
		return err
	}
	req := ted.SendMessageRequest{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "HTML",
	}
	_, err = b.Telegram.Do(req)
	if err != nil {
		return errors.Wrap(err, "error sending telegram message")
	}
	return nil
}

func (b Bete) etaMessageText(stopID string, filter []string) (string, error) {
	t := b.Clock.Now()
	arrivals, err := b.DataMall.GetBusArrival(stopID, "")
	if err != nil {
		return "", errors.Wrap(err, "error getting bus arrivals")
	}
	var stop BusStop
	stop, err = b.BusStops.Find(stopID)
	if err != nil {
		log.Printf("error getting bus stop: %v", err)
		stop = BusStop{ID: stopID}
	}
	return FormatArrivalsByService(ArrivalInfo{
		Stop:     stop,
		Time:     t,
		Services: arrivals.Services,
		Filter:   filter,
	})
}