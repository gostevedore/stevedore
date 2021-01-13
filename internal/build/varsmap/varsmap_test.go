package varsmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	v := Varsmap{
		VarMappingImageBuilderNameKey:              VarMappingImageBuilderNameDefaultValue,
		VarMappingImageBuilderTagKey:               VarMappingImageBuilderTagDefaultValue,
		VarMappingImageBuilderRegistryNamespaceKey: VarMappingImageBuilderRegistryNamespaceDefaultValue,
		VarMappingImageBuilderRegistryHostKey:      VarMappingImageBuilderRegistryHostDefaultValue,
		VarMappingImageBuilderLabelKey:             VarMappingImageBuilderLabelDefaultValue,
		VarMappingImageFromNameKey:                 VarMappingImageFromNameDefaultValue,
		VarMappingImageFromTagKey:                  VarMappingImageFromTagDefaultValue,
		VarMappingImageFromRegistryNamespaceKey:    VarMappingImageFromRegistryNamespaceDefaultValue,
		VarMappingImageFromRegistryHostKey:         VarMappingImageFromRegistryHostDefaultValue,
		VarMappingImageNameKey:                     VarMappingImageNameDefaultValue,
		VarMappingImageTagKey:                      VarMappingImageTagDefaultValue,
		VarMappingImageExtraTagsKey:                VarMappingImageExtraTagsDefaultValue,
		VarMappingRegistryNamespaceKey:             VarMappingRegistryNamespaceDefaultValue,
		VarMappingRegistryHostKey:                  VarMappingRegistryHostDefaultValue,
		VarMappingPushImagetKey:                    VarMappingPushImagetDefaultValue,
	}

	t.Log("Testing create a new varsmap")
	assert.Equal(t, v, New())
}

func TestGetUnderlyingMap(t *testing.T) {
	a := New()
	expected := map[string]string{
		VarMappingImageBuilderNameKey:              VarMappingImageBuilderNameDefaultValue,
		VarMappingImageBuilderTagKey:               VarMappingImageBuilderTagDefaultValue,
		VarMappingImageBuilderRegistryNamespaceKey: VarMappingImageBuilderRegistryNamespaceDefaultValue,
		VarMappingImageBuilderRegistryHostKey:      VarMappingImageBuilderRegistryHostDefaultValue,
		VarMappingImageBuilderLabelKey:             VarMappingImageBuilderLabelDefaultValue,
		VarMappingImageFromNameKey:                 VarMappingImageFromNameDefaultValue,
		VarMappingImageFromTagKey:                  VarMappingImageFromTagDefaultValue,
		VarMappingImageFromRegistryNamespaceKey:    VarMappingImageFromRegistryNamespaceDefaultValue,
		VarMappingImageFromRegistryHostKey:         VarMappingImageFromRegistryHostDefaultValue,
		VarMappingImageNameKey:                     VarMappingImageNameDefaultValue,
		VarMappingImageTagKey:                      VarMappingImageTagDefaultValue,
		VarMappingImageExtraTagsKey:                VarMappingImageExtraTagsDefaultValue,
		VarMappingRegistryNamespaceKey:             VarMappingRegistryNamespaceDefaultValue,
		VarMappingRegistryHostKey:                  VarMappingRegistryHostDefaultValue,
		VarMappingPushImagetKey:                    VarMappingPushImagetDefaultValue,
	}

	t.Log("Testing get underlying map")
	underlyingMapA := a.GetUnderlyingMap()
	assert.Equal(t, expected, underlyingMapA)

}

