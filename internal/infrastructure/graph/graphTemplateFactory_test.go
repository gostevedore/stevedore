package graph

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGraphTemplate(t *testing.T) {
	tests := []struct {
		desc    string
		factory *GraphTemplateFactory
		res     GraphTemplater
	}{
		{
			desc:    "Testing create new GraphTemplate",
			factory: NewGraphTemplateFactory(false),
			res:     NewGraphTemplate(),
		},
		{
			desc:    "Testing create new MockGraphTemplate",
			factory: NewGraphTemplateFactory(true),
			res:     NewMockGraphTemplate(),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res := test.factory.NewGraphTemplate()
			assert.Equal(t, reflect.TypeOf(test.res), reflect.TypeOf(res))
		})
	}
}

func TestNewGraphTemplateNode(t *testing.T) {
	tests := []struct {
		desc    string
		factory *GraphTemplateFactory
		res     GraphTemplateNoder
	}{
		{
			desc:    "Testing create new GraphTemplateNode",
			factory: NewGraphTemplateFactory(false),
			res:     NewGraphTemplateNode("node"),
		},
		{
			desc:    "Testing create new MockGraphTemplateNode",
			factory: NewGraphTemplateFactory(true),
			res:     NewGraphTemplateNode("node"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res := test.factory.NewGraphTemplateNode("node")
			assert.Equal(t, reflect.TypeOf(test.res), reflect.TypeOf(res))
		})
	}
}
