package config

type ConfigTOML struct {
	General            GeneralTOML       `toml:"general"`
	Style              StyleTOML         `toml:"style"`
	Media              MediaTOML         `toml:"media"`
	OpenPattern        OpenPatternTOML   `toml:"open-pattern"`
	OpenCustom         OpenCustomTOML    `toml:"open-custom"`
	NotificationConfig NotificationsTOML `toml:"desktop-notification"`
	Input              InputTOML         `toml:"input"`
}

type GeneralTOML struct {
	Editor              *string             `toml:"editor"`
	Confirmation        *bool               `toml:"confirmation"`
	MouseSupport        *bool               `toml:"mouse-support"`
	DateFormat          *string             `toml:"date-format"`
	DateTodayFormat     *string             `toml:"date-today-format"`
	DateRelative        *int                `toml:"date-relative"`
	MaxWidth            *int                `toml:"max-width"`
	QuoteReply          *bool               `toml:"quote-reply"`
	ShortHints          *bool               `toml:"short-hints"`
	ShowFilterPhrase    *bool               `toml:"show-filter-phrase"`
	ListPlacement       *string             `toml:"list-placement"`
	ListSplit           *string             `toml:"list-split"`
	ListProportion      *int                `toml:"list-proportion"`
	ContentProportion   *int                `toml:"content-proportion"`
	TerminalTitle       *int                `toml:"terminal-title"`
	ShowIcons           *bool               `toml:"show-icons"`
	ShowHelp            *bool               `toml:"show-help"`
	RedrawUI            *bool               `toml:"redraw-ui"`
	LeaderKey           *string             `toml:"leader-key"`
	LeaderTimeout       *int64              `toml:"leader-timeout"`
	Timelines           *[]TimelineTOML     `toml:"timelines"`
	LeaderActions       *[]LeaderActionTOML `toml:"leader-actions"`
	StickToTop          *bool               `toml:"stick-to-top"`
	NotificationsToHide *[]string           `toml:"notifications-to-hide"`
	ShowBoostedUser     *bool               `toml:"show-boosted-user"`
	DynamicTimelineName *bool               `toml:"dynamic-timeline-name"`
	CommandsInNewPane   *bool               `toml:"commands-in-new-pane"`
}

type TimelineTOML struct {
	Name        *string   `toml:"name"`
	Type        *string   `toml:"type"`
	Data        *string   `toml:"data"`
	Keys        *[]string `toml:"keys"`
	SpecialKeys *[]string `toml:"special-keys"`
	Shortcut    *string   `toml:"shortcut"`
	HideBoosts  *bool     `toml:"hide-boosts"`
	HideReplies *bool     `toml:"hide-replies"`

	Closed           *bool   `toml:"closed"`
	OnCreationClosed *string `toml:"on-creation-closed"`
	OnFocus          *string `toml:"on-focus"`
}

type LeaderActionTOML struct {
	Type     *string `toml:"type"`
	Data     *string `toml:"data"`
	Shortcut *string `toml:"shortcut"`
}

type StyleTOML struct {
	Theme *string `toml:"theme"`

	XrdbPrefix *string `toml:"xrdb-prefix"`

	Background *string `toml:"background"`
	Text       *string `toml:"text"`

	Subtle      *string `toml:"subtle"`
	WarningText *string `toml:"warning-text"`

	TextSpecial1 *string `toml:"text-special-one"`
	TextSpecial2 *string `toml:"text-special-two"`

	TopBarBackground *string `toml:"top-bar-background"`
	TopBarText       *string `toml:"top-bar-text"`

	StatusBarBackground *string `toml:"status-bar-background"`
	StatusBarText       *string `toml:"status-bar-text"`

	StatusBarViewBackground *string `toml:"status-bar-view-background"`
	StatusBarViewText       *string `toml:"status-bar-view-text"`

	ListSelectedBackground *string `toml:"list-selected-background"`
	ListSelectedText       *string `toml:"list-selected-text"`

	ListSelectedInactiveBackground *string `toml:"list-selected-inactive-background"`
	ListSelectedInactiveText       *string `toml:"list-selected-inactive-text"`

	ControlsText      *string `toml:"controls-text"`
	ControlsHighlight *string `toml:"controls-highlight"`

	AutocompleteBackground *string `toml:"autocomplete-background"`
	AutocompleteText       *string `toml:"autocomplete-text"`

	AutocompleteSelectedBackground *string `toml:"autocomplete-selected-background"`
	AutocompleteSelectedText       *string `toml:"autocomplete-selected-text"`

	ButtonColorOne *string `toml:"button-color-one"`
	ButtonColorTwo *string `toml:"button-color-two"`

	TimelineNameBackground *string `toml:"timeline-name-background"`
	TimelineNameText       *string `toml:"timeline-name-text"`

	IconColor *string `toml:"icon-color"`

	CommandText *string `toml:"command-text"`
}

