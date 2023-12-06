package file

const (
	configurationTemplate string = `#
# Stevedore is a Docker images factory, a tool that helps you to manage bunches of Docker image builds in just one command. It is not an alternative to Dockerfile or Buildkit, but a way to improve your building and promote experience

#
# Images tree location path
#  default value:
#    images_path: stevedore.yaml
{{ with .ImagesPath  -}}
images_path: {{ . }}
{{ else -}}
#
# images_path: stevedore.yaml
{{ end }}
#
# Builders location path
#  default value:
#    builders_path: stevedore.yaml
{{ with .BuildersPath -}}
builders_path: {{ . }}
{{ else -}}
#
# builders_path: stevedore.yaml
{{ end }}
#
# Log file location path
#  default value: 
#    log_path: /var/log/stevedore.log
{{ with .LogPathFile -}}
log_path: {{ . }}
{{ else -}}
#
# log_path: /var/log/stevedore.log
{{ end }}
#
# It defines the number of workers to build images which corresponds to the number of images that can be build concurrently
#  default value: 
#    concurrency: 4
{{ with .Concurrency -}}
concurrency: {{ . }}
{{ else -}}
#
# concurrency: 4
{{ end }}
#
# Push images automatically after build
#  default value: 
#    push_images: false
{{ if .PushImages -}}
push_images: {{ .PushImages }}
{{ else -}}
#
# push_images: false
{{ end }}
#
# Credentials storage
#   default value:
#     credentials:
#       storage_type: local
#       local_storage_path: /var/lib/stevedore/credentials
#       format: json
# 
{{ with .Credentials -}}
credentials:
  storage_type: {{ .StorageType }}
  format: {{ .Format }}
  {{ if eq .StorageType "local" -}}
  local_storage_path: {{ .LocalStoragePath }}
  {{ end -}}
  {{ if ne .EncryptionKey "" -}}
  encryption_key: {{ .EncryptionKey }}
  {{ end -}}
{{ else }}
# credentials:
#   storage_type: local
#	  local_storage_path: /var/lib/stevedore/credentials
#	  format: json
{{ end }}
#
# Generate extra tags when the main image tags is semver 2.0.0 compliance
#  default value: false
#    semantic_version_tags_enabled: false
{{ with .EnableSemanticVersionTags -}}
semantic_version_tags_enabled: {{ . }}
{{ else -}}
#
# semantic_version_tags_enabled: false
{{ end }}
# List of templates which define those extra tags to generate when 'semantic_version_tags_enabled' is enabled
# Parser will use the SemVer struct to generate the template, and all SemVer attributes could be used to define each template
# 
#   // SemVer is a sematinc version representation
#   type SemVer struct {
#       Major      string
#       Minor      string
#       Patch      string
#       PreRelease string
#       Build      string
#   }
#
#  default value: 
#    semantic_version_tags_templates:
#      - #u007b#u007b .Major #u007d#u007d.#u007b#u007b .Minor #u007d#u007d.#u007b#u007b .Patch #u007d#u007d
{{ if not .SemanticVersionTagsTemplates -}}
# semantic_version_tags_templates:
#   - #u007b#u007b .Major #u007d#u007d.#u007b#u007b .Minor #u007d#u007d.#u007b#u007b .Patch #u007d#u007d
{{ else -}}
semantic_version_tags_templates:
{{ range .SemanticVersionTagsTemplates -}}
  - "{{ . }}"
{{ end -}}
{{ end }}
#
# Define builder types
# You could define builders on its own file. Stevedore will look up for builders on the file set at 'builders_path'
# 
# examples:
#   1) Define one builder named 'infrastructure' which use docker as driver and current folder as Docker build context
#
#      builders:
#       infrastructure:
#         driver: docker
#         options:
#           context:
#              path: .

#
# Define images tree
# You could define an images tree on its own file. Stevedore will look up for images tree on the file set at 'images_path'
# 
#  examples:
#    1) Define an images tree that has one image named 'ubuntu' which has two versions, '20.04' and '22.04'.
#       The image '20.04' has two children: 'php-fpm' and 'php-cli', both having the version '8.0' defined. And the image '22.04' has only one child: 'php-fpm', with version '8.1' defined.
#       
#       images:
#         ubuntu:
#           22.04:
#             namespace: infrastructure
#             builder: infrastructure
#             children:
#             - php-fpm:
#               - 8.1
#           20.04:
#             namespace: infrastructure
#             builder: infrastructure
#         php-fpm:
#           8.1:
#             builder:
#               driver: docker
#               context:
#                 path: php-fpm
#             vars:
#               version: 8.1
#           8.0:
#             builder:
#               driver: docker
#               context:
#                 path: php-fpm
#               vars:
#                 version: 8.0
#             parents:
#               ubuntu:
#                 - 20.04
#                 - 22.04
#         php-cli:
#           8.0:
#             builder:
#               driver: docker
#               context:
#                 path: php-fpm
#             vars:
#               version: 8.0
#             parents:
#               ubuntu:
#                 - 20.04
#                 - 22.04
`
)
