package views

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func NewsletterThanksPage(path string) g.Node {
	return Page(
		"Thanks for signing up!",
		path,
		H1(g.Text(`Thanks for signing up!`)),
		P(g.Raw(`Please check your email for a confirmation link!`)),
	)
}