type ViewerTOML struct {
	Program  *string `toml:"program"`
	Args     *string `toml:"args"`
	Terminal *bool   `toml:"terminal"`
	Single   *bool   `toml:"single"`
	Reverse  *bool   `toml:"reverse"`
}

type MediaTOML struct {
	DeleteTmpFiles *bool       `toml:"delete-temp-files"`
	Image          *ViewerTOML `toml:"image"`
	Video          *ViewerTOML `toml:"video"`
	Audio          *ViewerTOML `toml:"audio"`
	Link           *ViewerTOML `toml:"link"`
}

type PatternTOML struct {
	Matching *string `toml:"matching"`
	Program  *string `toml:"program"`
	Args     *string `toml:"args"`
	Terminal *bool   `toml:"terminal"`
}

type OpenPatternTOML struct {
	Patterns *[]PatternTOML `toml:"patterns"`
}

type CustomTOML struct {
	Program     *string   `toml:"program"`
	Args        *string   `toml:"args"`
	Terminal    *bool     `toml:"terminal"`
	Hint        *string   `toml:"hint"`
	Keys        *[]string `toml:"keys"`
	SpecialKeys *[]string `toml:"special-keys"`
}

type OpenCustomTOML struct {
	Programs *[]CustomTOML `toml:"programs"`
}

type NotificationsTOML struct {
	Followers *bool `toml:"followers"`
	Favorite  *bool `toml:"favorite"`
	Mention   *bool `toml:"mention"`
	Update    *bool `toml:"update"`
	Boost     *bool `toml:"boost"`
	Poll      *bool `toml:"poll"`
	Posts     *bool `toml:"posts"`
}

type InputActionTOML struct {
	// key description
	Description *string   `toml:"desc"`
	Hint        *string   `toml:"hint"`
	HintAlt     *string   `toml:"hint-alt"`
	Keys        *[]string `toml:"keys"`
	SpecialKeys *[]string `toml:"special-keys"`
}

