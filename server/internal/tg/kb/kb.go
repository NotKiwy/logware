package kb

type App struct {
	URL string `json:"url"`
}

type Btn struct {
	Text string `json:"text"`
	Icon string `json:"icon_custom_emoji_id,omitempty"`
	App  *App   `json:"web_app,omitempty"`
	URL  string `json:"url,omitempty"`
	Data string `json:"callback_data,omitempty"`
}

type Kb struct {
	Rows [][]Btn `json:"inline_keyboard"`
}

func Web(__text, __url string) Btn {
	return Btn{Text: __text, App: &App{URL: __url}}
}

func (b Btn) With(__icon string) Btn {
	b.Icon = __icon
	return b
}

func Url(__text, __url string) Btn {
	return Btn{Text: __text, URL: __url}
}

func Rows(__r ...[]Btn) *Kb {
	return &Kb{Rows: __r}
}

func One(__b ...Btn) *Kb {
	return &Kb{Rows: [][]Btn{__b}}
}

func Data(__text, __data string) Btn {
	return Btn{Text: __text, Data: __data}
}
