package varsmap

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	v := Varsmap{
		VarMappingImageBuilderLabelKey:             VarMappingImageBuilderLabelDefaultValue,
		VarMappingImageBuilderNameKey:              VarMappingImageBuilderNameDefaultValue,
		VarMappingImageBuilderRegistryHostKey:      VarMappingImageBuilderRegistryHostDefaultValue,
		VarMappingImageBuilderRegistryNamespaceKey: VarMappingImageBuilderRegistryNamespaceDefaultValue,
		VarMappingImageBuilderTagKey:               VarMappingImageBuilderTagDefaultValue,
		VarMappingImageExtraTagsKey:                VarMappingImageExtraTagsDefaultValue,
		VarMappingImageFromFullyQualifiedNameKey:   VarMappingImageFromFullyQualifiedNameValue,
		VarMappingImageFromNameKey:                 VarMappingImageFromNameDefaultValue,
		VarMappingImageFromRegistryHostKey:         VarMappingImageFromRegistryHostDefaultValue,
		VarMappingImageFromRegistryNamespaceKey:    VarMappingImageFromRegistryNamespaceDefaultValue,
		VarMappingImageFromTagKey:                  VarMappingImageFromTagDefaultValue,
		VarMappingImageFullyQualifiedNameKey:       VarMappingImageFullyQualifiedNameValue,
		VarMappingImageLabelsKey:                   VarMappingImageLabelsDefaultValue,
		VarMappingImageNameKey:                     VarMappingImageNameDefaultValue,
		VarMappingImageTagKey:                      VarMappingImageTagDefaultValue,
		VarMappingPullParentImageKey:               VarMappingPullParentImageDefaultValue,
		VarMappingPushImagetKey:                    VarMappingPushImagetDefaultValue,
		VarMappingRegistryHostKey:                  VarMappingRegistryHostDefaultValue,
		VarMappingRegistryNamespaceKey:             VarMappingRegistryNamespaceDefaultValue,
	}

	t.Log("Testing create a new varsmap")
	assert.Equal(t, v, New())
}

func TestGetUnderlyingMap(t *testing.T) {
	a := New()
	expected := map[string]string{
		VarMappingImageBuilderLabelKey:             VarMappingImageBuilderLabelDefaultValue,
		VarMappingImageBuilderNameKey:              VarMappingImageBuilderNameDefaultValue,
		VarMappingImageBuilderRegistryHostKey:      VarMappingImageBuilderRegistryHostDefaultValue,
		VarMappingImageBuilderRegistryNamespaceKey: VarMappingImageBuilderRegistryNamespaceDefaultValue,
		VarMappingImageBuilderTagKey:               VarMappingImageBuilderTagDefaultValue,
		VarMappingImageExtraTagsKey:                VarMappingImageExtraTagsDefaultValue,
		VarMappingImageFromFullyQualifiedNameKey:   VarMappingImageFromFullyQualifiedNameValue,
		VarMappingImageFromNameKey:                 VarMappingImageFromNameDefaultValue,
		VarMappingImageFromRegistryHostKey:         VarMappingImageFromRegistryHostDefaultValue,
		VarMappingImageFromRegistryNamespaceKey:    VarMappingImageFromRegistryNamespaceDefaultValue,
		VarMappingImageFromTagKey:                  VarMappingImageFromTagDefaultValue,
		VarMappingImageFullyQualifiedNameKey:       VarMappingImageFullyQualifiedNameValue,
		VarMappingImageLabelsKey:                   VarMappingImageLabelsDefaultValue,
		VarMappingImageNameKey:                     VarMappingImageNameDefaultValue,
		VarMappingImageTagKey:                      VarMappingImageTagDefaultValue,
		VarMappingPullParentImageKey:               VarMappingPullParentImageDefaultValue,
		VarMappingPushImagetKey:                    VarMappingPushImagetDefaultValue,
		VarMappingRegistryHostKey:                  VarMappingRegistryHostDefaultValue,
		VarMappingRegistryNamespaceKey:             VarMappingRegistryNamespaceDefaultValue,
	}

	t.Log("Testing get underlying map")
	underlyingMapA := a.GetUnderlyingMap()
	assert.Equal(t, expected, underlyingMapA)

}

