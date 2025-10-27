package feed

import (
	"context"
	"errors"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Arceliar/phony"
	"github.com/RasmusLindroth/go-mastodon"
	"golang.org/x/exp/slices"

	"github.com/RasmusLindroth/tut/api"
	"github.com/RasmusLindroth/tut/config"
)

type apiFunc func(pg *mastodon.Pagination) ([]api.Item, error)
type apiFuncNotification func(nth []config.NotificationToHide, pg *mastodon.Pagination) ([]api.Item, error)
type apiEmptyFunc func() ([]api.Item, error)
type apiIDFunc func(pg *mastodon.Pagination, id mastodon.ID) ([]api.Item, error)
type apiIDFuncData func(pg *mastodon.Pagination, id mastodon.ID, data interface{}) ([]api.Item, error)
type apiSearchFunc func(search string) ([]api.Item, error)
type apiSearchPGFunc func(pg *mastodon.Pagination, search string) ([]api.Item, error)
type apiThreadFunc func(status *mastodon.Status) ([]api.Item, error)
type apiHistoryFunc func(status *mastodon.Status) ([]api.Item, error)

type LoadingLock struct {
	mux  sync.Mutex
	last time.Time
}

type DesktopNotificationType uint

const (
	DesktopNotificationNone DesktopNotificationType = iota
	DesktopNotificationFollower
	DesktopNotificationFavorite
	DesktopNotificationMention
	DesktopNotificationUpdate
	DesktopNotificationBoost
	DesktopNotificationPoll
	DesktopNotificationPost
)

type DesktopNotification struct {
	Type DesktopNotificationType
	Data string
}

type FeedUpdateFunc = func(from phony.Actor, notif DesktopNotification)

type Feed struct {
	phony.Inbox
	OnUpdate      FeedUpdateFunc
	accountClient *api.AccountClient
	config        *config.Config
	feedType      config.FeedType
	sticky        []api.Item
	items         []api.Item
	itemsMux      sync.RWMutex
	loadingNewer  *LoadingLock
	loadingOlder  *LoadingLock
	loadNewer     func(*Feed)
	loadOlder     func(*Feed)
	apiData       *api.RequestData
	apiDataMux    sync.Mutex
	receivers     []*api.Listener
	name          string
	hideBoosts    bool
	hideReplies   bool
}

func (f *Feed) Type() config.FeedType {
	return f.feedType
}

func (f *Feed) filteredList() []api.Item {
	f.itemsMux.RLock()
	defer f.itemsMux.RUnlock()
	filtered := []api.Item{}
	for _, fd := range f.items {
		switch x := fd.Raw().(type) {
		case *mastodon.Status:
			if f.Type() == config.TimelineHomeSpecial && x.Reblog == nil && x.InReplyToID == nil {
				continue
			}
			if x.Reblog != nil && f.hideBoosts {
				continue
			}
			if x.InReplyToID != nil && f.hideReplies {
				continue
			}
		}
		inUse, fType, _, _ := fd.Filtered(f.feedType)
		if inUse && fType == "hide" {
			continue
		}
		filtered = append(filtered, fd)
	}
	r := f.sticky
	return append(r, filtered...)
}

func (f *Feed) List() []api.Item {
	return f.filteredList()
}

func (f *Feed) Delete(id uint) {
	f.itemsMux.Lock()
	defer f.itemsMux.Unlock()
	var items []api.Item
	for _, item := range f.items {
		if item.ID() != id {
			items = append(items, item)
		}
	}
	f.items = items
	f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
}

func (f *Feed) Clear() {
	f.itemsMux.Lock()
	defer f.itemsMux.Unlock()
	f.items = []api.Item{}
	f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
}

func (f *Feed) Item(index int) (api.Item, error) {
	/*
		f.itemsMux.RLock()
		defer f.itemsMux.RUnlock()
		if f.StickyCount() > 0 && index < f.StickyCount() {
			return f.sticky[index], nil
		}
		if index < 0 || index >= len(f.items)+f.StickyCount() {
			return nil, errors.New("item out of range")
		}
	*/
	filtered := f.filteredList()
	if len(filtered) == 0 {
		return nil, errors.New("item out of range")
	}
	return filtered[index], nil
}

func (f *Feed) LoadNewer() {
	if f.loadNewer == nil {
		return
	}
	lock := f.loadingNewer.mux.TryLock()
	if !lock {
		return
	}
	if time.Since(f.loadingNewer.last) < (500 * time.Millisecond) {
		f.loadingNewer.mux.Unlock()
		return
	}
	f.loadNewer(f)
	f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
	f.loadingNewer.last = time.Now()
	f.loadingNewer.mux.Unlock()
}

