package ui

import (
	"github.com/Arceliar/phony"
	"github.com/rivo/tview"
)

type ModelViewOnceAction struct {
	fnYes func()
	fnNo  func()
}

type ModalView struct {
	phony.Inbox
	onceAction *ModelViewOnceAction
	tutView    *TutView
	View       *tview.Modal
}

func NewModalView(tv *TutView) *ModalView {
	mv := &ModalView{
		tutView: tv,
		View:    NewModal(tv.tut.Config),
	}
	mv.View.SetText("Are you sure?").
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				mv.Be_Done(mv, true)
			} else {
				mv.Be_Done(mv, false)
			}
		})
	return mv
}

func (mv *ModalView) Be_Done(from phony.Actor, res bool) {
	mv.Act(from, func() {
		action := mv.onceAction
		mv.onceAction = nil
		if res {
			action.fnYes()
		} else {
			action.fnNo()
		}
		mv.tutView.tut.App.QueueUpdateDraw(func() {
			mv.tutView.PrevFocus()
		})
	})
}

func (mv *ModalView) Be_OpenConfirm(from phony.Actor, text string, fn func()) {
	if !mv.tutView.tut.Config.General.Confirmation {
		fn()
		return
	}

	mv.Be_Open(from, text, fn, func() {})
}

func (mv *ModalView) Be_Open(from phony.Actor, text string, fnYes func(), fnNo func()) {
	mv.Act(from, func() {
		mv.onceAction = &ModelViewOnceAction{fnYes, fnNo}
		mv.View.SetFocus(0)
		mv.View.SetText(text)
		mv.tutView.tut.App.QueueUpdateDraw(func() {
			mv.tutView.SetPage(ModalFocus)
		})
	})
}

// no idea what this is for
func (mv *ModalView) Stop(fn func()) {
	fn()
}