type InputTOML struct {
	GlobalDown  *InputActionTOML `toml:"global-down"`
	GlobalUp    *InputActionTOML `toml:"global-up"`
	GlobalEnter *InputActionTOML `toml:"global-enter"`
	GlobalBack  *InputActionTOML `toml:"global-back"`
	GlobalExit  *InputActionTOML `toml:"global-exit"`

	MainHome        *InputActionTOML `toml:"main-home"`
	MainEnd         *InputActionTOML `toml:"main-end"`
	MainPrevFeed    *InputActionTOML `toml:"main-prev-feed"`
	MainNextFeed    *InputActionTOML `toml:"main-next-feed"`
	MainPrevPane    *InputActionTOML `toml:"main-prev-pane"`
	MainNextPane    *InputActionTOML `toml:"main-next-pane"`
	MainCompose     *InputActionTOML `toml:"main-compose"`
	MainNextAccount *InputActionTOML `toml:"main-next-account"`
	MainPrevAccount *InputActionTOML `toml:"main-prev-account"`

	StatusAvatar       *InputActionTOML `toml:"status-avatar"`
	StatusBoost        *InputActionTOML `toml:"status-boost"`
	StatusDelete       *InputActionTOML `toml:"status-delete"`
	StatusEdit         *InputActionTOML `toml:"status-edit"`
	StatusFavorite     *InputActionTOML `toml:"status-favorite"`
	StatusMedia        *InputActionTOML `toml:"status-media"`
	StatusLinks        *InputActionTOML `toml:"status-links"`
	StatusPoll         *InputActionTOML `toml:"status-poll"`
	StatusReply        *InputActionTOML `toml:"status-reply"`
	StatusBookmark     *InputActionTOML `toml:"status-bookmark"`
	StatusThread       *InputActionTOML `toml:"status-thread"`
	StatusUser         *InputActionTOML `toml:"status-user"`
	StatusViewFocus    *InputActionTOML `toml:"status-view-focus"`
	StatusYank         *InputActionTOML `toml:"status-yank"`
	StatusToggleCW     *InputActionTOML `toml:"status-toggle-cw"`
	StatusShowFiltered *InputActionTOML `toml:"status-show-filtered"`

	UserAvatar              *InputActionTOML `toml:"user-avatar"`
	UserBlock               *InputActionTOML `toml:"user-block"`
	UserFollow              *InputActionTOML `toml:"user-follow"`
	UserFollowRequestDecide *InputActionTOML `toml:"user-follow-request-decide"`
	UserMute                *InputActionTOML `toml:"user-mute"`
	UserLinks               *InputActionTOML `toml:"user-links"`
	UserUser                *InputActionTOML `toml:"user-user"`
	UserViewFocus           *InputActionTOML `toml:"user-view-focus"`
	UserYank                *InputActionTOML `toml:"user-yank"`

	ListOpenFeed   *InputActionTOML `toml:"list-open-feed"`
	ListUserList   *InputActionTOML `toml:"list-user-list"`
	ListUserAdd    *InputActionTOML `toml:"list-user-add"`
	ListUserDelete *InputActionTOML `toml:"list-user-delete"`

	TagOpenFeed *InputActionTOML `toml:"tag-open-feed"`
	TagFollow   *InputActionTOML `toml:"tag-follow"`

	LinkOpen *InputActionTOML `toml:"link-open"`
	LinkYank *InputActionTOML `toml:"link-yank"`

	ComposeEditCW               *InputActionTOML `toml:"compose-edit-cw"`
	ComposeEditText             *InputActionTOML `toml:"compose-edit-text"`
	ComposeIncludeQuote         *InputActionTOML `toml:"compose-include-quote"`
	ComposeMediaFocus           *InputActionTOML `toml:"compose-media-focus"`
	ComposePost                 *InputActionTOML `toml:"compose-post"`
	ComposeToggleContentWarning *InputActionTOML `toml:"compose-toggle-content-warning"`
	ComposeVisibility           *InputActionTOML `toml:"compose-visibility"`
	ComposeLanguage             *InputActionTOML `toml:"compose-language"`
	ComposePoll                 *InputActionTOML `toml:"compose-poll"`

	MediaDelete   *InputActionTOML `toml:"media-delete"`
	MediaEditDesc *InputActionTOML `toml:"media-edit-desc"`
	MediaAdd      *InputActionTOML `toml:"media-add"`

	VoteVote   *InputActionTOML `toml:"vote-vote"`
	VoteSelect *InputActionTOML `toml:"vote-select"`

	PollAdd         *InputActionTOML `toml:"poll-add"`
	PollEdit        *InputActionTOML `toml:"poll-edit"`
	PollDelete      *InputActionTOML `toml:"poll-delete"`
	PollMultiToggle *InputActionTOML `toml:"poll-multi-toggle"`
	PollExpiration  *InputActionTOML `toml:"poll-expiration"`

	PreferenceName         *InputActionTOML `toml:"preference-name"`
	PreferenceVisibility   *InputActionTOML `toml:"preference-visibility"`
	PreferenceBio          *InputActionTOML `toml:"preference-bio"`
	PreferenceSave         *InputActionTOML `toml:"preference-save"`
	PreferenceFields       *InputActionTOML `toml:"preference-fields"`
	PreferenceFieldsAdd    *InputActionTOML `toml:"preference-fields-add"`
	PreferenceFieldsEdit   *InputActionTOML `toml:"preference-fields-edit"`
	PreferenceFieldsDelete *InputActionTOML `toml:"preference-fields-delete"`

	EditorExit *InputActionTOML `toml:"editor-exit"`
}
