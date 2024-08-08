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
)
