package views

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func FrontPage() g.Node {
	return Page(
		"Canvas",
		"/",
		H1(g.Text(`Solution to problems.`)),
		P(g.Text(`Do you have problems? We also had problems.`)),
		P(g.Raw(`Then we created the <em>canvas</em>, and now we have triple times more.`)),
	)
}
