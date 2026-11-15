package text

import (
	"fmt"
	"html"
	"strings"
	"time"

	"logware/internal/tg/types"
)

const (
	RegBtn  = "Начать"
	SetBtn  = "Настройки"
	ShowBtn = "Показать"
	MenuBtn = "Logware"
	ProfBtn = "Перейти в настройки"
	ProfURL = "tg://settings"
)

const (
	__chan__ = "https://t.me/logwaree"
	__docs__ = "https://docs.logware.dev"
)

func Links() string {
	return `<a href="` + __chan__ + `">Канал</a>  ·  <a href="` + __docs__ + `">Документация</a>`
}

func quote(__s string) string {
	return "<blockquote>" + __s + "</blockquote>"
}

func Hello() string {
	return E(Logo, "⚓️") + " <b>Logware</b>\n\n" +
		"Всё, что удалят или изменят в твоих чатах, останется с тобой. Под замком, ключ от которого есть <b>только у тебя</b>.\n\n" +
		Links()
}

func Ready(__linked bool, __nick string) string {
	if __linked {
		return E(Check, "✅") + " <b>Всё работает</b>\n\n" +
			"Я подключён к аккаунту и сохраняю сообщения. Настройки и история в приложении.\n\n" + Links()
	}
	return "🔐 <b>Регистрация пройдена</b>\n\n" + HowTo(__nick) + "\n\n" + Links()
}

func Path(__nick string) string {
	return "Профиль &gt; Изменить &gt; ИИ-автоматизация &gt; @" + html.EscapeString(__nick)
}

func HowTo(__nick string) string {
	return "Осталось подключить меня к аккаунту:\n" + quote(Path(__nick)) + "\nИ выдай доступ к личным чатам."
}

func Linked(__r types.Rights, __reg bool, __nick string) string {
	_b := &strings.Builder{}
	_b.WriteString(E(Check, "✅") + " <b>Бот подключён</b>\n\n")
	if __reg {
		_b.WriteString("Сохраняю новые сообщения. Настройки и история в приложении.")
	} else {
		_b.WriteString("Пока ты не пройдёшь регистрацию, я ничего не сохраняю: без твоего ключа шифровать нечем.")
	}
	var _miss []string
	if !__r.Reply {
		_miss = append(_miss, "Отвечать на сообщения")
	}
	if !__r.DelAll {
		_miss = append(_miss, "Удалять сообщения")
	}
	if len(_miss) > 0 {
		_b.WriteString("\n\n" + E(Alert, "⚠️") + " <b>Включи ещё права</b>" + quote(strings.Join(_miss, "\n")) +
			"Там же, где подключал бота:\n" + Path(__nick))
	}
	_b.WriteString("\n\n" + Links())
	return _b.String()
}

func Unlinked() string {
	return "🔌 <b>Бот отключён</b>\n" + quote("Новые сообщения больше не сохраняются. Всё сохранённое осталось в приложении.")
}

func chat(__name string) string {
	return "Чат с <b>" + __name + "</b>"
}

func Del(__name string, __n int) string {
	_word := plural(__n, "сообщение", "сообщения", "сообщений")
	if __n >= 50 {
		return fmt.Sprintf("%s <b>Переписку очистили</b>\n%s", E(Spark, "✨"), quote(fmt.Sprintf("%s · сохранил %d %s", chat(__name), __n, _word)))
	}
	if __n == 1 {
		return fmt.Sprintf("%s <b>Удалено сообщение</b>\n%s", E(Trash, "🗑"), quote(chat(__name)))
	}
	return fmt.Sprintf("%s <b>Удалено %d %s</b>\n%s", E(Trash, "🗑"), __n, _word, quote(chat(__name)))
}

func Edit(__name string) string {
	return fmt.Sprintf("%s <b>Изменено сообщение</b>\n%s", E(Pencil, "✏️"), quote(chat(__name)))
}

func Kw(__name string) string {
	return fmt.Sprintf("%s <b>Ключевое слово</b>\n%s", E(Bell, "🔔"), quote(chat(__name)))
}

func Panic(__name string, __n int) string {
	return fmt.Sprintf("%s <b>Массовое удаление</b>\n%s", E(Alert, "❗️"), quote(fmt.Sprintf("%s · разом удалили %d %s", chat(__name), __n, plural(__n, "сообщение", "сообщения", "сообщений"))))
}

func Saved(__name string) string {
	return fmt.Sprintf("%s <b>Одноразовое сохранено</b>\n%s", E(Camera, "📸"), quote(chat(__name)))
}

func Name(__c types.Chat) string {
	_n := strings.TrimSpace(__c.First + " " + __c.Last)
	if _n == "" {
		_n = __c.Title
	}
	if _n == "" && __c.Nick != "" {
		_n = "@" + __c.Nick
	}
	if _n == "" {
		_n = "собеседником"
	}
	return html.EscapeString(_n)
}

func plural(__n int, __one, __few, __many string) string {
	_m10, _m100 := __n%10, __n%100
	switch {
	case _m10 == 1 && _m100 != 11:
		return __one
	case _m10 >= 2 && _m10 <= 4 && (_m100 < 12 || _m100 > 14):
		return __few
	}
	return __many
}

func Paid(__plan string, __until time.Time, __again bool) string {
	_h := " подключён"
	if __again {
		_h = " продлён"
	}
	return E(Spark, "✨") + " <b>" + html.EscapeString(__plan) + _h + "</b>\n\n" +
		"Действует до " + __until.Format("02.01.2006") + ".\n\n" + Links()
}

