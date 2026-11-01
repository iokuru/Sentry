package caregiver

import "github.com/iokuru/sentry/internal/model"

type AlertDispatch string

const (
	DispatchUrgentAlert   AlertDispatch = "urgent_alert"
	DispatchStandardAlert AlertDispatch = "standard_alert"
	DispatchFeedOnly       AlertDispatch = "feed_only"
	DispatchQuiet          AlertDispatch = "quiet"
)

type NotificationChannels struct {
	Telegram    bool
	SNS         bool
	FeedOnly    bool
	Quiet       bool
	UrgentBadge bool
}

func DetermineAlertDispatch(level model.RiskLevel, isDamageControl bool) (AlertDispatch, NotificationChannels) {
	if isDamageControl {
		return DispatchUrgentAlert, NotificationChannels{
			Telegram:    true,
			SNS:         true,
			FeedOnly:    false,
			Quiet:       false,
			UrgentBadge: true,
		}
	}

	switch level {
	case model.RiskHigh:
		return DispatchStandardAlert, NotificationChannels{
			Telegram:    true,
			SNS:         true,
			FeedOnly:    false,
			Quiet:       false,
			UrgentBadge: false,
		}
	case model.RiskMedium:
		return DispatchFeedOnly, NotificationChannels{
			Telegram:    false,
			SNS:         false,
			FeedOnly:    true,
			Quiet:       false,
			UrgentBadge: false,
		}
	default:
		return DispatchQuiet, NotificationChannels{
			Telegram:    false,
			SNS:         false,
			FeedOnly:    false,
			Quiet:       true,
			UrgentBadge: false,
		}
	}
}
