package config

var tvar = true
var fvar = false

var bt = &tvar
var bf = &fvar

func sp(s string) *string {
	return &s
}
func ip(i int) *int {
	return &i
}

func ip64(i int64) *int64 {
	return &i
}

var ConfigDefault = ConfigTOML{
	General: GeneralTOML{
		Editor:              sp("TUT_USE_INTERNAL"),
		Confirmation:        bt,
		MouseSupport:        bf,
		DateFormat:          sp("2006-01-02 15:04"),
		DateTodayFormat:     sp("15:04"),
		DateRelative:        ip(-1),
		QuoteReply:          bf,
		MaxWidth:            ip(0),
		ShortHints:          bf,
		ShowFilterPhrase:    bt,
		ShowIcons:           bt,
		ShowHelp:            bt,
		RedrawUI:            bt,
		StickToTop:          bf,
		ShowBoostedUser:     bf,
		DynamicTimelineName: bt,
		CommandsInNewPane:   bt,
		ListPlacement:       sp("left"),
		ListSplit:           sp("row"),
		ListProportion:      ip(1),
		ContentProportion:   ip(2),
		TerminalTitle:       ip(0),
		LeaderKey:           sp(""),
		LeaderTimeout:       ip64(1000),
		NotificationsToHide: &[]string{},
		Timelines: &[]TimelineTOML{
			{
				Name:        sp("Home"),
				Type:        sp("home"),
				HideBoosts:  bf,
				HideReplies: bf,
			},
			{
				Name: sp("Notifications"),
				Type: sp("notifications"),
				Keys: &[]string{"n", "N"},
			},
		},
	},
	Style: StyleTOML{
		Theme:                          sp("none"),
		XrdbPrefix:                     sp("guess"),
		Background:                     sp("#272822"),
		Text:                           sp("#f8f8f2"),
		Subtle:                         sp("#808080"),
		WarningText:                    sp("#f92672"),
		TextSpecial1:                   sp("#ae81ff"),
		TextSpecial2:                   sp("#a6e22e"),
		TopBarBackground:               sp("#f92672"),
		TopBarText:                     sp("#f8f8f2"),
		StatusBarBackground:            sp("#f92672"),
		StatusBarText:                  sp("#f8f8f2"),
		StatusBarViewBackground:        sp("#ae81ff"),
		StatusBarViewText:              sp("#f8f8f2"),
		CommandText:                    sp("#f8f8f2"),
		ListSelectedBackground:         sp("#f92672"),
		ListSelectedText:               sp("#f8f8f2"),
		ListSelectedInactiveBackground: sp("#ae81ff"),
		ListSelectedInactiveText:       sp("#f8f8f2"),
		ControlsText:                   sp("#f8f8f2"),
		ControlsHighlight:              sp("#a6e22e"),
		AutocompleteBackground:         sp("#272822"),
		AutocompleteText:               sp("#f8f8f2"),
		AutocompleteSelectedBackground: sp("#ae81ff"),
		AutocompleteSelectedText:       sp("#f8f8f2"),
		ButtonColorOne:                 sp("#f92672"),
		ButtonColorTwo:                 sp("#272822"),
		TimelineNameBackground:         sp("#272822"),
		TimelineNameText:               sp("#808080"),
	},
	Media: MediaTOML{
		DeleteTmpFiles: bt,
		Image: &ViewerTOML{
			Program:  sp("TUT_OS_DEFAULT"),
			Args:     sp(""),
			Terminal: bf,
			Single:   bt,
			Reverse:  bf,
		},
		Video: &ViewerTOML{
			Program:  sp("TUT_OS_DEFAULT"),
			Args:     sp(""),
			Terminal: bf,
			Single:   bt,
			Reverse:  bf,
		},
		Audio: &ViewerTOML{
			Program:  sp("TUT_OS_DEFAULT"),
			Args:     sp(""),
			Terminal: bf,
			Single:   bt,
			Reverse:  bf,
		},
		Link: &ViewerTOML{
			Program:  sp("TUT_OS_DEFAULT"),
			Args:     sp(""),
			Terminal: bf,
			Single:   bt,
			Reverse:  bf,
		},
	},
	NotificationConfig: NotificationsTOML{
		Followers: bf,
		Favorite:  bf,
		Mention:   bf,
		Update:    bf,
		Boost:     bf,
		Poll:      bf,
		Posts:     bf,
	},
	Input: InputTOML{
		GlobalDown: &InputActionTOML{
			Keys:        &[]string{"j", "J"},
			SpecialKeys: &[]string{"Down"},
		},
		GlobalUp: &InputActionTOML{
			Keys:        &[]string{"k", "K"},
			SpecialKeys: &[]string{"Up"},
		},
		GlobalEnter: &InputActionTOML{
			SpecialKeys: &[]string{"Enter"},
		},
		GlobalBack: &InputActionTOML{
			Hint:        sp("[Esc]"),
			SpecialKeys: &[]string{"Esc"},
		},
		GlobalExit: &InputActionTOML{
			Hint: sp("[Q]uit"),
			Keys: &[]string{"q", "Q"},
		},
		MainHome: &InputActionTOML{
			Hint:        sp(""),
			Keys:        &[]string{"g"},
			SpecialKeys: &[]string{"Home"},
		},
		MainEnd: &InputActionTOML{
			Hint:        sp(""),
			Keys:        &[]string{"G"},
			SpecialKeys: &[]string{"End"},
		},
		MainPrevFeed: &InputActionTOML{
			Hint:        sp(""),
			Keys:        &[]string{"h", "H"},
			SpecialKeys: &[]string{"Left"},
		},
		MainNextFeed: &InputActionTOML{
			Hint:        sp(""),
			Keys:        &[]string{"l", "L"},
			SpecialKeys: &[]string{"Right"},
		},
		MainPrevPane: &InputActionTOML{
			Hint:        sp(""),
			SpecialKeys: &[]string{"Backtab"},
		},
		MainNextPane: &InputActionTOML{
			Hint:        sp(""),
			SpecialKeys: &[]string{"Tab"},
		},
		MainCompose: &InputActionTOML{
			Hint: sp(""),
			Keys: &[]string{"c", "C"},
		},
		MainNextAccount: &InputActionTOML{
			Hint:        sp(""),
			SpecialKeys: &[]string{"Ctrl-N"},
		},
		MainPrevAccount: &InputActionTOML{
			Hint:        sp(""),
			SpecialKeys: &[]string{"Ctrl-P"},
		},
		StatusAvatar: &InputActionTOML{
			Hint: sp("[A]vatar"),
			Keys: &[]string{"a", "A"},
		},
		StatusBoost: &InputActionTOML{
			Hint:    sp("[B]oost"),
			HintAlt: sp("Un[B]oost"),
			Keys:    &[]string{"b", "B"},
		},
		StatusEdit: &InputActionTOML{
			Hint: sp("[E]dit"),
			Keys: &[]string{"e", "E"},
		},
		StatusDelete: &InputActionTOML{
			Hint: sp("[D]elete"),
			Keys: &[]string{"d", "D"},
		},
		StatusFavorite: &InputActionTOML{
			Hint:    sp("[F]avorite"),
			HintAlt: sp("Un[F]avorite"),
			Keys:    &[]string{"f", "F"},
		},
		StatusMedia: &InputActionTOML{
			Hint: sp("[M]edia"),
			Keys: &[]string{"m", "M"},
		},
		StatusLinks: &InputActionTOML{
			Hint: sp("[O]pen"),
			Keys: &[]string{"o", "O"},
		},
		StatusPoll: &InputActionTOML{
			Hint: sp("[P]oll"),
			Keys: &[]string{"p", "P"},
		},
		StatusReply: &InputActionTOML{
			Hint: sp("[R]eply"),
			Keys: &[]string{"r", "R"},
		},
		StatusBookmark: &InputActionTOML{
			Hint:    sp("[S]ave"),
			HintAlt: sp("Un[S]ave"),
			Keys:    &[]string{"s", "S"},
		},
		StatusThread: &InputActionTOML{
			Hint: sp("[T]hread"),
			Keys: &[]string{"t", "T"},
		},
		StatusUser: &InputActionTOML{
			Hint: sp("[U]ser"),
			Keys: &[]string{"u", "U"},
		},
		StatusViewFocus: &InputActionTOML{
			Hint: sp("[V]iew"),
			Keys: &[]string{"v", "V"},
		},
		StatusYank: &InputActionTOML{
			Hint: sp("[Y]ank"),
			Keys: &[]string{"y", "Y"},
		},
		StatusToggleCW: &InputActionTOML{
			Hint: sp("Press [Z] to toggle cw"),
			Keys: &[]string{"z", "Z"},
		},
		StatusShowFiltered: &InputActionTOML{
			Hint: sp("Press [Z] to view filtered toot"),
			Keys: &[]string{"z", "Z"},
		},
		UserAvatar: &InputActionTOML{
			Hint: sp("[A]vatar"),
			Keys: &[]string{"a", "A"},
		},
		UserBlock: &InputActionTOML{
			Hint:    sp("[B]lock"),
			HintAlt: sp("Un[B]lock"),
			Keys:    &[]string{"b", "B"},
		},
		UserFollow: &InputActionTOML{
			Hint:    sp("[F]ollow"),
			HintAlt: sp("Un[F]ollow"),
			Keys:    &[]string{"f", "F"},
		},
		UserFollowRequestDecide: &InputActionTOML{
			Hint:    sp("Follow [R]equest"),
			HintAlt: sp("Follow [R]equest"),
			Keys:    &[]string{"r", "R"},
		},
		UserMute: &InputActionTOML{
			Hint:    sp("[M]ute"),
			HintAlt: sp("Un[M]ute"),
			Keys:    &[]string{"m", "M"},
		},
		UserLinks: &InputActionTOML{
			Hint: sp("[O]pen"),
			Keys: &[]string{"o", "O"},
		},
		UserUser: &InputActionTOML{
			Hint: sp("[U]ser"),
			Keys: &[]string{"u", "U"},
		},
		UserViewFocus: &InputActionTOML{
			Hint: sp("[V]iew"),
			Keys: &[]string{"v", "V"},
		},
		UserYank: &InputActionTOML{
			Hint: sp("[Y]ank"),
			Keys: &[]string{"y", "Y"},
		},
		ListOpenFeed: &InputActionTOML{
			Hint: sp("[O]pen"),
			Keys: &[]string{"o", "O"},
		},
		ListUserList: &InputActionTOML{
			Hint: sp("[U]sers"),
			Keys: &[]string{"u", "U"},
		},
		ListUserAdd: &InputActionTOML{
			Hint: sp("[A]dd"),
			Keys: &[]string{"a", "A"},
		},
		ListUserDelete: &InputActionTOML{
			Hint: sp("[D]elete"),
			Keys: &[]string{"d", "D"},
		},
		LinkOpen: &InputActionTOML{
			Hint: sp("[O]pen"),
			Keys: &[]string{"o", "O"},
		},
		LinkYank: &InputActionTOML{
			Hint: sp("[Y]ank"),
			Keys: &[]string{"y", "Y"},
		},
		TagOpenFeed: &InputActionTOML{
			Hint: sp("[O]pen"),
			Keys: &[]string{"o", "O"},
		},
		TagFollow: &InputActionTOML{
			Hint:    sp("[F]ollow"),
			HintAlt: sp("Un[F]ollow"),
			Keys:    &[]string{"f", "F"},
		},
		ComposeEditCW: &InputActionTOML{
			Hint: sp("[C]W text"),
			Keys: &[]string{"c", "C"},
		},
		ComposeEditText: &InputActionTOML{
			Hint: sp("[E]dit text"),
			Keys: &[]string{"e", "E"},
		},
		ComposeIncludeQuote: &InputActionTOML{
			Hint: sp("[I]nclude quote"),
			Keys: &[]string{"i", "I"},
		},
		ComposeMediaFocus: &InputActionTOML{
			Hint: sp("[M]edia"),
			Keys: &[]string{"m", "M"},
		},
		ComposePost: &InputActionTOML{
			Hint: sp("[P]ost"),
			Keys: &[]string{"p", "P"},
		},
		ComposeToggleContentWarning: &InputActionTOML{
			Hint: sp("[T]oggle CW"),
			Keys: &[]string{"t", "T"},
		},
		ComposeVisibility: &InputActionTOML{
			Hint: sp("[V]isibility"),
			Keys: &[]string{"v", "V"},
		},
		ComposeLanguage: &InputActionTOML{
			Hint: sp("[L]ang"),
			Keys: &[]string{"l", "L"},
		},
		ComposePoll: &InputActionTOML{
			Hint: sp("P[O]ll"),
			Keys: &[]string{"o", "O"},
		},
		MediaDelete: &InputActionTOML{
			Hint: sp("[D]elete"),
			Keys: &[]string{"d", "D"},
		},
		MediaEditDesc: &InputActionTOML{
			Hint: sp("[E]dit desc"),
			Keys: &[]string{"e", "E"},
		},
		MediaAdd: &InputActionTOML{
			Hint: sp("[A]dd"),
			Keys: &[]string{"a", "A"},
		},
		VoteVote: &InputActionTOML{
			Hint: sp("[V]ote"),
			Keys: &[]string{"v", "V"},
		},
		VoteSelect: &InputActionTOML{
			Hint:        sp("[Enter] to select"),
			Keys:        &[]string{" "},
			SpecialKeys: &[]string{"Enter"},
		},
		PollAdd: &InputActionTOML{
			Hint: sp("[A]dd"),
			Keys: &[]string{"a", "A"},
		},
		PollEdit: &InputActionTOML{
			Hint: sp("[E]dit"),
			Keys: &[]string{"e", "E"},
		},
		PollDelete: &InputActionTOML{
			Hint: sp("[D]elete"),
			Keys: &[]string{"d", "D"},
		},
		PollMultiToggle: &InputActionTOML{
			Hint: sp("Toggle [M]ultiple"),
			Keys: &[]string{"m", "M"},
		},
		PollExpiration: &InputActionTOML{
			Hint: sp("E[X]pires"),
			Keys: &[]string{"x", "X"},
		},
		PreferenceName: &InputActionTOML{
			Hint: sp("[N]ame"),
			Keys: &[]string{"n", "N"},
		},
		PreferenceVisibility: &InputActionTOML{
			Hint: sp("[V]isibility"),
			Keys: &[]string{"v", "V"},
		},
		PreferenceBio: &InputActionTOML{
			Hint: sp("[B]io"),
			Keys: &[]string{"b", "B"},
		},
		PreferenceSave: &InputActionTOML{
			Hint: sp("[S]ave"),
			Keys: &[]string{"s", "S"},
		},
		PreferenceFields: &InputActionTOML{
			Hint: sp("[F]ields"),
			Keys: &[]string{"f", "F"},
		},
		PreferenceFieldsAdd: &InputActionTOML{
			Hint: sp("[A]dd"),
			Keys: &[]string{"a", "A"},
		},
		PreferenceFieldsEdit: &InputActionTOML{
			Hint: sp("[E]dit"),
			Keys: &[]string{"e", "E"},
		},
		PreferenceFieldsDelete: &InputActionTOML{
			Hint: sp("[D]elete"),
			Keys: &[]string{"d", "D"},
		},
		EditorExit: &InputActionTOML{
			Hint:        sp("[Esc] when done"),
			SpecialKeys: &[]string{"Esc"},
		},
	},
}