func (f *Feed) LoadOlder() {
	if f.loadOlder == nil {
		return
	}
	lock := f.loadingOlder.mux.TryLock()
	if !lock {
		return
	}
	if time.Since(f.loadingOlder.last) < (500 * time.Microsecond) {
		f.loadingOlder.mux.Unlock()
		return
	}
	f.loadOlder(f)
	f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
	f.loadingOlder.last = time.Now()
	f.loadingOlder.mux.Unlock()
}

func (f *Feed) HasStream() bool {
	return len(f.receivers) > 0
}

func (f *Feed) Name() string {
	return f.name
}

func (f *Feed) StickyCount() int {
	return len(f.sticky)
}

func (f *Feed) singleNewerSearch(fn apiSearchFunc, search string) {
	items, err := fn(search)
	if err != nil {
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(items, f.items...)
		f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
	}
	f.itemsMux.Unlock()
}

func (f *Feed) singleThread(fn apiThreadFunc, status *mastodon.Status) {
	items, err := fn(status)
	if err != nil {
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(items, f.items...)
		f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
	}
	f.itemsMux.Unlock()
}

func (f *Feed) singleHistory(fn apiHistoryFunc, status *mastodon.Status) {
	items, err := fn(status)
	if err != nil {
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(items, f.items...)
		f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
	}
	f.itemsMux.Unlock()
}

func (f *Feed) normalNewer(fn apiFunc) {
	pg := mastodon.Pagination{}
	f.apiDataMux.Lock()
	if f.apiData.MinID != mastodon.ID("") {
		pg.MinID = f.apiData.MinID
	}
	items, err := fn(&pg)
	if err != nil {
		f.apiDataMux.Unlock()
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		switch item := items[0].Raw().(type) {
		case *mastodon.Status:
			f.apiData.MinID = item.ID
		case *api.NotificationData:
			f.apiData.MinID = item.Item.ID
		}
		if f.apiData.MaxID == mastodon.ID("") {
			switch item := items[len(items)-1].Raw().(type) {
			case *mastodon.Status:
				f.apiData.MaxID = item.ID
			case *api.NotificationData:
				f.apiData.MaxID = item.Item.ID
			}
		}
		f.items = append(items, f.items...)
		f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
	}
	f.itemsMux.Unlock()
	f.apiDataMux.Unlock()
}

func (f *Feed) normalOlder(fn apiFunc) {
	pg := mastodon.Pagination{}
	f.apiDataMux.Lock()
	if f.apiData.MaxID == mastodon.ID("") {
		f.apiDataMux.Unlock()
		f.loadNewer(f)
		return
	}
	pg.MaxID = f.apiData.MaxID
	items, err := fn(&pg)
	if err != nil {
		f.apiDataMux.Unlock()
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		switch item := items[len(items)-1].Raw().(type) {
		case *mastodon.Status:
			f.apiData.MaxID = item.ID
		case *api.NotificationData:
			f.apiData.MaxID = item.Item.ID
		}
		f.items = append(f.items, items...)
		f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
	}
	f.itemsMux.Unlock()
	f.apiDataMux.Unlock()
}

func (f *Feed) normalNewerNotification(fn apiFuncNotification, nth []config.NotificationToHide) {
	pg := mastodon.Pagination{}
	f.apiDataMux.Lock()
	if f.apiData.MinID != mastodon.ID("") {
		pg.MinID = f.apiData.MinID
	}
	items, err := fn(nth, &pg)
	if err != nil {
		f.apiDataMux.Unlock()
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		switch item := items[0].Raw().(type) {
		case *mastodon.Status:
			f.apiData.MinID = item.ID
		case *api.NotificationData:
			f.apiData.MinID = item.Item.ID
		}
		if f.apiData.MaxID == mastodon.ID("") {
			switch item := items[len(items)-1].Raw().(type) {
			case *mastodon.Status:
				f.apiData.MaxID = item.ID
			case *api.NotificationData:
				f.apiData.MaxID = item.Item.ID
			}
		}
		f.items = append(items, f.items...)
		f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
	}
	f.itemsMux.Unlock()
	f.apiDataMux.Unlock()
}

func (f *Feed) normalOlderNotification(fn apiFuncNotification, nth []config.NotificationToHide) {
	pg := mastodon.Pagination{}
	f.apiDataMux.Lock()
	if f.apiData.MaxID == mastodon.ID("") {
		f.apiDataMux.Unlock()
		f.loadNewer(f)
		return
	}
	pg.MaxID = f.apiData.MaxID
	items, err := fn(nth, &pg)
	if err != nil {
		f.apiDataMux.Unlock()
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		switch item := items[len(items)-1].Raw().(type) {
		case *mastodon.Status:
			f.apiData.MaxID = item.ID
		case *api.NotificationData:
			f.apiData.MaxID = item.Item.ID
		}
		f.items = append(f.items, items...)
		f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
	}
	f.itemsMux.Unlock()
	f.apiDataMux.Unlock()
}

