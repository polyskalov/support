package internal

type Dialog struct {
	User User `json:"user"`
}

type Dialogs struct {
	dialogs []Dialog
}
