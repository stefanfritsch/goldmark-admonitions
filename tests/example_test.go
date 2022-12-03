package admonitions_test

import (
	"os"

	admonitions "github.com/stefanfritsch/goldmark-admonitions"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func Example() {
	src := []byte(`
## Hello

The following is an admonition:

!!!!note This is a note
The body
!!!!

This is the end.
  `)

	markdown := goldmark.New(
		goldmark.WithExtensions(
			&admonitions.Extender{},
		),
	)

	doc := markdown.Parser().Parse(text.NewReader(src))
	markdown.Renderer().Render(os.Stdout, src, doc)

	// Output:
	// 	<h2>Hello</h2>
	// <p>The following is an admonition:</p>
	// <div class="admonition adm-note">
	// <div class="adm-title">This is a note</div>
	//   <div class="adm-body">
	// <p>The body</p>
	//   </div>
	// </div>
	// <p>This is the end.</p>
}

func Example_noTitle() {
	src := []byte(`
## Hello

This is no admonition (no class)

!!!
Not an admonition
!!!

The following is an admonition:

!!!!note 
The body
!!!!

!!! note With title!
The body
!!!

!!!!danger With Attributes!{.otherclass}
The body
!!!!

This is the end.
  `)

	markdown := goldmark.New(
		goldmark.WithExtensions(
			&admonitions.Extender{},
		),
	)

	doc := markdown.Parser().Parse(text.NewReader(src))
	markdown.Renderer().Render(os.Stdout, src, doc)

	// Output:
	//   <h2>Hello</h2>
	// <p>This is no admonition (no class)</p>
	// <p>!!!
	// Not an admonition
	// !!!</p>
	// <p>The following is an admonition:</p>
	// <div class="admonition adm-note">
	// <div class="adm-title"></div>
	//   <div class="adm-body">
	// <p>The body</p>
	//   </div>
	// </div>
	// <div class="admonition adm-note">
	// <div class="adm-title">With title!</div>
	//   <div class="adm-body">
	// <p>The body</p>
	//   </div>
	// </div>
	// <div class="admonition adm-danger otherclass">
	// <div class="adm-title">With Attributes!</div>
	//   <div class="adm-body">
	// <p>The body</p>
	//   </div>
	// </div>
	// <p>This is the end.</p>
}

func Example_nested() {
	src := []byte(`
## Hello

The following contains an id and a class

!!!note This is a note {#big-div .add-border}

And the next admonition contains two classes.

!!!danger And this is danger {.background-green .font-big}
## This is nested within nested admonitions

here we close the inner admonition:
!!!

and finally the outer one:
!!!`)

	markdown := goldmark.New(
		goldmark.WithExtensions(
			&admonitions.Extender{},
		),
	)

	doc := markdown.Parser().Parse(text.NewReader(src))
	markdown.Renderer().Render(os.Stdout, src, doc)

	// Output:
	// <h2>Hello</h2>
	// <p>The following contains an id and a class</p>
	// <div id="big-div" class="admonition adm-note add-border">
	// <div class="adm-title">This is a note</div>
	//   <div class="adm-body">
	// <p>And the next admonition contains two classes.</p>
	// <div class="admonition adm-danger background-green font-big">
	// <div class="adm-title">And this is danger</div>
	//   <div class="adm-body">
	// <h2>This is nested within nested admonitions</h2>
	// <p>here we close the inner admonition:</p>
	//   </div>
	// </div>
	// <p>and finally the outer one:</p>
	//   </div>
	// </div>
}
