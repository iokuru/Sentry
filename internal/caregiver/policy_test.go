package caregiver

import (
	"testing"

	"github.com/iokuru/sentry/internal/model"
)

func TestDetermineAlertDispatch(t *testing.T) {
	tests := []struct {
		name            string
		level           model.RiskLevel
		isDamageControl bool
		wantDispatch    AlertDispatch
		wantTelegram    bool
		wantSNS         bool
		wantFeedOnly    bool
		wantQuiet       bool
		wantUrgent      bool
	}{
		{
			name:            "damage control triggers urgent telegram and sns",
			level:           model.RiskHigh,
			isDamageControl: true,
			wantDispatch:    DispatchUrgentAlert,
			wantTelegram:    true,
			wantSNS:         true,
			wantFeedOnly:    false,
			wantQuiet:       false,
			wantUrgent:      true,
		},
		{
			name:            "damage control at medium risk still escalates to urgent alert",
			level:           model.RiskMedium,
			isDamageControl: true,
			wantDispatch:    DispatchUrgentAlert,
			wantTelegram:    true,
			wantSNS:         true,
			wantFeedOnly:    false,
			wantQuiet:       false,
			wantUrgent:      true,
		},
		{
			name:            "high risk scam triggers standard telegram and sns",
			level:           model.RiskHigh,
			isDamageControl: false,
			wantDispatch:    DispatchStandardAlert,
			wantTelegram:    true,
			wantSNS:         true,
			wantFeedOnly:    false,
			wantQuiet:       false,
			wantUrgent:      false,
		},
		{
			name:            "medium risk stays in feed only without push notifications",
			level:           model.RiskMedium,
			isDamageControl: false,
			wantDispatch:    DispatchFeedOnly,
			wantTelegram:    false,
			wantSNS:         false,
			wantFeedOnly:    true,
			wantQuiet:       false,
			wantUrgent:      false,
		},
		{
			name:            "low risk benign message stays quiet",
			level:           model.RiskLow,
			isDamageControl: false,
			wantDispatch:    DispatchQuiet,
			wantTelegram:    false,
			wantSNS:         false,
			wantFeedOnly:    false,
			wantQuiet:       true,
			wantUrgent:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dispatch, channels := DetermineAlertDispatch(tt.level, tt.isDamageControl)
			if dispatch != tt.wantDispatch {
				t.Fatalf("dispatch = %v, want %v", dispatch, tt.wantDispatch)
			}
			if channels.Telegram != tt.wantTelegram {
				t.Errorf("Telegram = %v, want %v", channels.Telegram, tt.wantTelegram)
			}
			if channels.SNS != tt.wantSNS {
				t.Errorf("SNS = %v, want %v", channels.SNS, tt.wantSNS)
			}
			if channels.FeedOnly != tt.wantFeedOnly {
				t.Errorf("FeedOnly = %v, want %v", channels.FeedOnly, tt.wantFeedOnly)
			}
			if channels.Quiet != tt.wantQuiet {
				t.Errorf("Quiet = %v, want %v", channels.Quiet, tt.wantQuiet)
			}
			if channels.UrgentBadge != tt.wantUrgent {
				t.Errorf("UrgentBadge = %v, want %v", channels.UrgentBadge, tt.wantUrgent)
			}
		})
	}
}
