package ast

import (
	"fmt"
	"strings"

	"github.com/alsritter/walk/walk"
)

type ASTree interface {
	// Get i th ASTree.
	Child(i int) ASTree
	// The Child number.
	NumChildren() int
	// return ASTree children Iterator.
	Children() *Iterator
	// return current location token string.
	Location() string

	ToString() string
}

type Ints []ASTree

func (i Ints) Iterator() *Iterator {
	return &Iterator{
		data:  i,
		index: 0,
	}
}

type Iterator struct {
	data  Ints
	index int
}

func (i *Iterator) HasNext() bool {
	return i.index < len(i.data)
}

func (i *Iterator) Next() (v ASTree) {
	v = i.data[i.index]
	i.index++
	return v
}

// ==========================ASTLeaf===============================
var empty = make(Ints, 0)

type ASTLeaf struct {
	token walk.Token
}

func NewASTLeaf(t walk.Token) ASTree {
	return &ASTLeaf{token: t}
}

func (a *ASTLeaf) Children() *Iterator {
	return empty.Iterator()
}

func (a *ASTLeaf) Child(i int) ASTree {
	walk.PanicError(walk.NewIndexOutOfBoundsException(fmt.Sprintf("%v %d", a.token, i), nil))
	return nil
}

func (a *ASTLeaf) NumChildren() int { return 0 }

func (a *ASTLeaf) Location() string { return fmt.Sprintf("at line %d", a.token.GetLineNumber()) }

func (a *ASTLeaf) ToString() string { return a.token.GetText() }

func (a *ASTLeaf) Token() walk.Token { return a.token }

// ============================ASTList=============================

type ASTList struct {
	children []ASTree
}

func NewASTList(list []ASTree) ASTree {
	return &ASTList{children: list}
}

func (a *ASTList) NumChildren() int { return len(a.children) }

func (a *ASTList) Child(i int) ASTree { return a.children[i] }

func (a *ASTList) Children() *Iterator {
	return Ints(a.children).Iterator()
}

func (a *ASTList) Location() string {
	for _, t := range a.children {
		s := t.Location()
		if s != "" {
			return s
		}
	}
	return ""
}

func (a *ASTList) ToString() string {
	sb := strings.Builder{}
	sb.WriteString("(")
	sep := ""
	for _, t := range a.children {
		sb.WriteString(sep)
		sep = " "
		sb.WriteString(t.ToString())
	}
	return sb.String()
}

// ====================BinaryExpr=========================
type BinaryExpr struct {
	ASTree
}

func NewBinaryExpr(list []ASTree) ASTree {
	return &BinaryExpr{
		ASTree: NewASTList(list),
	}
}

func (e *BinaryExpr) Left() ASTree {
	return e.Child(0)
}

func (e *BinaryExpr) Right() ASTree {
	return e.Child(2)
}

func (e *BinaryExpr) Operator() string {
	leaf := e.Child(1).(*ASTLeaf)
	return leaf.token.GetText()
}

// =====================NumberLiteral=====================
type NumberLiteral struct {
	ASTree
}

func NewNumberLiteral(t walk.Token) ASTree {
	return &NumberLiteral{
		ASTree: NewASTLeaf(t),
	}
}

func (e *NumberLiteral) Value() int32 {
	t := e.ASTree.(*ASTLeaf)
	return t.Token().GetNumber()
}

// ==========================Name=============================
type Name struct {
	ASTree
}

func NewName(t walk.Token) ASTree {
	return &Name{
		ASTree: NewASTLeaf(t),
	}
}

func (e *Name) Name() string {
	t := e.ASTree.(*ASTLeaf)
	return t.Token().GetText()
}
