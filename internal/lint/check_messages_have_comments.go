// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package lint

import (
	"github.com/emicklei/proto"
	"github.com/xutaox/prototool/internal/text"
)

var messagesHaveCommentsLinter = NewLinter(
	"MESSAGES_HAVE_COMMENTS",
	`Verifies that all non-extended messages have a comment of the form "// MessageName ...".`,
	checkMessagesHaveComments,
)

func checkMessagesHaveComments(add func(*text.Failure), dirPath string, descriptors []*FileDescriptor) error {
	return runVisitor(messagesHaveCommentsVisitor{baseAddVisitor: newBaseAddVisitor(add)}, descriptors)
}

type messagesHaveCommentsVisitor struct {
	baseAddVisitor
}

func (v messagesHaveCommentsVisitor) VisitMessage(message *proto.Message) {
	// for nested messages
	for _, child := range message.Elements {
		child.Accept(v)
	}
	if message.IsExtend {
		return
	}
	if !hasGolangStyleComment(message.Comment, message.Name) {
		v.AddFailuref(message.Position, `Message %q needs a comment of the form "// %s ..."`, message.Name, message.Name)
	}
}
