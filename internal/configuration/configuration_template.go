package configuration

//

const (
	configurationTemplate string = `#
# Stevedore manages and governs the Docker's image's building process

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
#    log_path: /dev/null
{{ with .LogPathFile -}}
log_path: {{ . }}
{{ else -}}
#
# log_path: /dev/null
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
# On build, push images automatically after it finishes
#  default value: 
#    push_images: false
{{ if not .PushImages -}}
push_images: {{ .PushImages }}
{{ else -}}
#
# push_images: false
{{ end }}
#
# Directory to store docker registry credentials
#  default value: 
#    docker_registry_credentials_dir: ~/.config/stevedore/credentials
{{ with .DockerCredentialsDir -}}
docker_registry_credentials_dir: {{ . }}
{{ else -}}
#
# docker_registry_credentials_dir: ~/.config/stevedore/credentials
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
#   	Major      string
#   	Minor      string
#   	Patch      string
#   	PreRelease string
#   	Build      string
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
  - {{ . }}
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
#    1) Define an images tree that has one image named 'ubuntu' which has two versions, '20.04' and '18.04'.
#       The image '20.04' has two children: 'php-fpm' and 'php-cli', both having the version '7.4' defined. And the image '18.04' has only one child: 'php-fpm', with version '7.3' defined.
#       
#       images:
#         ubuntu:
#           18.04:
#             builder: infrastructure
#             children:
#             - php-fpm:
#               - 7.3
#           20.04:
#             builder: infrastructure
#             children:
#             - php-fpm:
#               - 7.4
#             - php-cli:
#               - 7.4
#         php-fpm:
#           7.3:
#             builder: infrastructure
#             builder:
#               driver: docker
#               context:
#                 path: php-fpm
#               vars:
#                 version: 7.3
#           7.4:
#             builder: infrastructure
#             builder:
#               driver: docker
#               context:
#                 path: php-fpm
#               vars:
#                 version: 7.4
#         php-cli:
#           7.4:
#             builder: infrastructure
#             builder:
#               driver: docker
#               context:
#                 path: php-fpm
#               vars:
#                 version: 7.4
`
)