const DevBtn = "Кабинет"

func repo(__url string) string {
	return `<a href="` + html.EscapeString(__url) + `">` + html.EscapeString(strings.TrimPrefix(__url, "https://")) + `</a>`
}

func Kind(__k string) string {
	switch __k {
	case "update":
		return "Обновление"
	case "edit":
		return "Изменение"
	}
	return "Новый плагин"
}

func Review(__nick string, __user int64, __kind, __slug, __name, __ver, __about, __note, __repo string, __pic bool) string {
	_b := &strings.Builder{}
	_b.WriteString(E(Bell, "🔔") + " <b>" + Kind(__kind) + " на проверке</b>\n\n")
	_b.WriteString("<b>" + html.EscapeString(__name) + "</b>")
	if __ver != "" {
		_b.WriteString(" · " + html.EscapeString(__ver))
	}
	fmt.Fprintf(_b, " · <code>%s</code>", html.EscapeString(__slug))
	if __about != "" {
		_b.WriteString("\n" + quote(html.EscapeString(__about)))
	} else {
		_b.WriteString("\n")
	}
	fmt.Fprintf(_b, "\nАвтор: <b>%s</b> · <code>%d</code>", html.EscapeString(__nick), __user)
	if __kind != "edit" {
		_b.WriteString("\nКод: " + repo(__repo))
	}
	if __pic {
		_b.WriteString("\nНовая аватарка, смотри в админке")
	}
	if __note != "" {
		_b.WriteString("\n\nКомментарий автора:" + quote(html.EscapeString(__note)))
	}
	_b.WriteString("\n\n<i>Ответь на это сообщение, чтобы отклонить с причиной.</i>")
	return _b.String()
}

func what(__kind, __name, __ver string) string {
	_n := "<b>" + html.EscapeString(__name) + "</b>"
	switch __kind {
	case "update":
		return "Обновление " + _n + " до " + html.EscapeString(__ver)
	case "edit":
		return "Изменения " + _n
	}
	return "Плагин " + _n
}

func Reviewed(__kind, __name, __ver, __repo string, __ok bool, __why string) string {
	if __ok {
		return E(Check, "✅") + " " + what(__kind, __name, __ver) + ": одобрено\n\nКод: " + repo(__repo)
	}
	_s := E(Alert, "⚠️") + " " + what(__kind, __name, __ver) + ": отклонено"
	if __why != "" {
		_s += "\n" + quote(html.EscapeString(__why))
	} else {
		_s += "\n"
	}
	return _s + "\nКод: " + repo(__repo)
}

func WhySent() string {
	return E(Check, "✅") + " Причина отправлена разработчику."
}

func Verdict(__kind, __name, __ver string, __ok bool, __why string) string {
	if __ok {
		_t := "Плагин появится в магазине со следующим обновлением бота."
		if __kind == "edit" {
			_t = "Новое название, описание и аватарка уже в магазине."
		}
		return E(Check, "✅") + " <b>" + what(__kind, __name, __ver) + " прошло проверку</b>\n\n" + _t + " Всё видно во вкладке Developer."
	}
	_s := E(Alert, "⚠️") + " <b>" + what(__kind, __name, __ver) + " не прошло проверку</b>\n"
	if __why != "" {
		_s += quote(html.EscapeString(__why))
	} else {
		_s += "\n"
	}
	return _s + "Исправь и отправь снова во вкладке Developer."
}

func Banned() string {
	return E(Alert, "⚠️") + " <b>Доступ к боту закрыт</b>\n\nЕсли это ошибка, напиши в наш канал." + "\n\n" + Links()
}

func PaidCoin(__plan string, __until time.Time) string {
	return E(Spark, "✨") + " <b>" + html.EscapeString(__plan) + " подключён</b>\n\n" +
		"Действует до " + __until.Format("02.01.2006") + ".\n\n" +
		"Так как ты оплатил криптой, автопродление не работает. За 3 дня до конца подписки придёт напоминание об оплате.\n\n" + Links()
}

func Renew(__plan string, __until time.Time) string {
	return E(Bell, "🔔") + " <b>" + html.EscapeString(__plan) + " заканчивается " + __until.Format("02.01") + "</b>\n\n" +
		"Продли подписку, чтобы плагины и время в имени продолжили работать. Если оплатить заранее, новый месяц добавится к текущему сроку."
}

const RenewBtn = "Продлить"

const (
	TermsBtn   = "Открыть соглашение"
	PrivacyBtn = "Открыть политику"
	SupportBtn = "Написать в поддержку"
)

func Terms() string {
	return "<b>Пользовательское соглашение</b>\n\nУсловия использования Logware, подписки и возврат денег."
}

func Privacy() string {
	return "<b>Политика конфиденциальности</b>\n\nКакие данные мы храним, как они зашифрованы и как их удалить."
}

func Support(__nick string) string {
	return E(Bell, "💬") + " <b>Поддержка</b>\n\nВопросы по работе бота и по оплате принимаем в @" + html.EscapeString(__nick) + ". Если вопрос про платёж, укажи дату и способ оплаты."
}

func NoSupport() string {
	return E(Bell, "💬") + " <b>Поддержка</b>\n\nПоддержка скоро заработает. Пока напиши в наш канал." + "\n\n" + Links()
}
