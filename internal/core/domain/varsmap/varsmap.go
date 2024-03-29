package varsmap

import (
	errors "github.com/apenella/go-common-utils/error"
)

const (
	VarMappingImageBuilderLabelKey             = "image_builder_label_key"
	VarMappingImageBuilderNameKey              = "image_builder_name_key"               // Not comming from build's command flag
	VarMappingImageBuilderRegistryHostKey      = "image_builder_registry_host_key"      // Not comming from build's command flag
	VarMappingImageBuilderRegistryNamespaceKey = "image_builder_registry_namespace_key" // Not comming from build's command flag
	VarMappingImageBuilderTagKey               = "image_builder_tag_key"                // Not comming from build's command flag
	VarMappingImageExtraTagsKey                = "image_extra_tags_key"
	VarMappingImageFromFullyQualifiedNameKey   = "image_from_fully_qualified_name_key"
	VarMappingImageFromNameKey                 = "image_from_name_key"
	VarMappingImageFromRegistryHostKey         = "image_from_registry_host_key"
	VarMappingImageFromRegistryNamespaceKey    = "image_from_registry_namespace_key"
	VarMappingImageFromTagKey                  = "image_from_tag_key"
	VarMappingImageFullyQualifiedNameKey       = "image_fully_qualified_name_key"
	VarMappingImageLabelsKey                   = "image_lables_key"
	VarMappingImageNameKey                     = "image_name_key"
	VarMappingImageTagKey                      = "image_tag_key"
	VarMappingPullParentImageKey               = "pull_parent_image_key"
	VarMappingPushImagetKey                    = "push_image_key"
	VarMappingRegistryHostKey                  = "image_registry_host_key"
	VarMappingRegistryNamespaceKey             = "image_registry_namespace_key"

	VarMappingImageBuilderLabelDefaultValue             = "image_builder_label"
	VarMappingImageBuilderNameDefaultValue              = "image_builder_name"               // Not comming from build's command flag
	VarMappingImageBuilderRegistryHostDefaultValue      = "image_builder_registry_host"      // Not comming from build's command flag
	VarMappingImageBuilderRegistryNamespaceDefaultValue = "image_builder_registry_namespace" // Not comming from build's command flag
	VarMappingImageBuilderTagDefaultValue               = "image_builder_tag"                // Not comming from build's command flag
	VarMappingImageExtraTagsDefaultValue                = "image_extra_tags"
	VarMappingImageFromFullyQualifiedNameValue          = "image_from_fully_qualified_name"
	VarMappingImageFromNameDefaultValue                 = "image_from_name"
	VarMappingImageFromRegistryHostDefaultValue         = "image_from_registry_host"
	VarMappingImageFromRegistryNamespaceDefaultValue    = "image_from_registry_namespace"
	VarMappingImageFromTagDefaultValue                  = "image_from_tag"
	VarMappingImageFullyQualifiedNameValue              = "image_fully_qualified_name"
	VarMappingImageLabelsDefaultValue                   = "image_labels"
	VarMappingImageNameDefaultValue                     = "image_name"
	VarMappingImageTagDefaultValue                      = "image_tag"
	VarMappingPullParentImageDefaultValue               = "pull_parent_image"
	VarMappingPushImagetDefaultValue                    = "push_image"
	VarMappingRegistryHostDefaultValue                  = "image_registry_host"
	VarMappingRegistryNamespaceDefaultValue             = "image_registry_namespace"
)

// Varsmap is a map[string]string that defines the variables names passed from builder to build drivers
type Varsmap map[string]string

