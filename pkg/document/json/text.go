/*
 * Copyright 2020 The Yorkie Authors. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package json

import (
	"fmt"
	"unicode/utf8"

	"github.com/yorkie-team/yorkie/pkg/document/time"
	"github.com/yorkie-team/yorkie/pkg/log"
)

func InitialTextNode() *RGATreeSplitNode {
	return NewRGATreeSplitNode(initialNodeID, &TextValue{
		value: "",
	})
}

type TextValue struct {
	value string
}

func NewTextValue(value string) *TextValue {
	return &TextValue{
		value: value,
	}
}

func (t *TextValue) Len() int {
	return utf8.RuneCountInString(t.value)
}

func (t *TextValue) String() string {
	return t.value
}

func (t *TextValue) AnnotatedString() string {
	return t.value
}

func (t *TextValue) Split(offset int) RGATreeSplitValue {
	value := t.value
	r := []rune(value)
	t.value = string(r[0:offset])
	return NewTextValue(string(r[offset:]))
}

// DeepCopy copies itself deeply.
func (t *TextValue) DeepCopy() RGATreeSplitValue {
	return &TextValue{
		value: t.value,
	}
}

// Text is an extended data type for the contents of a text editor.
type Text struct {
	rgaTreeSplit *RGATreeSplit
	selectionMap map[string]*Selection
	createdAt    *time.Ticket
	movedAt      *time.Ticket
	removedAt    *time.Ticket
}

// NewText creates a new instance of Text.
func NewText(elements *RGATreeSplit, createdAt *time.Ticket) *Text {
	return &Text{
		rgaTreeSplit: elements,
		selectionMap: make(map[string]*Selection),
		createdAt:    createdAt,
	}
}

func (t *Text) Marshal() string {
	return fmt.Sprintf("\"%s\"", t.rgaTreeSplit.marshal())
}

// DeepCopy copies itself deeply.
func (t *Text) DeepCopy() Element {
	rgaTreeSplit := NewRGATreeSplit(InitialTextNode())

	current := rgaTreeSplit.InitialHead()
	for _, node := range t.Nodes() {
		current = rgaTreeSplit.InsertAfter(current, node.DeepCopy())
		insPrevID := node.InsPrevID()
		if insPrevID != nil {
			insPrevNode := rgaTreeSplit.FindNode(insPrevID)
			if insPrevNode == nil {
				log.Logger.Warn("insPrevNode should be presence")
			}
			current.SetInsPrev(insPrevNode)
		}
	}

	return NewText(rgaTreeSplit, t.createdAt)
}

// CreatedAt returns the creation time of this Text.
func (t *Text) CreatedAt() *time.Ticket {
	return t.createdAt
}

// RemovedAt returns the removal time of this Text.
func (t *Text) RemovedAt() *time.Ticket {
	return t.removedAt
}

// MovedAt returns the move time of this Text.
func (t *Text) MovedAt() *time.Ticket {
	return t.movedAt
}

// SetMovedAt sets the move time of this Text.
func (t *Text) SetMovedAt(movedAt *time.Ticket) {
	t.movedAt = movedAt
}

// Remove removes this Text.
func (t *Text) Remove(executedAt *time.Ticket) bool {
	if t.removedAt == nil || executedAt.After(t.removedAt) {
		t.removedAt = executedAt
		return true
	}
	return false
}

// CreateRange returns pair of RGATreeSplitNodePos of the given integer offsets.
func (t *Text) CreateRange(from, to int) (*RGATreeSplitNodePos, *RGATreeSplitNodePos) {
	return t.rgaTreeSplit.createRange(from, to)
}

func (t *Text) Edit(
	from,
	to *RGATreeSplitNodePos,
	latestCreatedAtMapByActor map[string]*time.Ticket,
	content string,
	executedAt *time.Ticket,
) (*RGATreeSplitNodePos, map[string]*time.Ticket) {
	cursorPos, latestCreatedAtMapByActor := t.rgaTreeSplit.edit(
		from,
		to,
		latestCreatedAtMapByActor,
		NewTextValue(content),
		executedAt,
	)
	log.Logger.Debugf(
		"EDIT: '%s' edits %s",
		executedAt.ActorID().String(),
		t.rgaTreeSplit.AnnotatedString(),
	)
	return cursorPos, latestCreatedAtMapByActor
}

func (t *Text) Select(
	from *RGATreeSplitNodePos,
	to *RGATreeSplitNodePos,
	executedAt *time.Ticket,
) {
	if _, ok := t.selectionMap[executedAt.ActorIDHex()]; !ok {
		t.selectionMap[executedAt.ActorIDHex()] = newSelection(from, to, executedAt)
		return
	}

	prevSelection := t.selectionMap[executedAt.ActorIDHex()]
	if executedAt.After(prevSelection.updatedAt) {
		log.Logger.Debugf(
			"SELT: '%s' selects %s",
			executedAt.ActorID().String(),
			t.rgaTreeSplit.AnnotatedString(),
		)

		t.selectionMap[executedAt.ActorIDHex()] = newSelection(from, to, executedAt)
	}
}

func (t *Text) Nodes() []*RGATreeSplitNode {
	return t.rgaTreeSplit.nodes()
}

// AnnotatedString returns a String containing the meta data of the text
// for debugging purpose.
func (t *Text) AnnotatedString() string {
	return t.rgaTreeSplit.AnnotatedString()
}
