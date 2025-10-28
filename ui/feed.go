package ui

import (
	"fmt"
	"strconv"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/RasmusLindroth/tut/api"
	"github.com/RasmusLindroth/tut/config"
	"github.com/RasmusLindroth/tut/feed"
	"github.com/RasmusLindroth/tut/util"
)

type FeedList struct {
	Text        *tview.List
	Symbol      *tview.List
	stickyCount int
}

func (fl *FeedList) InFocus(style config.Style) {
	inFocus(fl.Text, style)
	inFocus(fl.Symbol, style)
}

func inFocus(l *tview.List, style config.Style) {
	l.SetBackgroundColor(style.Background)
	l.SetMainTextColor(style.Text)
	l.SetSelectedBackgroundColor(style.ListSelectedBackground)
	l.SetSelectedTextColor(style.ListSelectedText)
}

func (fl *FeedList) OutFocus(style config.Style) {
	outFocus(fl.Text, style)
	outFocus(fl.Symbol, style)
}

func outFocus(l *tview.List, style config.Style) {
	l.SetBackgroundColor(style.Background)
	l.SetMainTextColor(style.Text)
	l.SetSelectedBackgroundColor(style.ListSelectedInactiveBackground)
	l.SetSelectedTextColor(style.ListSelectedInactiveText)
}

type Feed struct {
	tutView  *TutView
	Data     *feed.Feed
	List     *FeedList
	Content  *FeedContent
	Timeline *config.Timeline
}

func (f *Feed) ListInFocus() {
	f.List.InFocus(f.tutView.tut.Config.Style)
}

func (f *Feed) ListOutFocus() {
	f.List.OutFocus(f.tutView.tut.Config.Style)
}

func (f *Feed) LoadOlder() {
	f.Data.LoadOlder()
}

func (f *Feed) LoadNewer(force bool) {
	if f.Data.HasStream() && !force {
		return
	}
	f.Data.LoadNewer()
}

func (f *Feed) Delete() {
	id := f.List.GetCurrentID()
	f.Data.Delete(id)
}

func (f *Feed) DrawContent() {
	id := f.List.GetCurrentID()
	for _, item := range f.Data.List() {
		if id != item.ID() {
			continue
		}
		DrawItem(f.tutView, item, f.Content.Main, f.Content.Controls, f.Data.Type())
		f.tutView.ShouldSync()
	}
}

func (ui *Feed) CreateUpdateCallback() feed.FeedUpdateCallBack {
	return func(logical *feed.Feed, nft feed.DesktopNotification) {
		feed.SendDesktopNotification(ui.tutView.tut.Config.NotificationConfig, nft)
		ui.tutView.tut.App.QueueUpdateDraw(func() {
			lLen := ui.List.GetItemCount()
			curr := ui.List.GetCurrentID()
			ui.List.Clear()
			for _, item := range ui.Data.List() {
				main, symbol := DrawListItem(ui.tutView.tut.Config, item)
				ui.List.AddItem(main, symbol, item.ID())
			}
			if ui.tutView.tut.Config.General.StickToTop {
				ui.List.SetCurrentItem(ui.List.stickyCount)
				ui.DrawContent()
			} else {
				ui.List.SetByID(curr)
			}
			if lLen == 0 {
				ui.DrawContent()
			}
		})
	}
}

func newPublicFeed(gen feed.NewFeedFunc) func(tv *TutView, tl *config.Timeline) *Feed {
	return func(tv *TutView, tl *config.Timeline) *Feed {
		fd := &Feed{}
		f := gen(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback() /* not called until call to f.LoadNewer() */, tl.HideBoosts, tl.HideReplies)
		*fd = Feed{
			tutView:  tv,
			Data:     f,
			Timeline: tl,
			List:     NewFeedList(tv.tut, f.StickyCount()),
			Content:  NewFeedContent(tv.tut),
		}
		f.LoadNewer()

		return fd
	}
}

var (
	NewHomeFeed         = newPublicFeed(feed.NewTimelineHome)
	NewHomeSpecialFeed  = newPublicFeed(feed.NewTimelineHomeSpecial)
	NewFederatedFeed    = newPublicFeed(feed.NewTimelineFederated)
	NewLocalFeed        = newPublicFeed(feed.NewTimelineLocal)
	NewNotificationFeed = newPublicFeed(feed.NewNotifications)
	NewMentionFeed      = newPublicFeed(feed.NewNotificationsMentions)
)

