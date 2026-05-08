package sharedtypes

type UpdateProfileEnter struct{}
type UpdateProfileExit struct{}
type AddProfileExit struct{}
type AddProfileEnter struct{}
type HelpViewEnter struct{}
type HelpViewExit struct{}
type ClearWarnings struct{}
type AddSelectorEnter struct{}
type AddSelectorExit struct{}
type PasteProfileEnter struct{}
type PasteProfileExit struct{}
type PasteProfileSubmitted struct{}
type ThemeViewEnter struct{}
type ThemeViewExit struct{}

type ClearTestResult struct {
	GroupId   int
	ProfileId int
}