func TestSetUnderlyingMap(t *testing.T) {

	a := New()
	newMapA := map[string]string{
		VarMappingImageBuilderNameKey:              "expectedVarMappingImageBuilderNameDefaultValue",
		VarMappingImageBuilderTagKey:               "expectedVarMappingImageBuilderTagDefaultValue",
		VarMappingImageBuilderRegistryNamespaceKey: "expectedVarMappingImageBuilderRegistryNamespaceDefaultValue",
		VarMappingImageBuilderRegistryHostKey:      "expectedVarMappingImageBuilderRegistryHostDefaultValue",
		VarMappingImageBuilderLabelKey:             "expectedVarMappingImageBuilderLabelDefaultValue",
		VarMappingImageFromNameKey:                 "expectedVarMappingImageFromNameDefaultValue",
		VarMappingImageFromTagKey:                  "expectedVarMappingImageFromTagDefaultValue",
		VarMappingImageFromRegistryNamespaceKey:    "expectedVarMappingImageFromRegistryNamespaceDefaultValue",
		VarMappingImageFromRegistryHostKey:         "expectedVarMappingImageFromRegistryHostDefaultValue",
		VarMappingImageNameKey:                     "expectedVarMappingImageNameDefaultValue",
		VarMappingImageTagKey:                      "expectedVarMappingImageTagDefaultValue",
		VarMappingImageExtraTagsKey:                "expectedVarMappingImageExtraTagsDefaultValue",
		VarMappingRegistryNamespaceKey:             "expectedVarMappingRegistryNamespaceDefaultValue",
		VarMappingRegistryHostKey:                  "expectedVarMappingRegistryHostDefaultValue",
		VarMappingPushImagetKey:                    "expectedVarMappingPushImagetDefaultValue",
	}
	expected := Varsmap{
		VarMappingImageBuilderNameKey:              "expectedVarMappingImageBuilderNameDefaultValue",
		VarMappingImageBuilderTagKey:               "expectedVarMappingImageBuilderTagDefaultValue",
		VarMappingImageBuilderRegistryNamespaceKey: "expectedVarMappingImageBuilderRegistryNamespaceDefaultValue",
		VarMappingImageBuilderRegistryHostKey:      "expectedVarMappingImageBuilderRegistryHostDefaultValue",
		VarMappingImageBuilderLabelKey:             "expectedVarMappingImageBuilderLabelDefaultValue",
		VarMappingImageFromNameKey:                 "expectedVarMappingImageFromNameDefaultValue",
		VarMappingImageFromTagKey:                  "expectedVarMappingImageFromTagDefaultValue",
		VarMappingImageFromRegistryNamespaceKey:    "expectedVarMappingImageFromRegistryNamespaceDefaultValue",
		VarMappingImageFromRegistryHostKey:         "expectedVarMappingImageFromRegistryHostDefaultValue",
		VarMappingImageNameKey:                     "expectedVarMappingImageNameDefaultValue",
		VarMappingImageTagKey:                      "expectedVarMappingImageTagDefaultValue",
		VarMappingImageExtraTagsKey:                "expectedVarMappingImageExtraTagsDefaultValue",
		VarMappingRegistryNamespaceKey:             "expectedVarMappingRegistryNamespaceDefaultValue",
		VarMappingRegistryHostKey:                  "expectedVarMappingRegistryHostDefaultValue",
		VarMappingPushImagetKey:                    "expectedVarMappingPushImagetDefaultValue",
	}

	t.Log("Testing set underlying map")
	a.SetUnderlyingMap(newMapA)
	assert.Equal(t, expected, a)
}

func TestCombine(t *testing.T) {
	a := New()
	newMapA := Varsmap{
		VarMappingImageBuilderNameKey: "expectedVarMappingImageBuilderNameDefaultValue",
		VarMappingImageBuilderTagKey:  "expectedVarMappingImageBuilderTagDefaultValue",
	}
	expected := Varsmap{
		VarMappingImageBuilderNameKey:              "expectedVarMappingImageBuilderNameDefaultValue",
		VarMappingImageBuilderTagKey:               "expectedVarMappingImageBuilderTagDefaultValue",
		VarMappingImageBuilderRegistryNamespaceKey: VarMappingImageBuilderRegistryNamespaceDefaultValue,
		VarMappingImageBuilderRegistryHostKey:      VarMappingImageBuilderRegistryHostDefaultValue,
		VarMappingImageBuilderLabelKey:             VarMappingImageBuilderLabelDefaultValue,
		VarMappingImageFromNameKey:                 VarMappingImageFromNameDefaultValue,
		VarMappingImageFromTagKey:                  VarMappingImageFromTagDefaultValue,
		VarMappingImageFromRegistryNamespaceKey:    VarMappingImageFromRegistryNamespaceDefaultValue,
		VarMappingImageFromRegistryHostKey:         VarMappingImageFromRegistryHostDefaultValue,
		VarMappingImageNameKey:                     VarMappingImageNameDefaultValue,
		VarMappingImageTagKey:                      VarMappingImageTagDefaultValue,
		VarMappingImageExtraTagsKey:                VarMappingImageExtraTagsDefaultValue,
		VarMappingRegistryNamespaceKey:             VarMappingRegistryNamespaceDefaultValue,
		VarMappingRegistryHostKey:                  VarMappingRegistryHostDefaultValue,
		VarMappingPushImagetKey:                    VarMappingPushImagetDefaultValue,
	}

	t.Log("Testing combine maps")
	newMapA.Combine(a)
	assert.Equal(t, expected, newMapA)
}