func NewThreadFeed(tv *TutView, item api.Item, tl *config.Timeline) *Feed {
	status := util.StatusOrReblog(item.Raw().(*mastodon.Status))
	fd := &Feed{
		tutView:  tv,
		Content:  NewFeedContent(tv.tut),
		Timeline: tl,
	}
	f := feed.NewThread(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback(), status)
	fd.Data = f
	fd.List = NewFeedList(tv.tut, f.StickyCount())
	f.LoadNewer()
	for i, s := range f.List() {
		main, symbol := DrawListItem(tv.tut.Config, s)
		fd.List.AddItem(main, symbol, s.ID())
		if s.Raw().(*mastodon.Status).ID == item.Raw().(*mastodon.Status).ID {
			fd.List.SetCurrentItem(i)
		}
	}
	fd.DrawContent()

	return fd
}

func NewHistoryFeed(tv *TutView, item api.Item, tl *config.Timeline) *Feed {
	status := util.StatusOrReblog(item.Raw().(*mastodon.Status))
	fd := &Feed{
		tutView:  tv,
		Content:  NewFeedContent(tv.tut),
		Timeline: tl,
	}
	f := feed.NewHistory(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback(), status)
	fd.Data = f
	fd.List = NewFeedList(tv.tut, f.StickyCount())
	f.LoadNewer()
	for _, s := range f.List() {
		main, symbol := DrawListItem(tv.tut.Config, s)
		fd.List.AddItem(main, symbol, s.ID())
	}
	fd.List.SetCurrentItem(0)
	fd.DrawContent()

	return fd
}

func NewConversationsFeed(tv *TutView, tl *config.Timeline) *Feed {
	fd := &Feed{
		tutView:  tv,
		Content:  NewFeedContent(tv.tut),
		Timeline: tl,
	}
	f := feed.NewConversations(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback(), false, false)
	fd.Data = f
	fd.List = NewFeedList(tv.tut, f.StickyCount())
	f.LoadNewer()

	return fd
}

func NewUserFeed(tv *TutView, item api.Item, tl *config.Timeline) *Feed {
	if item.Type() != api.UserType && item.Type() != api.ProfileType {
		panic("Can't open user. Wrong type.\n")
	}
	u := item.Raw().(*api.User)
	fd := &Feed{
		tutView:  tv,
		Content:  NewFeedContent(tv.tut),
		Timeline: tl,
	}
	f := feed.NewUserProfile(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback(), u)
	fd.Data = f
	fd.List = NewFeedList(tv.tut, f.StickyCount())
	f.LoadNewer()

	return fd
}

func NewUserSearchFeed(tv *TutView, tl *config.Timeline) *Feed {
	fd := &Feed{
		tutView:  tv,
		Content:  NewFeedContent(tv.tut),
		Timeline: tl,
	}
	f := feed.NewUserSearch(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback(), tl.Subaction)
	fd.Data = f
	fd.List = NewFeedList(tv.tut, f.StickyCount())
	f.LoadNewer()
	for _, s := range f.List() {
		main, symbol := DrawListItem(tv.tut.Config, s)
		fd.List.AddItem(main, symbol, s.ID())
	}
	fd.DrawContent()

	return fd
}

func NewTagFeed(tv *TutView, tl *config.Timeline) *Feed {
	fd := &Feed{
		tutView:  tv,
		Content:  NewFeedContent(tv.tut),
		Timeline: tl,
	}
	f := feed.NewTag(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback(), tl.Subaction, tl.HideBoosts, tl.HideReplies)
	fd.Data = f
	fd.List = NewFeedList(tv.tut, f.StickyCount())
	f.LoadNewer()

	return fd
}

func NewTagsFeed(tv *TutView, tl *config.Timeline) *Feed {
	fd := &Feed{
		tutView:  tv,
		Content:  NewFeedContent(tv.tut),
		Timeline: tl,
	}
	f := feed.NewTags(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback())
	fd.Data = f
	fd.List = NewFeedList(tv.tut, f.StickyCount())
	f.LoadNewer()

	return fd
}

func NewListsFeed(tv *TutView, tl *config.Timeline) *Feed {
	fd := &Feed{
		tutView:  tv,
		Content:  NewFeedContent(tv.tut),
		Timeline: tl,
	}
	f := feed.NewListList(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback())
	fd.Data = f
	fd.List = NewFeedList(tv.tut, f.StickyCount())
	f.LoadNewer()

	return fd
}

