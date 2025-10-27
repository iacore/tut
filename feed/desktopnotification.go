package feed

import (
	"fmt"

	"github.com/gen2brain/beeep"

	"github.com/RasmusLindroth/tut/config"
)

func SendDesktopNotification(config config.Notification, notif DesktopNotification) {
	switch notif.Type {
	case DesktopNotificationFollower:
		if config.NotificationFollower {
			beeep.Notify(fmt.Sprintf("%s follows you", notif.Data), "", "")
		}
	case DesktopNotificationFavorite:
		if config.NotificationFavorite {
			beeep.Notify(fmt.Sprintf("%s favorited your toot", notif.Data), "", "")
		}
	case DesktopNotificationMention:
		if config.NotificationMention {
			beeep.Notify(fmt.Sprintf("%s mentioned you", notif.Data), "", "")
		}
	case DesktopNotificationUpdate:
		if config.NotificationUpdate {
			beeep.Notify(fmt.Sprintf("%s changed their toot", notif.Data), "", "")
		}
	case DesktopNotificationBoost:
		if config.NotificationBoost {
			beeep.Notify(fmt.Sprintf("%s boosted your toot", notif.Data), "", "")
		}
	case DesktopNotificationPoll:
		if config.NotificationPoll {
			beeep.Notify("Poll has ended", "", "")
		}
	case DesktopNotificationPost:
		if config.NotificationPost {
			beeep.Notify("New post", "", "")
		}
	}
}
