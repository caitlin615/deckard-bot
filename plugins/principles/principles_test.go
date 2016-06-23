package principles

import (
	"fmt"

	"github.com/handwritingio/deckard-bot/message"
)

func ExamplePlugin_HandleMessage() {
	p := new(Plugin)
	principleList, _ := buildPrinciples(principleData)
	p.List = principleList

	fmt.Println(rePrinciple.MatchString("!principles"))
	fmt.Println(rePrinciple.MatchString("!principle"))
	fmt.Println(rePrincipleNum.MatchString("!principle 99"))
	fmt.Println(rePrincipleKeyword.FindStringSubmatch("!principle clever code")[1])
	fmt.Println(p.HandleMessage(format("!principle 1")).Text)
	fmt.Println(p.HandleMessage(format("!principles 12")).Text)
	fmt.Println(p.HandleMessage(format("!principle clever code")).Text)
	fmt.Println(p.HandleMessage(format("!principles mary had a little lamb")).Text)
	// Output:
	// true
	// true
	// true
	// clever code
	// *Build what matters*: Engineering effort is a scarce commodity. It should only be applied to problems that "move the needle" for the company.
	// Sorry, the principle you requested does not exist
	// *Clear Code Beats Clever Code*: Don't write code you can't debug at 3AM while drunk. Never name a variable 'data' or 'info'. If the implementation is hard to explain, it's a bad idea.
	// Sorry, no principles match keyword `mary had a little lamb`
}

func format(text string) message.Basic {
	return message.Basic{
		ID:       1,
		Text:     text,
		Finished: true,
	}
}

var principleData = []byte(`Engineering Principles
=======================

1. **Build what matters** Engineering effort is a scarce commodity.
It should only be applied to problems that "move the needle" for
the company.

1. **Be a scientist** Science is a systematic enterprise that builds
and organizes knowledge in the form of testable explanations and
predictions about the universe. We should base our decisions on
objective data obtained through research rather than hunches or
superstition. To improve is to change.

1. **Know your enemies** Conway's Law, sleep deprivation, subjective
beliefs, lax validation, crappy dependencies, unreliable networks,
etc... Know them and have a plan to beat them.

1. **Don’t live with broken windows** Fix bad designs, wrong
decisions, and poor code when you see them.

1. **Don’t repeat yourself** Every piece of knowledge/process must
have a single, unambiguous, authoritative representation within a
system.

1. **Clear Code Beats Clever Code** Don't write code you can't debug
at 3AM while drunk. Never name a variable 'data' or 'info'. If the
implementation is hard to explain, it's a bad idea.

1. **It’s Both What You Say and the Way You Say It** There’s no
point in having great ideas if you don’t communicate them effectively.
Be conscious of both your words and your tone.

1. **Sign Your Work** Craftsmen of an earlier age were proud to
sign their work. You should be, too. This can manifest itself as
well documented modules, open-sourcing our tools, or just giving
brown-bags on things you've made.

1. **Protect your Culture** A diversity of opinions is a good thing.
An objective respect for the scientific approach to finding the
best answers to organizational problems leads to a stronger bond
between all of us. Only hire people we would want in the office
*tomorrow*.

1. **Collaborate, Don't Corroborate** Asking for input on 80% done
work isn't fair to the reviewer. Don't get peer support just to
make yourself feel better. Encourage peer support when you're 20%
done, to allow for more meaningful collaboration.

1. **Don’t Gather Requirements, Dig for Them** Requirements rarely
lie on the surface. They’re buried deep beneath layers of assumptions,
misconceptions, and politics. To think like a user, work with a
user.`)