func (f *Feed) newerSearchPG(fn apiSearchPGFunc, search string) {
	pg := mastodon.Pagination{}
	f.apiDataMux.Lock()
	if f.apiData.MinID != mastodon.ID("") {
		pg.MinID = f.apiData.MinID
	}
	items, err := fn(&pg, search)
	if err != nil {
		f.apiDataMux.Unlock()
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		item := items[0].Raw().(*mastodon.Status)
		f.apiData.MinID = item.ID
		f.items = append(items, f.items...)
		f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
		if f.apiData.MaxID == mastodon.ID("") {
			item = items[len(items)-1].Raw().(*mastodon.Status)
			f.apiData.MaxID = item.ID
		}
	}
	f.itemsMux.Unlock()
	f.apiDataMux.Unlock()
}

func (f *Feed) olderSearchPG(fn apiSearchPGFunc, search string) {
	pg := mastodon.Pagination{}
	f.apiDataMux.Lock()
	if f.apiData.MaxID == mastodon.ID("") {
		f.apiDataMux.Unlock()
		f.loadNewer(f)
		return
	}
	pg.MaxID = f.apiData.MaxID
	items, err := fn(&pg, search)
	if err != nil {
		f.apiDataMux.Unlock()
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		item := items[len(items)-1].Raw().(*mastodon.Status)
		f.apiData.MaxID = item.ID
		f.items = append(f.items, items...)
		f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
	}
	f.itemsMux.Unlock()
	f.apiDataMux.Unlock()
}

func (f *Feed) normalNewerUser(fn apiIDFunc, id mastodon.ID) {
	pg := mastodon.Pagination{}
	f.apiDataMux.Lock()
	if f.apiData.MinID != mastodon.ID("") {
		pg.MinID = f.apiData.MinID
	}
	items, err := fn(&pg, id)
	if err != nil {
		f.apiDataMux.Unlock()
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		item := items[0].Raw().(*mastodon.Status)
		f.apiData.MinID = item.ID
		f.items = append(items, f.items...)
		f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
		if f.apiData.MaxID == mastodon.ID("") {
			item = items[len(items)-1].Raw().(*mastodon.Status)
			f.apiData.MaxID = item.ID
		}
	}
	f.itemsMux.Unlock()
	f.apiDataMux.Unlock()
}

func (f *Feed) normalOlderUser(fn apiIDFunc, id mastodon.ID) {
	pg := mastodon.Pagination{}
	f.apiDataMux.Lock()
	if f.apiData.MaxID == mastodon.ID("") {
		f.apiDataMux.Unlock()
		f.loadNewer(f)
		return
	}
	pg.MaxID = f.apiData.MaxID
	items, err := fn(&pg, id)
	if err != nil {
		f.apiDataMux.Unlock()
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		item := items[len(items)-1].Raw().(*mastodon.Status)
		f.apiData.MaxID = item.ID
		f.items = append(f.items, items...)
		f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
	}
	f.itemsMux.Unlock()
	f.apiDataMux.Unlock()
}

func (f *Feed) normalNewerID(fn apiIDFunc, id mastodon.ID) {
	pg := mastodon.Pagination{}
	f.apiDataMux.Lock()
	if f.apiData.MinID != mastodon.ID("") {
		pg.MinID = f.apiData.MinID
	}
	items, err := fn(&pg, id)
	if err != nil {
		f.apiDataMux.Unlock()
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		item := items[0].Raw().(*mastodon.Status)
		f.apiData.MinID = item.ID
		f.items = append(items, f.items...)
		f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
		if f.apiData.MaxID == mastodon.ID("") {
			item = items[len(items)-1].Raw().(*mastodon.Status)
			f.apiData.MaxID = item.ID
		}
	}
	f.itemsMux.Unlock()
	f.apiDataMux.Unlock()
}

func (f *Feed) normalOlderID(fn apiIDFunc, id mastodon.ID) {
	pg := mastodon.Pagination{}
	f.apiDataMux.Lock()
	if f.apiData.MaxID == mastodon.ID("") {
		f.apiDataMux.Unlock()
		f.loadNewer(f)
		return
	}
	pg.MaxID = f.apiData.MaxID
	items, err := fn(&pg, id)
	if err != nil {
		f.apiDataMux.Unlock()
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		item := items[len(items)-1].Raw().(*mastodon.Status)
		f.apiData.MaxID = item.ID
		f.items = append(f.items, items...)
		f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
	}
	f.itemsMux.Unlock()
	f.apiDataMux.Unlock()
}

