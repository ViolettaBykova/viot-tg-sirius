package bot

import (
	"github.com/ViolettaBykova/viot-tg-sirius/pkg/weather"
	"github.com/robfig/cron/v3"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Option func(s *Service)

type Service struct {
	store      Storage
	bot        *tb.Bot
	cron       *cron.Cron
	cronJobs   map[int64]int // Карта для хранения задач по ID пользователей
	weatherAPI *weather.Client
}

func New(store Storage, bot *tb.Bot, cron *cron.Cron, weatherAPI *weather.Client, opts ...Option) *Service {
	s := &Service{
		store:      store,
		bot:        bot,
		cron:       cron,
		cronJobs:   make(map[int64]int),
		weatherAPI: weatherAPI,
	}

	for _, applyOpt := range opts {
		applyOpt(s)
	}

	s.loadScheduledJobs()
	return s
}
