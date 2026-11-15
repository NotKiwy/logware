package text

const (
	Logo   = "5458491451363073497"
	Trash  = "5458618350466802842"
	Pencil = "5458662369586623553"
	Camera = "5456429158391393224"
	Bell   = "5456143654735358732"
	Alert  = "5456596997123389820"
	Spark  = "5456128278752439217"
	Check  = "5458695866036565296"
	Eye    = "5456257862210726370"
	Gear   = "5458443777226093464"
)

func E(__id, __alt string) string {
	return `<tg-emoji emoji-id="` + __id + `">` + __alt + `</tg-emoji>`
}