func (f *Feed) normalEmpty(fn apiEmptyFunc) {
	items, err := fn()
	if err != nil {
		return
	}
	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(f.items, items...)
		f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
	}
	f.itemsMux.Unlock()
}

func (f *Feed) linkNewer(fn apiFunc) {
	f.apiDataMux.Lock()
	pg := &mastodon.Pagination{}
	pg.MinID = f.apiData.MinID
	maxTmp := f.apiData.MaxID

	items, err := fn(pg)
	if err != nil {
		f.apiDataMux.Unlock()
		return
	}
	f.apiData.MinID = pg.MinID
	if pg.MaxID == "" {
		f.apiData.MaxID = maxTmp
	} else {
		f.apiData.MaxID = pg.MaxID
	}
	f.apiDataMux.Unlock()
	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(items, f.items...)
		f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
	}
	f.itemsMux.Unlock()
}

func (f *Feed) linkOlder(fn apiFunc) {
	f.apiDataMux.Lock()
	pg := &mastodon.Pagination{}
	pg.MaxID = f.apiData.MaxID
	if pg.MaxID == "" {
		f.apiDataMux.Unlock()
		return
	}

	items, err := fn(pg)
	if err != nil {
		f.apiDataMux.Unlock()
		return
	}
	f.apiData.MaxID = pg.MaxID
	f.apiDataMux.Unlock()

	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(f.items, items...)
		f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
	}
	f.itemsMux.Unlock()
}

func (f *Feed) linkNewerID(fn apiIDFunc, id mastodon.ID) {
	f.apiDataMux.Lock()
	pg := &mastodon.Pagination{}
	pg.MinID = f.apiData.MinID
	maxTmp := f.apiData.MaxID

	items, err := fn(pg, id)
	if err != nil {
		f.apiDataMux.Unlock()
		return
	}
	f.apiData.MinID = pg.MinID
	if pg.MaxID == "" {
		f.apiData.MaxID = maxTmp
	} else {
		f.apiData.MaxID = pg.MaxID
	}
	f.apiDataMux.Unlock()
	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(items, f.items...)
		f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
	}
	f.itemsMux.Unlock()
}

func (f *Feed) linkOlderID(fn apiIDFunc, id mastodon.ID) {
	f.apiDataMux.Lock()
	pg := &mastodon.Pagination{}
	pg.MaxID = f.apiData.MaxID
	if pg.MaxID == "" {
		f.apiDataMux.Unlock()
		return
	}

	items, err := fn(pg, id)
	if err != nil {
		f.apiDataMux.Unlock()
		return
	}
	f.apiData.MaxID = pg.MaxID
	f.apiDataMux.Unlock()

	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(f.items, items...)
		f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
	}
	f.itemsMux.Unlock()
}

func (f *Feed) linkNewerIDdata(fn apiIDFuncData, id mastodon.ID, data interface{}) {
	f.apiDataMux.Lock()
	pg := &mastodon.Pagination{}
	pg.MinID = f.apiData.MinID
	maxTmp := f.apiData.MaxID

	items, err := fn(pg, id, data)
	if err != nil {
		f.apiDataMux.Unlock()
		return
	}
	f.apiData.MinID = pg.MinID
	if pg.MaxID == "" {
		f.apiData.MaxID = maxTmp
	} else {
		f.apiData.MaxID = pg.MaxID
	}
	f.apiDataMux.Unlock()
	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(items, f.items...)
		f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
	}
	f.itemsMux.Unlock()
}

func (f *Feed) linkOlderIDdata(fn apiIDFuncData, id mastodon.ID, data interface{}) {
	f.apiDataMux.Lock()
	pg := &mastodon.Pagination{}
	pg.MaxID = f.apiData.MaxID
	if pg.MaxID == "" {
		f.apiDataMux.Unlock()
		return
	}

	items, err := fn(pg, id, data)
	if err != nil {
		f.apiDataMux.Unlock()
		return
	}
	f.apiData.MaxID = pg.MaxID
	f.apiDataMux.Unlock()

	f.itemsMux.Lock()
	if len(items) > 0 {
		f.items = append(f.items, items...)
		f.OnUpdate(f, DesktopNotification{Type: DesktopNotificationNone})
	}
	f.itemsMux.Unlock()
}

func (f *Feed) Close() {
	for _, r := range f.receivers {
		f.accountClient.RemoveListener(r)
	}
}