func NewListFeed(tv *TutView, l *mastodon.List, tl *config.Timeline) *Feed {
	fd := &Feed{
		tutView:  tv,
		Content:  NewFeedContent(tv.tut),
		Timeline: tl,
	}
	f := feed.NewList(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback(), l, tl.HideBoosts, tl.HideReplies)
	fd.Data = f
	fd.List = NewFeedList(tv.tut, f.StickyCount())
	f.LoadNewer()

	return fd
}

func NewUsersInListFeed(tv *TutView, l *mastodon.List, tl *config.Timeline) *Feed {
	fd := &Feed{
		tutView:  tv,
		Content:  NewFeedContent(tv.tut),
		Timeline: tl,
	}
	f := feed.NewUsersInList(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback(), l)
	fd.Data = f
	fd.List = NewFeedList(tv.tut, f.StickyCount())
	f.LoadNewer()

	return fd
}

func NewUsersAddListFeed(tv *TutView, l *mastodon.List, tl *config.Timeline) *Feed {
	fd := &Feed{
		tutView:  tv,
		Content:  NewFeedContent(tv.tut),
		Timeline: tl,
	}
	f := feed.NewUsersAddList(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback(), l)
	fd.Data = f
	fd.List = NewFeedList(tv.tut, f.StickyCount())
	f.LoadNewer()

	return fd
}

func NewFavoritedFeed(tv *TutView, tl *config.Timeline) *Feed {
	fd := &Feed{
		tutView:  tv,
		Content:  NewFeedContent(tv.tut),
		Timeline: tl,
	}
	f := feed.NewFavorites(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback())
	fd.Data = f
	fd.List = NewFeedList(tv.tut, f.StickyCount())
	f.LoadNewer()

	return fd
}

func NewBookmarksFeed(tv *TutView, tl *config.Timeline) *Feed {
	fd := &Feed{
		tutView:  tv,
		Content:  NewFeedContent(tv.tut),
		Timeline: tl,
	}
	f := feed.NewBookmarks(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback())
	fd.Data = f
	fd.List = NewFeedList(tv.tut, f.StickyCount())
	f.LoadNewer()

	return fd
}

func NewFavoritesStatus(tv *TutView, id mastodon.ID, tl *config.Timeline) *Feed {
	fd := &Feed{
		tutView:  tv,
		Content:  NewFeedContent(tv.tut),
		Timeline: tl,
	}
	f := feed.NewFavoritesStatus(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback(), id)
	fd.Data = f
	fd.List = NewFeedList(tv.tut, f.StickyCount())
	f.LoadNewer()

	return fd
}

func NewBoosts(tv *TutView, id mastodon.ID, tl *config.Timeline) *Feed {
	fd := &Feed{
		tutView:  tv,
		Content:  NewFeedContent(tv.tut),
		Timeline: tl,
	}
	f := feed.NewBoosts(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback(), id)
	fd.Data = f
	fd.List = NewFeedList(tv.tut, f.StickyCount())
	f.LoadNewer()

	return fd
}

func NewFollowers(tv *TutView, id mastodon.ID, tl *config.Timeline) *Feed {
	fd := &Feed{
		tutView:  tv,
		Content:  NewFeedContent(tv.tut),
		Timeline: tl,
	}
	f := feed.NewFollowers(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback(), id)
	fd.Data = f
	fd.List = NewFeedList(tv.tut, f.StickyCount())
	f.LoadNewer()

	return fd
}

func NewFollowing(tv *TutView, id mastodon.ID, tl *config.Timeline) *Feed {
	fd := &Feed{
		tutView:  tv,
		Content:  NewFeedContent(tv.tut),
		Timeline: tl,
	}
	f := feed.NewFollowing(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback(), id)
	fd.Data = f
	fd.List = NewFeedList(tv.tut, f.StickyCount())
	f.LoadNewer()

	return fd
}

func NewBlocking(tv *TutView, tl *config.Timeline) *Feed {
	fd := &Feed{
		tutView:  tv,
		Content:  NewFeedContent(tv.tut),
		Timeline: tl,
	}
	f := feed.NewBlocking(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback())
	fd.Data = f
	fd.List = NewFeedList(tv.tut, f.StickyCount())
	f.LoadNewer()

	return fd
}