func TestSetUnderlyingMap(t *testing.T) {

	a := New()
	newMapA := map[string]string{
		VarMappingImageBuilderLabelKey:             "expectedVarMappingImageBuilderLabelDefaultValue",
		VarMappingImageBuilderNameKey:              "expectedVarMappingImageBuilderNameDefaultValue",
		VarMappingImageBuilderRegistryHostKey:      "expectedVarMappingImageBuilderRegistryHostDefaultValue",
		VarMappingImageBuilderRegistryNamespaceKey: "expectedVarMappingImageBuilderRegistryNamespaceDefaultValue",
		VarMappingImageBuilderTagKey:               "expectedVarMappingImageBuilderTagDefaultValue",
		VarMappingImageExtraTagsKey:                "expectedVarMappingImageExtraTagsDefaultValue",
		VarMappingImageFromFullyQualifiedNameKey:   "expectedVarMappingImageFromFullyQualifiedNameValue",
		VarMappingImageFromNameKey:                 "expectedVarMappingImageFromNameDefaultValue",
		VarMappingImageFromRegistryHostKey:         "expectedVarMappingImageFromRegistryHostDefaultValue",
		VarMappingImageFromRegistryNamespaceKey:    "expectedVarMappingImageFromRegistryNamespaceDefaultValue",
		VarMappingImageFromTagKey:                  "expectedVarMappingImageFromTagDefaultValue",
		VarMappingImageFullyQualifiedNameKey:       "expectedVarMappingImageFullyQualifiedNameValue",
		VarMappingImageLabelsKey:                   "expectedVarMappingImageLabelsDefaultValue",
		VarMappingImageNameKey:                     "expectedVarMappingImageNameDefaultValue",
		VarMappingImageTagKey:                      "expectedVarMappingImageTagDefaultValue",
		VarMappingPullParentImageKey:               "expectedVarMappingPullParentImageDefaultValue",
		VarMappingPushImagetKey:                    "expectedVarMappingPushImagetDefaultValue",
		VarMappingRegistryHostKey:                  "expectedVarMappingRegistryHostDefaultValue",
		VarMappingRegistryNamespaceKey:             "expectedVarMappingRegistryNamespaceDefaultValue",
	}
	expected := Varsmap{
		VarMappingImageBuilderLabelKey:             "expectedVarMappingImageBuilderLabelDefaultValue",
		VarMappingImageBuilderNameKey:              "expectedVarMappingImageBuilderNameDefaultValue",
		VarMappingImageBuilderRegistryHostKey:      "expectedVarMappingImageBuilderRegistryHostDefaultValue",
		VarMappingImageBuilderRegistryNamespaceKey: "expectedVarMappingImageBuilderRegistryNamespaceDefaultValue",
		VarMappingImageBuilderTagKey:               "expectedVarMappingImageBuilderTagDefaultValue",
		VarMappingImageExtraTagsKey:                "expectedVarMappingImageExtraTagsDefaultValue",
		VarMappingImageFromFullyQualifiedNameKey:   "expectedVarMappingImageFromFullyQualifiedNameValue",
		VarMappingImageFromNameKey:                 "expectedVarMappingImageFromNameDefaultValue",
		VarMappingImageFromRegistryHostKey:         "expectedVarMappingImageFromRegistryHostDefaultValue",
		VarMappingImageFromRegistryNamespaceKey:    "expectedVarMappingImageFromRegistryNamespaceDefaultValue",
		VarMappingImageFromTagKey:                  "expectedVarMappingImageFromTagDefaultValue",
		VarMappingImageFullyQualifiedNameKey:       "expectedVarMappingImageFullyQualifiedNameValue",
		VarMappingImageLabelsKey:                   "expectedVarMappingImageLabelsDefaultValue",
		VarMappingImageNameKey:                     "expectedVarMappingImageNameDefaultValue",
		VarMappingImageTagKey:                      "expectedVarMappingImageTagDefaultValue",
		VarMappingPullParentImageKey:               "expectedVarMappingPullParentImageDefaultValue",
		VarMappingPushImagetKey:                    "expectedVarMappingPushImagetDefaultValue",
		VarMappingRegistryHostKey:                  "expectedVarMappingRegistryHostDefaultValue",
		VarMappingRegistryNamespaceKey:             "expectedVarMappingRegistryNamespaceDefaultValue",
	}

	t.Log("Testing set underlying map")
	a.SetUnderlyingMap(newMapA)
	assert.Equal(t, expected, a)
}