func (f *Feed) ListenToStream(id api.StreamID, handler HandlerFunc) {
	rec := &api.Listener{
		StreamID: id,
		OnEvent: func(from phony.Actor, e mastodon.Event) {
			handler(f, from, e)
		},
	}
	err := f.accountClient.AddListener(rec)
	if err != nil {
		log.Fatalln("Couldn't open stream")
	}
	f.receivers = append(f.receivers, rec)
}

// handles Conversations (DM) and Update (new posts)
func be_HandleUpdatesAndDM(f *Feed, from phony.Actor, e mastodon.Event) {
	f.Act(from, func() {
		switch t := e.(type) {
		case *mastodon.ConversationEvent:
			if t.Conversation.LastStatus == nil {
				return
			}
			s := api.NewStatusItem(t.Conversation.LastStatus, false)
			f.itemsMux.Lock()
			found := false
			if len(f.receivers) > 0 {
				for _, item := range f.items {
					switch v := item.Raw().(type) {
					case *mastodon.Status:
						if t.Conversation.LastStatus.ID == v.ID {
							found = true
							break
						}
					}
				}
			}
			if !found {
				f.items = append([]api.Item{s}, f.items...)
				f.OnUpdate(f, DesktopNotification{
					Type: DesktopNotificationMention,
				})
				f.apiData.MinID = t.Conversation.LastStatus.ID
			}
			f.itemsMux.Unlock()
		case *mastodon.UpdateEvent:
			s := api.NewStatusItem(t.Status, false)
			f.itemsMux.Lock()
			found := false
			if len(f.receivers) > 0 {
				for _, item := range f.items {
					switch v := item.Raw().(type) {
					case *mastodon.Status:
						if t.Status.ID == v.ID {
							found = true
							break
						}
					}
				}
			}
			if !found {
				f.items = append([]api.Item{s}, f.items...)
				f.OnUpdate(f, DesktopNotification{
					Type: DesktopNotificationPost,
				})
				f.apiData.MinID = t.Status.ID
			}
			f.itemsMux.Unlock()
		}
	})
}

func be_HandleNotifications(f *Feed, from phony.Actor, e mastodon.Event) {
	f.Act(from, func() {
		handleNotifications(f, e, false)
	})
}
func be_HandleNotificationsMentions(f *Feed, from phony.Actor, e mastodon.Event) {
	f.Act(from, func() {
		handleNotifications(f, e, true)
	})
}
func handleNotifications(f *Feed, e mastodon.Event, mentions bool) {
	switch t := e.(type) {
	case *mastodon.NotificationEvent:
		switch t.Notification.Type {
		case "follow":
			if slices.Contains(f.config.General.NotificationsToHide, config.HideFollow) || mentions {
				return
			}
		case "follow_request":
			if slices.Contains(f.config.General.NotificationsToHide, config.HideFollowRequest) || mentions {
				return
			}
		case "favourite":
			if slices.Contains(f.config.General.NotificationsToHide, config.HideFavorite) || mentions {
				return
			}
		case "reblog":
			if slices.Contains(f.config.General.NotificationsToHide, config.HideBoost) || mentions {
				return
			}
		case "mention":
			if slices.Contains(f.config.General.NotificationsToHide, config.HideMention) && !mentions {
				return
			}
		case "update":
			if slices.Contains(f.config.General.NotificationsToHide, config.HideEdited) || mentions {
				return
			}
		case "status":
			if slices.Contains(f.config.General.NotificationsToHide, config.HideStatus) || mentions {
				return
			}
		case "poll":
			if slices.Contains(f.config.General.NotificationsToHide, config.HidePoll) || mentions {
				return
			}
		}
		rel, err := f.accountClient.Client.GetAccountRelationships(context.Background(), []string{string(t.Notification.Account.ID)})
		if err != nil {
			return
		}
		if len(rel) == 0 {
			log.Fatalln(t.Notification.Account.Acct)
			return
		}
		s := api.NewNotificationItem(t.Notification,
			&api.User{
				Data:     &t.Notification.Account,
				Relation: rel[0],
			})
		f.itemsMux.Lock()
		f.items = append([]api.Item{s}, f.items...)
		nft := DesktopNotificationNone
		data := t.Notification.Account.DisplayName
		switch t.Notification.Type {
		case "follow", "follow_request":
			nft = DesktopNotificationFollower
		case "favourite":
			nft = DesktopNotificationFavorite
		case "reblog":
			nft = DesktopNotificationBoost
		case "mention":
			nft = DesktopNotificationMention
		case "update":
			nft = DesktopNotificationUpdate
		case "status":
			nft = DesktopNotificationPost
		case "poll":
			nft = DesktopNotificationPoll
		default:
			nft = DesktopNotificationNone
		}
		f.OnUpdate(f, DesktopNotification{
			Type: nft,
			Data: data,
		})
		f.itemsMux.Unlock()
	}
}