// New return a Varsmap object
func New() Varsmap {
	return Varsmap{
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
}

// GetUnderlyingMap return the map[string]string behind Varsmap
func (v Varsmap) GetUnderlyingMap() map[string]string {
	return (map[string]string)(v)
}

// SetUnderlyingMap return the map[string]string behind Varsmap
func (v Varsmap) SetUnderlyingMap(underlyingmap map[string]string) {

	for key, value := range underlyingmap {
		v[key] = value
	}

}

// Combine include c varsmap values over v varsmsp but does not overrides values when a key already exists
func (v Varsmap) Combine(c Varsmap) error {
	var existsV, existsC bool

	errContext := "(core::domain::varsmap::Combine)"

	if v == nil {
		return errors.New(errContext, "Variables mapping is nil")
	}

	if c == nil {
		return errors.New(errContext, "Variables mapping to combine is nil")
	}

	auxV := v.GetUnderlyingMap()
	auxC := c.GetUnderlyingMap()

	_, existsV = auxV[VarMappingImageBuilderNameKey]
	if !existsV {
		_, existsC = auxC[VarMappingImageBuilderNameKey]
		if existsC {
			auxV[VarMappingImageBuilderNameKey] = auxC[VarMappingImageBuilderNameKey]
		}
	}

	_, existsV = auxV[VarMappingImageBuilderTagKey]
	if !existsV {
		_, existsC = auxC[VarMappingImageBuilderTagKey]
		if existsC {
			auxV[VarMappingImageBuilderTagKey] = auxC[VarMappingImageBuilderTagKey]
		}
	}

	_, existsV = auxV[VarMappingImageBuilderRegistryNamespaceKey]
	if !existsV {
		_, existsC = auxC[VarMappingImageBuilderRegistryNamespaceKey]
		if existsC {
			auxV[VarMappingImageBuilderRegistryNamespaceKey] = auxC[VarMappingImageBuilderRegistryNamespaceKey]
		}
	}

	_, existsV = auxV[VarMappingImageBuilderRegistryHostKey]
	if !existsV {
		_, existsC = auxC[VarMappingImageBuilderRegistryHostKey]
		if existsC {
			auxV[VarMappingImageBuilderRegistryHostKey] = auxC[VarMappingImageBuilderRegistryHostKey]
		}
	}

	_, existsV = auxV[VarMappingImageBuilderLabelKey]
	if !existsV {
		_, existsC = auxC[VarMappingImageBuilderLabelKey]
		if existsC {
			auxV[VarMappingImageBuilderLabelKey] = auxC[VarMappingImageBuilderLabelKey]
		}
	}

	_, existsV = auxV[VarMappingImageFromFullyQualifiedNameKey]
	if !existsV {
		_, existsC = auxC[VarMappingImageFromFullyQualifiedNameKey]
		if existsC {
			auxV[VarMappingImageFromFullyQualifiedNameKey] = auxC[VarMappingImageFromFullyQualifiedNameKey]
		}
	}

	_, existsV = auxV[VarMappingImageFromNameKey]
	if !existsV {
		_, existsC = auxC[VarMappingImageFromNameKey]
		if existsC {
			auxV[VarMappingImageFromNameKey] = auxC[VarMappingImageFromNameKey]
		}
	}

	_, existsV = auxV[VarMappingImageFullyQualifiedNameKey]
	if !existsV {
		_, existsC = auxC[VarMappingImageFullyQualifiedNameKey]
		if existsC {
			auxV[VarMappingImageFullyQualifiedNameKey] = auxC[VarMappingImageFullyQualifiedNameKey]
		}
	}

	_, existsV = auxV[VarMappingImageFromTagKey]
	if !existsV {
		_, existsC = auxC[VarMappingImageFromTagKey]
		if existsC {
			auxV[VarMappingImageFromTagKey] = auxC[VarMappingImageFromTagKey]
		}
	}

	_, existsV = auxV[VarMappingImageFromRegistryNamespaceKey]
	if !existsV {
		_, existsC = auxC[VarMappingImageFromRegistryNamespaceKey]
		if existsC {
			auxV[VarMappingImageFromRegistryNamespaceKey] = auxC[VarMappingImageFromRegistryNamespaceKey]
		}
	}

	_, existsV = auxV[VarMappingImageFromRegistryHostKey]
	if !existsV {
		_, existsC = auxC[VarMappingImageFromRegistryHostKey]
		if existsC {
			auxV[VarMappingImageFromRegistryHostKey] = auxC[VarMappingImageFromRegistryHostKey]
		}
	}

	_, existsV = auxV[VarMappingImageNameKey]
	if !existsV {
		_, existsC = auxC[VarMappingImageNameKey]
		if existsC {
			auxV[VarMappingImageNameKey] = auxC[VarMappingImageNameKey]
		}
	}

	_, existsV = auxV[VarMappingImageTagKey]
	if !existsV {
		_, existsC = auxC[VarMappingImageTagKey]
		if existsC {
			auxV[VarMappingImageTagKey] = auxC[VarMappingImageTagKey]
		}
	}

	_, existsV = auxV[VarMappingRegistryNamespaceKey]
	if !existsV {
		_, existsC = auxC[VarMappingRegistryNamespaceKey]
		if existsC {
			auxV[VarMappingRegistryNamespaceKey] = auxC[VarMappingRegistryNamespaceKey]
		}
	}

	_, existsV = auxV[VarMappingRegistryHostKey]
	if !existsV {
		_, existsC = auxC[VarMappingRegistryHostKey]
		if existsC {
			auxV[VarMappingRegistryHostKey] = auxC[VarMappingRegistryHostKey]
		}
	}

	_, existsV = auxV[VarMappingPullParentImageKey]
	if !existsV {
		_, existsC = auxC[VarMappingPullParentImageKey]
		if existsC {
			auxV[VarMappingPullParentImageKey] = auxC[VarMappingPullParentImageKey]
		}
	}

	_, existsV = auxV[VarMappingPushImagetKey]
	if !existsV {
		_, existsC = auxC[VarMappingPushImagetKey]
		if existsC {
			auxV[VarMappingPushImagetKey] = auxC[VarMappingPushImagetKey]
		}
	}

	_, existsV = auxV[VarMappingImageExtraTagsKey]
	if !existsV {
		_, existsC = auxC[VarMappingImageExtraTagsKey]
		if existsC {
			auxV[VarMappingImageExtraTagsKey] = auxC[VarMappingImageExtraTagsKey]
		}
	}

	_, existsV = auxV[VarMappingImageLabelsKey]
	if !existsV {
		_, existsC = auxC[VarMappingImageLabelsKey]
		if existsC {
			auxV[VarMappingImageLabelsKey] = auxC[VarMappingImageLabelsKey]
		}
	}

	v.SetUnderlyingMap(auxV)

	return nil
}
