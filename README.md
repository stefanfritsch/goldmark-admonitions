# Goldmark-Admonitions

## Overview

[Goldmark](github.com/yuin/goldmark) is a fast markdown renderer for Go. Admonitions are a markdown extension that allows you to style markdown as nice boxes with a title.

### An Example

```markdown
!!!note This is a note
This is the body

## You can use other markdown elements in admonitions!

We are now inside a div with the css-class "adm-body" wrapped inside a div with "admonition" and "adm-note". This can be used to style this block
!!!

!!!danger This a warning
The same as above but instead of "adm-note" you have "adm-danger"
!!!
```

Now add the css to your stylesheet:

```css
.admonition {
  border: 1px solid black;
  border-radius: 0.25rem;
}

.adm-title {color: #efefef;}
.adm-body {background-color: lightgrey;}

.adm-note .adm-title {background-color: darkblue;}

.adm-danger .adm-title {background-color: darkred;}

```

If you do that the admonitions will look like this (this is an image as GitHub doesn't allow custom css in READMEs):

![](assets/Screenshot%202022-10-14%20001453.png)

### Full Example

A full code example that renders to stdout could look like this:

```go
func main() {
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
}
```

## Indented Style

Instead of opening and closing blocks, i.e.

```markdown
!!!note This is a note
   
   This is the body
!!!

and this isn't
```

you can also use indentation for compatibility with markdown-it:

```markdown
!!!note This is a note
   
   this is the body

and this isn't
```
