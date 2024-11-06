package notifier

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"time"
)

type electionFetcher interface {
	Fetch() (int, int, bool, error)
}

type Candidates struct {
	FirstCandidateName  string
	SecondCandidateName string
}

type Telegram struct {
	fetcher     electionFetcher
	updateDelay time.Duration
	botToken    string
	channelId   string
	botApi      *tgbotapi.BotAPI
	candidates  Candidates
}

func NewTelegram(fetcher electionFetcher, botToken, channelId string, updateDelay time.Duration, candidates Candidates) (*Telegram, error) {
	botApi, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}
	return &Telegram{fetcher: fetcher, botToken: botToken, channelId: channelId, updateDelay: updateDelay, botApi: botApi, candidates: candidates}, nil
}

func (t *Telegram) Start(ctx context.Context) error {
	if err := t.sendMessage(fmt.Sprintf("Bot started: %s vs %s", t.candidates.FirstCandidateName, t.candidates.SecondCandidateName)); err != nil {
		return err
	}

	log.Println("tg fetcher started!")
	for {
		time.Sleep(time.Second * 1)
		log.Println("fetching candidates data...")
		firstCandidateRes, secondCandidateRes, post, err := t.fetcher.Fetch()
		if err != nil {
			log.Println("failed to fetch candidates data:", err)
			continue
		}

		if post {
			if err = t.sendMessage(fmt.Sprintf("%s: %d\n%s: %d", t.candidates.FirstCandidateName, firstCandidateRes, t.candidates.SecondCandidateName, secondCandidateRes)); err != nil {
				log.Printf("failed to send message: %v", err)
			} else {
				log.Printf("posted %d:%d", firstCandidateRes, secondCandidateRes)
			}
		} else {
			log.Println("post canceled with same data from prev request")
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
}

func (t *Telegram) sendMessage(content string) error {
	msg := tgbotapi.NewMessageToChannel(t.channelId, content)
	_, err := t.botApi.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
