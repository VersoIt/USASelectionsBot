package main

import (
	"context"
	"errors"
	"github.com/VersoIt/learning/internal/notifier"
	"github.com/VersoIt/learning/internal/parser"
	"github.com/VersoIt/learning/internal/service"
	"github.com/VersoIt/learning/pkg/cfg"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := cfg.Get()
	electionsParser := parser.NewElection(config.ParseUrl, config.FirstCandidateContainer, config.SecondCandidateContainer)
	electionService := service.NewElectionFetcher(electionsParser)
	tgNotifier, err := notifier.NewTelegram(electionService, config.BotToken, config.ChannelId, config.UpdateDelay, notifier.Candidates{FirstCandidateName: config.FirstCandidateName, SecondCandidateName: config.SecondCandidateName})
	if err != nil {
		log.Println(err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := tgNotifier.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Println(err)
			}
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-exit
	cancel()
	log.Println("fetcher closed!")
}