func NewMuting(tv *TutView, tl *config.Timeline) *Feed {
	fd := &Feed{
		tutView:  tv,
		Content:  NewFeedContent(tv.tut),
		Timeline: tl,
	}
	f := feed.NewMuting(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback())
	fd.Data = f
	fd.List = NewFeedList(tv.tut, f.StickyCount())
	f.LoadNewer()

	return fd
}

func NewFollowRequests(tv *TutView, tl *config.Timeline) *Feed {
	fd := &Feed{
		tutView:  tv,
		Content:  NewFeedContent(tv.tut),
		Timeline: tl,
	}
	f := feed.NewFollowRequests(tv.tut.Client, tv.tut.Config, fd.CreateUpdateCallback())
	fd.Data = f
	fd.List = NewFeedList(tv.tut, f.StickyCount())
	f.LoadNewer()

	return fd
}

func NewFeedList(t *Tut, stickyCount int) *FeedList {
	fl := &FeedList{
		Text:        NewList(t.Config),
		Symbol:      NewList(t.Config),
		stickyCount: stickyCount,
	}
	return fl
}

func (fl *FeedList) AddItem(text string, symbols string, id uint) {
	fl.Text.AddItem(text, fmt.Sprintf("%d", id), 0, nil)
	fl.Symbol.AddItem(symbols, fmt.Sprintf("%d", id), 0, nil)
}

func (fl *FeedList) Set(index int) (loadOlder bool, loadNewer bool) {
	ni := index
	if ni >= fl.Text.GetItemCount() {
		ni = fl.Text.GetItemCount() - 1
	}
	if ni < 0 {
		ni = 0
	}
	fl.Text.SetCurrentItem(ni)
	fl.Symbol.SetCurrentItem(ni)
	return fl.Text.GetItemCount()-(ni+1) < 5, ni-fl.stickyCount < 4
}

func (fl *FeedList) Next() (loadOlder bool) {
	ni := fl.Text.GetCurrentItem() + 1
	if ni >= fl.Text.GetItemCount() {
		ni = fl.Text.GetItemCount() - 1
		if ni < 0 {
			ni = 0
		}
	}
	fl.Text.SetCurrentItem(ni)
	fl.Symbol.SetCurrentItem(ni)
	return fl.Text.GetItemCount()-(ni+1) < 5
}

func (fl *FeedList) Prev() (loadNewer bool) {
	ni := fl.Text.GetCurrentItem() - 1
	if ni < 0 {
		ni = 0
	}
	fl.Text.SetCurrentItem(ni)
	fl.Symbol.SetCurrentItem(ni)
	return ni-fl.stickyCount < 4
}

func (fl *FeedList) Clear() {
	fl.Text.Clear()
	fl.Symbol.Clear()
}

func (fl *FeedList) GetItemCount() int {
	return fl.Text.GetItemCount()
}

func (fl *FeedList) SetCurrentItem(index int) {
	fl.Text.SetCurrentItem(index)
	fl.Symbol.SetCurrentItem(index)
}

func (fl *FeedList) GetCurrentID() uint {
	if fl.GetItemCount() == 0 {
		return 0
	}
	i := fl.Text.GetCurrentItem()
	_, sec := fl.Text.GetItemText(i)
	id, err := strconv.ParseUint(sec, 10, 32)
	if err != nil {
		return 0
	}
	return uint(id)
}

func (fl *FeedList) SetByID(id uint) {
	if fl.Text.GetItemCount() == 0 {
		return
	}
	s := fmt.Sprintf("%d", id)
	items := fl.Text.FindItems("", s, false, false)
	for _, i := range items {
		_, sec := fl.Text.GetItemText(i)
		if sec == s {
			fl.Text.SetCurrentItem(i)
			fl.Symbol.SetCurrentItem(i)
			break
		}
	}
}

type FeedContent struct {
	Main     *tview.TextView
	Controls *tview.Flex
}

func NewFeedContent(t *Tut) *FeedContent {
	m := NewTextView(t.Config)
	m.SetWordWrap(true)

	if t.Config.General.MaxWidth > 0 {
		mw := t.Config.General.MaxWidth
		m.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
			rWidth := width
			if rWidth > mw {
				rWidth = mw
			}
			return x, y, rWidth, height
		})
	}
	c := NewControlView(t.Config)
	fc := &FeedContent{
		Main:     m,
		Controls: c,
	}
	return fc
}
