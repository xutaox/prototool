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
	"strings"

	"github.com/emicklei/proto"
	"github.com/xutaox/prototool/internal/text"
)

var messageFieldNamesFilepathLinter = NewLinter(
	"MESSAGE_FIELD_NAMES_FILEPATH",
	`Verifies that all message field names do not contain "file_path" as "filepath" should be used instead.`,
	checkMessageFieldNamesFilepath,
)

func checkMessageFieldNamesFilepath(add func(*text.Failure), dirPath string, descriptors []*FileDescriptor) error {
	return runVisitor(messageFieldNamesFilepathVisitor{baseAddVisitor: newBaseAddVisitor(add)}, descriptors)
}

type messageFieldNamesFilepathVisitor struct {
	baseAddVisitor
}

func (v messageFieldNamesFilepathVisitor) VisitMessage(message *proto.Message) {
	for _, element := range message.Elements {
		element.Accept(v)
	}
}

func (v messageFieldNamesFilepathVisitor) VisitOneof(oneof *proto.Oneof) {
	for _, element := range oneof.Elements {
		element.Accept(v)
	}
}

func (v messageFieldNamesFilepathVisitor) VisitNormalField(field *proto.NormalField) {
	v.handleField(field.Field)
}

func (v messageFieldNamesFilepathVisitor) VisitOneofField(field *proto.OneOfField) {
	v.handleField(field.Field)
}

func (v messageFieldNamesFilepathVisitor) VisitMapField(field *proto.MapField) {
	v.handleField(field.Field)
}

func (v messageFieldNamesFilepathVisitor) handleField(field *proto.Field) {
	if strings.Contains(field.Name, "file_path") {
		// This might break down in practice but let's try it out
		v.AddFailuref(field.Position, "Field name %q should be %q.", field.Name, strings.Replace(field.Name, "file_path", "filepath", -1))
	}
}