func newFeed(ac *api.AccountClient, ft config.FeedType, cnf *config.Config, hideBoosts bool, hideReplies bool) *Feed {
	feed := &Feed{
		accountClient: ac,
		config:        cnf,
		sticky:        make([]api.Item, 0),
		items:         make([]api.Item, 0),
		feedType:      ft,
		loadNewer:     func(*Feed) {},
		loadOlder:     func(*Feed) {},
		apiData:       &api.RequestData{},
		loadingNewer:  &LoadingLock{},
		loadingOlder:  &LoadingLock{},
		hideBoosts:    hideBoosts,
		hideReplies:   hideReplies,
	}
	runtime.SetFinalizer(feed, (*Feed).Close)
	return feed
}

type NewFeedFunc = func(ac *api.AccountClient, cnf *config.Config, hideBoosts bool, hideReplies bool) *Feed
type TimelineItemsFunc = func(ac *api.AccountClient, pg *mastodon.Pagination) ([]api.Item, error)
type HandlerFunc = func(f *Feed, from phony.Actor, e mastodon.Event)
type LoaderFuncs struct {
	newer func(*Feed)
	older func(*Feed)
}

func feedGenerator(load LoaderFuncs, feedType config.FeedType, streamId api.StreamID, handler HandlerFunc) NewFeedFunc {
	return func(ac *api.AccountClient, cnf *config.Config, hideBoosts bool, hideReplies bool) *Feed {
		feed := newFeed(ac, feedType, cnf, hideBoosts, hideReplies)
		feed.loadNewer = load.newer
		feed.loadOlder = load.older
		feed.ListenToStream(streamId, handler)
		return feed
	}
}

func loader_Timeline(getTimeline TimelineItemsFunc) LoaderFuncs {
	return LoaderFuncs{
		newer: func(feed *Feed) {
			getTlFn := func(pg *mastodon.Pagination) ([]api.Item, error) { return getTimeline(feed.accountClient, pg) }
			feed.normalNewer(getTlFn)
		},
		older: func(feed *Feed) {
			getTlFn := func(pg *mastodon.Pagination) ([]api.Item, error) { return getTimeline(feed.accountClient, pg) }
			feed.normalOlder(getTlFn)
		},
	}
}

var NewTimelineHome = feedGenerator(loader_Timeline((*api.AccountClient).GetTimeline), config.TimelineHome, api.MakeStreamID(api.HomeStream, ""), be_HandleUpdatesAndDM)
var NewTimelineHomeSpecial = feedGenerator(loader_Timeline((*api.AccountClient).GetTimeline), config.TimelineHomeSpecial, api.MakeStreamID(api.HomeStream, ""), be_HandleUpdatesAndDM)
var NewTimelineFederated = feedGenerator(loader_Timeline((*api.AccountClient).GetTimelineFederated), config.TimelineFederated, api.MakeStreamID(api.FederatedStream, ""), be_HandleUpdatesAndDM)
var NewTimelineLocal = feedGenerator(loader_Timeline((*api.AccountClient).GetTimelineLocal), config.TimelineLocal, api.MakeStreamID(api.LocalStream, ""), be_HandleUpdatesAndDM)
var NewConversations = feedGenerator(loader_Timeline((*api.AccountClient).GetConversations), config.Conversations, api.MakeStreamID(api.DirectStream, ""), be_HandleUpdatesAndDM)

func onlyShowMentions() []config.NotificationToHide {
	return []config.NotificationToHide{config.HideStatus, config.HideBoost, config.HideFollow, config.HideFollowRequest, config.HideFavorite, config.HidePoll, config.HideEdited}
}
func loader_Notifications(mentions bool) LoaderFuncs {
	return LoaderFuncs{
		newer: func(feed *Feed) {
			cnf := feed.config
			hide := cnf.General.NotificationsToHide
			if mentions {
				hide = onlyShowMentions()
			}
			feed.normalNewerNotification(feed.accountClient.GetNotifications, hide)
		},
		older: func(feed *Feed) {
			cnf := feed.config
			hide := cnf.General.NotificationsToHide
			if mentions {
				hide = onlyShowMentions()
			}
			feed.normalOlderNotification(feed.accountClient.GetNotifications, hide)
		},
	}
}

var NewNotifications = feedGenerator(loader_Notifications(false), config.Conversations, api.MakeStreamID(api.HomeStream, ""), be_HandleNotifications)
var NewNotificationsMentions = feedGenerator(loader_Notifications(true), config.Conversations, api.MakeStreamID(api.HomeStream, ""), be_HandleNotificationsMentions)

