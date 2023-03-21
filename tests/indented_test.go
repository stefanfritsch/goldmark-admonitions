package admonitions_test

import (
	"math/rand"
	"os"

	admonitions "github.com/stefanfritsch/goldmark-admonitions"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func Example_indented() {
	src := []byte(`
## Hello

The following contains an id and a class

!!!note This is a note {#big-div .add-border}

   And the content is indented.
   
   ## This is still in the note
	 
And this isn't.
`)

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
	// <div id="big-div" class="admonition adm-note add-border" data-admonition="0">
	//   <div class="adm-title">This is a note</div>
	//   <div class="adm-body">
	// <p>And the content is indented.</p>
	// <h2>This is still in the note</h2>
	//   </div>
	// </div>
	// <p>And this isn't.</p>
}

func Example_nested_indented() {
	// one of the admonitions isn't closed correctly and so keeps the random
	// admonition id. The seed keeps this repeatable, independently of the number
	// of tests.
	rand.Seed(1)
	
	src := []byte(`
## Hello

The following contains an id and a class

!!!note This is a note {#big-div .add-border}

   And the content is indented.
   
   !!!danger This is still in the note
      
      this is a deeper error
  
   this is level 1 again
   
   !!!note This is another note
      in level 2
   
And this isn't.
`)

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
  // <div id="big-div" class="admonition adm-note add-border" data-admonition="0">
  //   <div class="adm-title">This is a note</div>
  //   <div class="adm-body">
  // <p>And the content is indented.</p>
  // <div class="admonition adm-danger" data-admonition="1">
  //   <div class="adm-title">This is still in the note</div>
  //   <div class="adm-body">
  // <p>this is a deeper error</p>
  //   </div>
  // </div>
  // <p>this is level 1 again</p>
  // <div class="admonition adm-note" data-admonition="xPLDnJObCsNVlgTeMaPEZQle">
  //   <div class="adm-title">This is another note</div>
  //   <div class="adm-body">
  // <p>in level 2</p>
  //   </div>
  // </div>
  //   </div>
  // </div>
  // <p>And this isn't.</p>
}