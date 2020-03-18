package bete

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/yi-jiayu/ted"
)

func TestBete_HandleTextMessage(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	stop := buildBusStop()
	filter := []string{"5", "24"}
	arrivals := buildDataMallBusArrival()
	chatID := randomInt64ID()
	text := must(FormatArrivalsByService(ArrivalInfo{
		Stop:     stop,
		Time:     refTime,
		Services: arrivals.Services,
		Filter:   filter,
	})).(string)
	req := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkup(stop.ID, filter),
	}

	b.Clock.(*MockClock).EXPECT().Now().Return(refTime)
	b.BusStops.(*MockBusStopRepository).EXPECT().Find(gomock.Any()).Return(stop, nil)
	b.Favourites.(*MockFavouriteRepository).EXPECT().FindByUserAndText(gomock.Any(), gomock.Any()).Return("")
	b.DataMall.(*MockDataMall).EXPECT().GetBusArrival(stop.ID, "").Return(arrivals, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: randomID()},
			Chat: ted.Chat{ID: chatID},
			Text: "96049 5 24",
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_HandleTextMessage_Favourite(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	stop := buildBusStop()
	filter := []string{"5", "24"}
	arrivals := buildDataMallBusArrival()
	userID := randomID()
	chatID := randomInt64ID()
	messageText := "SUTD"
	replyText := must(FormatArrivalsByService(ArrivalInfo{
		Stop:     stop,
		Time:     refTime,
		Services: arrivals.Services,
		Filter:   filter,
	})).(string)
	req := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        replyText,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkup(stop.ID, filter),
	}

	b.Clock.(*MockClock).EXPECT().Now().Return(refTime)
	b.BusStops.(*MockBusStopRepository).EXPECT().Find(gomock.Any()).Return(stop, nil)
	b.Favourites.(*MockFavouriteRepository).EXPECT().FindByUserAndText(userID, messageText).Return("96049 5 24")
	b.DataMall.(*MockDataMall).EXPECT().GetBusArrival(stop.ID, "").Return(arrivals, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			Text: messageText,
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_HandleReply_AddFavouriteQuery(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageText := "96049 5 24"
	req := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        fmt.Sprintf(AddFavouritePromptForName, messageText),
		ReplyMarkup: ted.ForceReply{},
	}

	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			ReplyToMessage: &ted.Message{
				Text: AddFavouritePromptForQuery,
			},
			Text: messageText,
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_HandleReply_AddFavouriteQuery_InvalidQuery(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageText := `Invalid Query: !@#$%^&*"`
	reportError := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        AddFavouriteReportQueryInvalid,
		ReplyMarkup: ted.ForceReply{},
	}
	askAgain := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        AddFavouritePromptForQuery,
		ReplyMarkup: ted.ForceReply{},
	}

	b.Telegram.(*MockTelegram).EXPECT().Do(reportError).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(askAgain).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			ReplyToMessage: &ted.Message{
				Text: AddFavouritePromptForQuery,
			},
			Text: messageText,
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_HandleCommand_Favourite(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	req := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        "What would you like to do?",
		ReplyMarkup: favouritesReplyMarkup(),
	}

	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID, Type: "private"},
			Text: "/favourites",
			Entities: []ted.MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: 11,
				},
			},
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_HandleCommand_Favourite_NonPrivateChat(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageID := randomID()
	req := ted.SendMessageRequest{
		ChatID:           chatID,
		Text:             "Sorry, you can only manage your favourites in a private chat.",
		ReplyToMessageID: messageID,
	}

	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			ID:   messageID,
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID, Type: "group"},
			Text: "/favourites",
			Entities: []ted.MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: 11,
				},
			},
		},
	}
	b.HandleUpdate(context.Background(), update)
}
