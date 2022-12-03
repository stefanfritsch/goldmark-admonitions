package admonitions

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type admonitionParser struct {
}

var defaultAdmonitionParser = &admonitionParser{}

// NewAdmonitionParser returns a new BlockParser that
// parses fenced admonition blocks.
func NewAdmonitionParser() parser.BlockParser {
	return defaultAdmonitionParser
}

type AdmonitionData struct {
	char   byte
	indent int
	length int
	node   ast.Node
}

var admonitionInfoKey = parser.NewContextKey()

func (b *admonitionParser) Trigger() []byte {
	return []byte{'!'}
}

func (b *admonitionParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, segment := reader.PeekLine()
	pos := pc.BlockOffset()
	if pos < 0 || line[pos] != '!' {
		return nil, parser.NoChildren
	}
	findent := pos

	// currently useless
	admonitionChar := line[pos]
	i := pos
	for ; i < len(line) && line[i] == admonitionChar; i++ {
	}
	oFenceLength := i - pos
	if oFenceLength < 3 {
		return nil, parser.NoChildren
	}

	// ========================================================================== //
	// 	Without attributes we return

	if i >= len(line)-1 {
		// If there are no attributes we can't create a div because we won't know
		// if a "!!!" ends the last admonition or opens a new one
		return nil, parser.NoChildren
	}

	rest := line[i:]
	left := i + util.TrimLeftSpaceLength(rest)
	right := len(line) - 1 - util.TrimRightSpaceLength(rest)

	if left >= right {
		// As above:
		// If there are no attributes we can't create a div because we won't know
		// if a "!!!" ends the last admonition or opens a new one
		return nil, parser.NoChildren
	}

	// ========================================================================== //
	// 	With attributes we construct the node
	node := parseOpeningLine(reader, left)

	fdata := &AdmonitionData{admonitionChar, findent, oFenceLength, node}
	var fdataMap []*AdmonitionData

	if oldData := pc.Get(admonitionInfoKey); oldData != nil {
		fdataMap = oldData.([]*AdmonitionData)
		fdataMap = append(fdataMap, fdata)
	} else {
		fdataMap = []*AdmonitionData{fdata}
	}
	pc.Set(admonitionInfoKey, fdataMap)

	// ========================================================================== //
	// 	 check if it's an empty block

	line, _ = reader.PeekLine()
	w, pos := util.IndentWidth(line, reader.LineOffset())

	if close, _ := b.closes(line, segment, w, pos, node, fdata); close {
		return node, parser.NoChildren
	}

	return node, parser.HasChildren
}

// Parse the opening line for
// * admonition class
// * admonition title
// * attributes
func parseOpeningLine(reader text.Reader, left int) *Admonition {
	node := NewAdmonition()
	reader.Advance(left)

	remainingLine, _ := reader.PeekLine()
	remainingLength := len(remainingLine)

	// ========================================================================== //
	// 	find class
	endClass := 0
	for ; endClass < remainingLength && remainingLine[endClass] != ' ' && remainingLine[endClass] != '{'; endClass++ {
	}
	if endClass > 0 {
		node.AdmonitionClass = remainingLine[0:endClass]
	}

	// ========================================================================== //
	// 	find title
	startTitle := endClass + util.TrimLeftSpaceLength(remainingLine[endClass:])
	endTitle := startTitle
	for ; endTitle < remainingLength && remainingLine[endTitle] != '{'; endTitle++ {
	}
	if endTitle > startTitle {
		endTitle = endTitle - util.TrimRightSpaceLength(remainingLine[startTitle:endTitle])
		if endTitle > startTitle {
			node.Title = remainingLine[startTitle:endTitle]
		}
	}

	// ========================================================================== //
	// 	find attributes
	reader.Advance(endTitle)
	attrs, ok := parser.ParseAttributes(reader)
	hasClass := false
	admClass := bytes.Join([][]byte{[]byte("admonition adm-"), node.AdmonitionClass}, []byte(""))

	if ok {

		for _, attr := range attrs {
			oldVal := attr.Value.([]byte)
			var val []byte

			if bytes.Equal(attr.Name, []byte("class")) {
				hasClass = true
				val = bytes.Join([][]byte{admClass, oldVal}, []byte(" "))
			} else {
				val = oldVal
			}

			node.SetAttribute(attr.Name, val)
		}
	}

	if !hasClass {
		node.SetAttribute([]byte("class"), admClass)
	}

	return node
}

func (b *admonitionParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	rawdata := pc.Get(admonitionInfoKey)
	fdataMap := rawdata.([]*AdmonitionData)
	fdata := fdataMap[len(fdataMap)-1]

	line, segment := reader.PeekLine()
	w, pos := util.IndentWidth(line, reader.LineOffset())

	if close, newline := b.closes(line, segment, w, pos, node, fdata); close {
		reader.Advance(segment.Stop - segment.Start - newline + segment.Padding)
		fdataMap = fdataMap[:len(fdataMap)-1]

		if len(fdataMap) == 0 {
			return parser.Close
		} else {
			pc.Set(admonitionInfoKey, fdataMap)
			return parser.Close
		}
	}

	return parser.Continue | parser.HasChildren
}

func (b *admonitionParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
}

func (b *admonitionParser) CanInterruptParagraph() bool {
	return true
}

func (b *admonitionParser) CanAcceptIndentedLine() bool {
	return false
}

func (b *admonitionParser) closes(line []byte, segment text.Segment, w int, pos int, node ast.Node, fdata *AdmonitionData) (bool, int) {

	// don't close anything but the last node
	if node != fdata.node {
		return false, 1
	}

	// If the indentation is lower, we assume the user forgot to close the block
	if w < fdata.indent {
		return true, 1
	}

	// else, check for the correct number of closing chars and provide the info
	// necessary to advance the reader
	if w == fdata.indent {
		i := pos
		for ; i < len(line) && line[i] == fdata.char; i++ {
		}
		length := i - pos

		if length >= fdata.length && util.IsBlank(line[i:]) {
			newline := 1
			if line[len(line)-1] != '\n' {
				newline = 0
			}

			return true, newline
		}
	}

	return false, 0
}