func TestCombine(t *testing.T) {

	t.Run("Testing Combine case A", func(t *testing.T) {
		a := New()
		newMapA := Varsmap{
			VarMappingImageBuilderNameKey: "expectedVarMappingImageBuilderNameDefaultValue",
			VarMappingImageBuilderTagKey:  "expectedVarMappingImageBuilderTagDefaultValue",
		}
		expected := Varsmap{
			VarMappingImageBuilderLabelKey:             VarMappingImageBuilderLabelDefaultValue,
			VarMappingImageBuilderNameKey:              "expectedVarMappingImageBuilderNameDefaultValue",
			VarMappingImageBuilderRegistryHostKey:      VarMappingImageBuilderRegistryHostDefaultValue,
			VarMappingImageBuilderRegistryNamespaceKey: VarMappingImageBuilderRegistryNamespaceDefaultValue,
			VarMappingImageBuilderTagKey:               "expectedVarMappingImageBuilderTagDefaultValue",
			VarMappingImageExtraTagsKey:                VarMappingImageExtraTagsDefaultValue,
			VarMappingImageFromFullyQualifiedNameKey:   VarMappingImageFromFullyQualifiedNameValue,
			VarMappingImageFromNameKey:                 VarMappingImageFromNameDefaultValue,
			VarMappingImageFromRegistryHostKey:         VarMappingImageFromRegistryHostDefaultValue,
			VarMappingImageFromRegistryNamespaceKey:    VarMappingImageFromRegistryNamespaceDefaultValue,
			VarMappingImageFromTagKey:                  VarMappingImageFromTagDefaultValue,
			VarMappingImageFullyQualifiedNameKey:       VarMappingImageFullyQualifiedNameValue,
			VarMappingImageLabelsKey:                   VarMappingImageLabelsDefaultValue,
			VarMappingImageNameKey:                     VarMappingImageNameDefaultValue,
			VarMappingImageTagKey:                      VarMappingImageTagDefaultValue,
			VarMappingPullParentImageKey:               VarMappingPullParentImageDefaultValue,
			VarMappingPushImagetKey:                    VarMappingPushImagetDefaultValue,
			VarMappingRegistryHostKey:                  VarMappingRegistryHostDefaultValue,
			VarMappingRegistryNamespaceKey:             VarMappingRegistryNamespaceDefaultValue,
		}

		t.Log("Testing combine case A")
		newMapA.Combine(a)
		assert.Equal(t, expected, newMapA)
	})

	t.Run("Testing Combine case B", func(t *testing.T) {
		b := New()
		newMapB := Varsmap{
			VarMappingRegistryHostKey:      "expectedVarMappingRegistryHostDefaultValue",
			VarMappingRegistryNamespaceKey: "expectedVarMappingRegistryNamespaceDefaultValue",
		}
		expected := Varsmap{
			VarMappingImageBuilderLabelKey:             VarMappingImageBuilderLabelDefaultValue,
			VarMappingImageBuilderNameKey:              VarMappingImageBuilderNameDefaultValue,
			VarMappingImageBuilderRegistryHostKey:      VarMappingImageBuilderRegistryHostDefaultValue,
			VarMappingImageBuilderRegistryNamespaceKey: VarMappingImageBuilderRegistryNamespaceDefaultValue,
			VarMappingImageBuilderTagKey:               VarMappingImageBuilderTagDefaultValue,
			VarMappingImageExtraTagsKey:                VarMappingImageExtraTagsDefaultValue,
			VarMappingImageFromFullyQualifiedNameKey:   VarMappingImageFromFullyQualifiedNameValue,
			VarMappingImageFromNameKey:                 VarMappingImageFromNameDefaultValue,
			VarMappingImageFromRegistryHostKey:         VarMappingImageFromRegistryHostDefaultValue,
			VarMappingImageFromRegistryNamespaceKey:    VarMappingImageFromRegistryNamespaceDefaultValue,
			VarMappingImageFromTagKey:                  VarMappingImageFromTagDefaultValue,
			VarMappingImageFullyQualifiedNameKey:       VarMappingImageFullyQualifiedNameValue,
			VarMappingImageLabelsKey:                   VarMappingImageLabelsDefaultValue,
			VarMappingImageNameKey:                     VarMappingImageNameDefaultValue,
			VarMappingImageTagKey:                      VarMappingImageTagDefaultValue,
			VarMappingPullParentImageKey:               VarMappingPullParentImageDefaultValue,
			VarMappingPushImagetKey:                    VarMappingPushImagetDefaultValue,
			VarMappingRegistryHostKey:                  "expectedVarMappingRegistryHostDefaultValue",
			VarMappingRegistryNamespaceKey:             "expectedVarMappingRegistryNamespaceDefaultValue",
		}

		t.Log("Testing combine case B")
		newMapB.Combine(b)
		assert.Equal(t, expected, newMapB)
	})

	t.Run("Testing error on Combine when input vars map is nil", func(t *testing.T) {
		var b Varsmap
		newMapB := New()

		errContext := "(core::domain::varsmap::Combine)"

		t.Log("Testing error on Combine when input vars map is nil")
		err := newMapB.Combine(b)
		assert.Error(t, errors.New(errContext, "Variables mapping to combine is nil"), err)
	})
}
