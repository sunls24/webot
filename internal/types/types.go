package types

type LoginMode string

const (
	LoginNormal  LoginMode = "normal"
	LoginDesktop LoginMode = "desktop"
)

type Func string

const (
	FuncHelp  Func = "help"
	FuncClear Func = "clear"
	FuncV2ex  Func = "v2ex"
)

type SettingsKey string

const (
	PushV2ex SettingsKey = "push_v2ex"
)
