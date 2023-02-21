package promote

import (
	"context"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/application/promote"
)

type Handler struct {
	app PromoteApplication
}

func NewHandler(a PromoteApplication) *Handler {
	return &Handler{
		app: a,
	}
}

func (h *Handler) Handler(ctx context.Context, options *Options) error {

	errContext := "(handler::promote::Handler)"

	if options.SourceImageName == "" {
		return errors.New(errContext, "Source images name must be provided")
	}

	applicationOptions := &promote.Options{
		SourceImageName:       options.SourceImageName,
		RemoveTargetImageTags: options.RemoveTargetImageTags,
	}

	applicationOptions.DryRun = options.DryRun
	applicationOptions.EnableSemanticVersionTags = options.EnableSemanticVersionTags
	applicationOptions.PromoteSourceImageTag = options.PromoteSourceImageTag
	applicationOptions.RemoteSourceImage = options.RemoteSourceImage
	applicationOptions.RemoveTargetImageTags = options.RemoveTargetImageTags
	applicationOptions.SemanticVersionTagsTemplates = options.SemanticVersionTagsTemplates
	applicationOptions.TargetImageName = options.TargetImageName
	applicationOptions.TargetImageRegistryHost = options.TargetImageRegistryHost
	applicationOptions.TargetImageRegistryNamespace = options.TargetImageRegistryNamespace

	if len(options.TargetImageTags) > 0 {
		applicationOptions.TargetImageTags = append([]string{}, options.TargetImageTags...)
	}

	err := h.app.Promote(ctx, applicationOptions)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