func NewFavorites(ac *api.AccountClient, cnf *config.Config) *Feed {
	feed := newFeed(ac, config.Favorited, cnf, false, false)
	feed.loadNewer = func(feed *Feed) { feed.linkNewer(feed.accountClient.GetFavorites) }
	feed.loadOlder = func(feed *Feed) { feed.linkOlder(feed.accountClient.GetFavorites) }

	return feed
}

func NewBookmarks(ac *api.AccountClient, cnf *config.Config) *Feed {
	feed := newFeed(ac, config.Saved, cnf, false, false)
	feed.loadNewer = func(feed *Feed) { feed.linkNewer(feed.accountClient.GetBookmarks) }
	feed.loadOlder = func(feed *Feed) { feed.linkOlder(feed.accountClient.GetBookmarks) }

	return feed
}

func NewUserSearch(ac *api.AccountClient, cnf *config.Config, search string) *Feed {
	feed := newFeed(ac, config.UserList, cnf, false, false)
	feed.name = search
	feed.loadNewer = func(feed *Feed) { feed.singleNewerSearch(feed.accountClient.GetUsers, search) }

	return feed
}

func NewUserProfile(ac *api.AccountClient, cnf *config.Config, user *api.User) *Feed {
	feed := newFeed(ac, config.User, cnf, false, false)
	feed.name = user.Data.Acct
	feed.sticky = append(feed.sticky, api.NewUserItem(user, true))
	pinned, err := ac.GetUserPinned(user.Data.ID)
	if err == nil {
		feed.sticky = append(feed.sticky, pinned...)
	}
	feed.loadNewer = func(feed *Feed) { feed.normalNewerUser(feed.accountClient.GetUser, user.Data.ID) }
	feed.loadOlder = func(feed *Feed) { feed.normalOlderUser(feed.accountClient.GetUser, user.Data.ID) }

	return feed
}

func NewThread(ac *api.AccountClient, cnf *config.Config, status *mastodon.Status) *Feed {
	feed := newFeed(ac, config.Thread, cnf, false, false)
	once := true
	feed.loadNewer = func(feed *Feed) {
		if once {
			feed.singleThread(feed.accountClient.GetThread, status)
			once = false
		}
	}

	return feed
}

func NewHistory(ac *api.AccountClient, cnf *config.Config, status *mastodon.Status) *Feed {
	feed := newFeed(ac, config.History, cnf, false, false)
	once := true
	feed.loadNewer = func(feed *Feed) {
		if once {
			feed.singleHistory(feed.accountClient.GetHistory, status)
			once = false
		}
	}
	return feed
}

func NewTag(ac *api.AccountClient, cnf *config.Config, search string, hideBoosts bool, hideReplies bool) *Feed {
	feed := newFeed(ac, config.Tag, cnf, hideBoosts, hideReplies)
	parts := strings.Split(search, " ")
	var tparts []string
	for _, p := range parts {
		p = strings.TrimPrefix(p, "#")
		if len(p) > 0 {
			tparts = append(tparts, p)
		}
	}
	joined := strings.Join(tparts, " ")
	feed.name = joined
	feed.loadNewer = func(feed *Feed) { feed.newerSearchPG(feed.accountClient.GetTagMultiple, joined) }
	feed.loadOlder = func(feed *Feed) { feed.olderSearchPG(feed.accountClient.GetTagMultiple, joined) }
	for _, t := range tparts {
		feed.ListenToStream(api.MakeStreamID(api.TagStream, t), be_HandleUpdatesAndDM)
	}
	return feed
}

func NewTags(ac *api.AccountClient, cnf *config.Config) *Feed {
	feed := newFeed(ac, config.Tags, cnf, false, false)
	once := true
	feed.loadNewer = func(feed *Feed) {
		if once {
			feed.normalNewer(feed.accountClient.GetTags)
		}
		once = false
	}
	feed.loadOlder = func(feed *Feed) { feed.normalOlder(feed.accountClient.GetTags) }

	return feed
}

func NewListList(ac *api.AccountClient, cnf *config.Config) *Feed {
	feed := newFeed(ac, config.Lists, cnf, false, false)
	once := true
	feed.loadNewer = func(feed *Feed) {
		if once {
			feed.normalEmpty(feed.accountClient.GetLists)
		}
		once = false
	}

	return feed
}

func NewList(ac *api.AccountClient, cnf *config.Config, list *mastodon.List, hideBoosts bool, hideReplies bool) *Feed {
	feed := newFeed(ac, config.List, cnf, hideBoosts, hideReplies)
	feed.name = list.Title
	feed.loadNewer = func(feed *Feed) { feed.normalNewerID(feed.accountClient.GetListStatuses, list.ID) }
	feed.loadOlder = func(feed *Feed) { feed.normalOlderID(feed.accountClient.GetListStatuses, list.ID) }
	feed.ListenToStream(api.MakeStreamID(api.ListStream, string(list.ID)), be_HandleUpdatesAndDM)
	return feed
}

