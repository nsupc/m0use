package main

import (
	"fmt"
	"log"
	"log/slog"
	"m0use/pkg/config"
	"m0use/pkg/ns"
	"os"
	"strconv"

	"github.com/nsupc/eurogo/client"
	"github.com/nsupc/eurogo/models"
)

func main() {
	var path string

	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		path = "config.yml"
	}

	conf, err := config.ReadConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	config.InitLogger(conf)

	var ignore []string

	if conf.Cache.IsActive {
		ignore, err = conf.ReadCache()
		if err != nil {
			slog.Error("cache read error", slog.Any("error", err))
			os.Exit(1)
		}

		slog.Info(fmt.Sprintf("skipping %d cached nations", len(ignore)))
	} else {
		ignore = []string{}
	}

	nsClient := ns.NewClient(conf.User, conf.RequestRate)

	canRecruit, cannotRecruit, err := nsClient.GetRecruitmentEligibleNations(conf.Region, conf.Region, ignore)
	if err != nil {
		slog.Error("unable to retrieve nations", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info(fmt.Sprintf("%d/%d checked nations have recruitment telegrams enabled", len(canRecruit), len(canRecruit)+len(cannotRecruit)))

	err = conf.WriteCache(append(ignore, cannotRecruit...))
	if err != nil {
		slog.Warn("unable to update cache", slog.Any("error", err))
	}

	eurocoreClient := client.New(conf.Eurocore.Username, conf.Eurocore.Password, conf.Eurocore.Url)

	telegrams := []models.NewTelegram{}

	if conf.Telegram.Template != "" {
		template, err := eurocoreClient.GetTemplate(conf.Telegram.Template)
		if err != nil {
			slog.Error("unable to retrieve telegram template", slog.Any("error", err))
			os.Exit(1)
		}

		for _, recipient := range canRecruit {
			telegram := models.NewTelegram{
				Sender:    template.Nation,
				Recipient: recipient,
				Id:        strconv.Itoa(template.Tgid),
				Secret:    template.Key,
				Type:      "standard",
			}

			telegrams = append(telegrams, telegram)
		}
	} else {
		for _, recipient := range canRecruit {
			telegram := models.NewTelegram{
				Sender:    conf.Telegram.Author,
				Recipient: recipient,
				Id:        strconv.Itoa(conf.Telegram.Id),
				Secret:    conf.Telegram.Key,
				Type:      "standard",
			}

			telegrams = append(telegrams, telegram)
		}
	}

	err = eurocoreClient.SendTelegrams(telegrams)
	if err != nil {
		slog.Error("unable to send telegrams", slog.Any("error", err))
	} else {
		slog.Info("telegram request sent successfully")
	}
}
