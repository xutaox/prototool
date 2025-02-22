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
	"github.com/xutaox/prototool/internal/strs"
	"github.com/xutaox/prototool/internal/text"
)

var enumFieldPrefixesLinter = NewLinter(
	"ENUM_FIELD_PREFIXES",
	"Verifies that all enum fields are prefixed with [NESTED_MESSAGE_NAME_]ENUM_NAME_.",
	checkEnumFieldPrefixes,
)

func checkEnumFieldPrefixes(add func(*text.Failure), dirPath string, descriptors []*FileDescriptor) error {
	return runVisitor(&enumFieldPrefixesVisitor{baseAddVisitor: newBaseAddVisitor(add)}, descriptors)
}

type enumFieldPrefixesVisitor struct {
	baseAddVisitor

	nestedNames []string
}

func (v *enumFieldPrefixesVisitor) VisitMessage(message *proto.Message) {
	v.nestedNames = append(v.nestedNames, strs.ToUpperSnakeCase(message.Name))
	for _, child := range message.Elements {
		child.Accept(v)
	}
	v.nestedNames = v.nestedNames[0 : len(v.nestedNames)-1]
}

func (v *enumFieldPrefixesVisitor) VisitEnum(enum *proto.Enum) {
	v.nestedNames = append(v.nestedNames, strs.ToUpperSnakeCase(enum.Name))
	for _, child := range enum.Elements {
		child.Accept(v)
	}
	v.nestedNames = v.nestedNames[0 : len(v.nestedNames)-1]
}

func (v *enumFieldPrefixesVisitor) VisitEnumField(enumField *proto.EnumField) {
	expectedPrefix := strings.Join(v.nestedNames, "_") + "_"
	if !strings.HasPrefix(enumField.Name, expectedPrefix) {
		v.AddFailuref(enumField.Position, "Enum field %q is expected to have the prefix %q.", enumField.Name, expectedPrefix)
	}
}