func NewUsersInList(ac *api.AccountClient, cnf *config.Config, list *mastodon.List) *Feed {
	feed := newFeed(ac, config.ListUsersIn, cnf, false, false)
	feed.name = list.Title
	once := true
	feed.loadNewer = func(feed *Feed) {
		if once {
			feed.linkNewerIDdata(feed.accountClient.GetListUsers, list.ID, list)
		}
		once = false
	}
	feed.loadOlder = func(feed *Feed) { feed.linkOlderIDdata(feed.accountClient.GetListUsers, list.ID, list) }

	return feed
}

func NewUsersAddList(ac *api.AccountClient, cnf *config.Config, list *mastodon.List) *Feed {
	feed := newFeed(ac, config.ListUsersAdd, cnf, false, false)
	feed.name = list.Title
	once := true
	feed.loadNewer = func(feed *Feed) {
		if once {
			feed.linkNewerIDdata(feed.accountClient.GetFollowingForList, ac.Me.ID, list)
		}
		once = false
	}
	feed.loadOlder = func(feed *Feed) { feed.linkOlderIDdata(feed.accountClient.GetFollowingForList, ac.Me.ID, list) }

	return feed
}

func NewFavoritesStatus(ac *api.AccountClient, cnf *config.Config, id mastodon.ID) *Feed {
	feed := newFeed(ac, config.Favorites, cnf, false, false)
	once := true
	feed.loadNewer = func(feed *Feed) {
		if once {
			feed.linkNewerID(feed.accountClient.GetFavoritesStatus, id)
		}
		once = false
	}

	return feed
}

func NewBoosts(ac *api.AccountClient, cnf *config.Config, id mastodon.ID) *Feed {
	feed := newFeed(ac, config.Boosts, cnf, false, false)
	once := true
	feed.loadNewer = func(feed *Feed) {
		if once {
			feed.linkNewerID(feed.accountClient.GetBoostsStatus, id)
		}
		once = false
	}

	return feed
}

func NewFollowers(ac *api.AccountClient, cnf *config.Config, id mastodon.ID) *Feed {
	feed := newFeed(ac, config.Followers, cnf, false, false)
	once := true
	feed.loadNewer = func(feed *Feed) {
		if once {
			feed.linkNewerID(feed.accountClient.GetFollowers, id)
		}
		once = false
	}
	feed.loadOlder = func(feed *Feed) { feed.linkOlderID(feed.accountClient.GetFollowers, id) }

	return feed
}

func NewFollowing(ac *api.AccountClient, cnf *config.Config, id mastodon.ID) *Feed {
	feed := newFeed(ac, config.Following, cnf, false, false)
	once := true
	feed.loadNewer = func(feed *Feed) {
		if once {
			feed.linkNewerID(feed.accountClient.GetFollowing, id)
		}
		once = false
	}
	feed.loadOlder = func(feed *Feed) { feed.linkOlderID(feed.accountClient.GetFollowing, id) }

	return feed
}

func NewBlocking(ac *api.AccountClient, cnf *config.Config) *Feed {
	feed := newFeed(ac, config.Blocking, cnf, false, false)
	once := true
	feed.loadNewer = func(feed *Feed) {
		if once {
			feed.linkNewer(feed.accountClient.GetBlocking)
		}
		once = false
	}
	feed.loadOlder = func(feed *Feed) { feed.linkOlder(feed.accountClient.GetBlocking) }

	return feed
}

func NewMuting(ac *api.AccountClient, cnf *config.Config) *Feed {
	feed := newFeed(ac, config.Muting, cnf, false, false)
	once := true
	feed.loadNewer = func(feed *Feed) {
		if once {
			feed.linkNewer(feed.accountClient.GetMuting)
		}
		once = false
	}
	feed.loadOlder = func(feed *Feed) { feed.linkOlder(feed.accountClient.GetMuting) }

	return feed
}

func NewFollowRequests(ac *api.AccountClient, cnf *config.Config) *Feed {
	feed := newFeed(ac, config.FollowRequests, cnf, false, false)
	once := true
	feed.loadNewer = func(feed *Feed) {
		if once {
			feed.linkNewer(feed.accountClient.GetFollowRequests)
		}
		once = false
	}
	feed.loadOlder = func(feed *Feed) { feed.linkOlder(feed.accountClient.GetFollowRequests) }

	return feed
}
