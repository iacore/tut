package feed

import (
	"fmt"

	"github.com/gen2brain/beeep"

	"github.com/RasmusLindroth/tut/config"
)

func SendDesktopNotification(config config.Notification, nft DesktopNotification) {
	switch nft.Type {
	case DesktopNotificationFollower:
		if config.NotificationFollower {
			beeep.Notify(fmt.Sprintf("%s follows you", nft.Data), "", "")
		}
	case DesktopNotificationFavorite:
		if config.NotificationFavorite {
			beeep.Notify(fmt.Sprintf("%s favorited your toot", nft.Data), "", "")
		}
	case DesktopNotificationMention:
		if config.NotificationMention {
			beeep.Notify(fmt.Sprintf("%s mentioned you", nft.Data), "", "")
		}
	case DesktopNotificationUpdate:
		if config.NotificationUpdate {
			beeep.Notify(fmt.Sprintf("%s changed their toot", nft.Data), "", "")
		}
	case DesktopNotificationBoost:
		if config.NotificationBoost {
			beeep.Notify(fmt.Sprintf("%s boosted your toot", nft.Data), "", "")
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
