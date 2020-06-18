package cron_test

import (
	"testing"

	"github.com/enrico5b1b4/telegram-bot/cron"
	"github.com/stretchr/testify/require"
)

func TestStartAndStopScheduler(t *testing.T) {
	scheduler := cron.NewScheduler()
	scheduler.Start()

	err := scheduler.Stop()
	require.NoError(t, err)
}
