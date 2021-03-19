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

package splay_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yorkie-team/yorkie/pkg/splay"
)

type stringValue struct {
	content string
}

func newSplayNode(content string) *splay.Node {
	return splay.NewNode(&stringValue{
		content: content,
	})
}

func (v *stringValue) Len() int {
	return len(v.content)
}

func (v *stringValue) String() string {
	return v.content
}

func TestSplayTree(t *testing.T) {
	t.Run("insert and splay test", func(t *testing.T) {
		tree := splay.NewTree(nil)

		nodeA := tree.Insert(newSplayNode("A2"))
		assert.Equal(t, "[2,2]A2", tree.AnnotatedString())
		nodeB := tree.Insert(newSplayNode("B23"))
		assert.Equal(t, "[2,2]A2[5,3]B23", tree.AnnotatedString())
		nodeC := tree.Insert(newSplayNode("C234"))
		assert.Equal(t, "[2,2]A2[5,3]B23[9,4]C234", tree.AnnotatedString())
		nodeD := tree.Insert(newSplayNode("D2345"))
		assert.Equal(t, "[2,2]A2[5,3]B23[9,4]C234[14,5]D2345", tree.AnnotatedString())

		tree.Splay(nodeB)
		assert.Equal(t, "[2,2]A2[14,3]B23[9,4]C234[5,5]D2345", tree.AnnotatedString())

		assert.Equal(t, 0, tree.IndexOf(nodeA))
		assert.Equal(t, 2, tree.IndexOf(nodeB))
		assert.Equal(t, 5, tree.IndexOf(nodeC))
		assert.Equal(t, 9, tree.IndexOf(nodeD))

		node, offset := tree.Find(1)
		assert.Equal(t, nodeA, node)
		assert.Equal(t, 1, offset)

		node, offset = tree.Find(7)
		assert.Equal(t, nodeC, node)
		assert.Equal(t, 2, offset)

		node, offset = tree.Find(11)
		assert.Equal(t, nodeD, node)
		assert.Equal(t, 2, offset)
	})

	t.Run("deletion test", func(t *testing.T) {
		tree := splay.NewTree(nil)

		nodeH := tree.Insert(newSplayNode("H"))
		assert.Equal(t, "[1,1]H", tree.AnnotatedString())
		nodeE := tree.Insert(newSplayNode("E"))
		assert.Equal(t, "[1,1]H[2,1]E", tree.AnnotatedString())
		nodeL := tree.Insert(newSplayNode("LL"))
		assert.Equal(t, "[1,1]H[2,1]E[4,2]LL", tree.AnnotatedString())
		nodeO := tree.Insert(newSplayNode("O"))
		assert.Equal(t, "[1,1]H[2,1]E[4,2]LL[5,1]O", tree.AnnotatedString())

		tree.Delete(nodeE)
		assert.Equal(t, "[4,1]H[3,2]LL[1,1]O", tree.AnnotatedString())

		assert.Equal(t, tree.IndexOf(nodeH), 0)
		assert.Equal(t, tree.IndexOf(nodeE), -1)
		assert.Equal(t, tree.IndexOf(nodeL), 1)
		assert.Equal(t, tree.IndexOf(nodeO), 3)
	})
}
