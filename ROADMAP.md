# ROADMAP

## v0.12.0
- [ ] apply Levenshtein distances on the image name to identify which image the user wants to build
- [ ] enable builds over HTTP
  
## v0.11.0
- [x] fix: promote images with subnamespaces fails example myregistry/mynamesapce/grafana/tempo:1.0.0
- [ ] fix: wildcard images could not use {{ .Name }} on images definition
- [ ] fix: wildcarded images do not inherit parent persistent_vars
- [ ] fix: container builders must match to  [a-zA-Z0-9][a-zA-Z0-9_.-]*
- [x] credentials from env-vars
- Redefine github.com/gostevedore/stevedore/internal/promote use cases
     [ ] - promote images from tree
     [x] - promote local and remote images
     [x] - to support test
- [x] inject Dockerfiles --> v0.11.0
- [x] git context should accept authentication
- [x] generate images_tree from multiples files located in a folder
- [x] define parents on image definitions
- [x] clean: remove image CheckCompatibility --> v0.11.0
- [x] rewrite tree package
  - engine::imagesEngine: findNodes has dependencies
  - tree::renderizeGraphRec include to wildcard index nodes
  - tree::imageIndex global
- [x] remove-promote-tag to remove-local-images-after-push

## v0.10.1
- [x] stevedore init creates a configuration file with execution permissions
- [x] promote does not use -S flag to generate semantic version tags
- [x] tags defined on the image are ignored

## v0.10.0
- [x] create an init command
- [x] automate semver tagging for images
- [x] inline builder definition on image definition
- [x] config precedence: current folder, user home, global folder
- [x] console message when no images to build are found
- [x] define a depth to cascade builds
- [x] on buildWorker notify the user that semver won't be created
- [x] fix set image from details to build images
- [x] accept nil builder
- [x] update image children to children

## v0.9.1
- [x] log_path is always written on stdout when it is not overwritten.
- [x] test promote
- [x] when a build fails return a status code != 0 --> v0.9.1
- [x] functional tests using os.Exec --> v0.9.1 (tests moved to command test)
- [x] override the image name given by argument with a flag --> v0.9.1
- [x] improve get_test.go:118: Functional test to get builders --> v0.9.1 (tests moved to command test)
- [x] improve test builder_test.go:35: Testing array generation from a builder conf. Not always return the array elements in the same order --> v0.9.1 (tests moved to command test)
- [x] test command package

## v0.9.0
- [x] define promote command
- [x] fix application panics when vars or pesistent_vars are not strings on image_tree definition
- [x] fix create a complete path when creating new credentials

## v0.8.1
- [x] fix remove wildcarded version images when all images are listed  
- [x] fix builders without options panics with a nil pointer exception
- [x] prefixing docker builder output
- [x] builders definition may accept remote repositories with definitions (git context)
- [x] map build's command flags builder variables
- [x] fix get builders to show non-string options

## v0.8.0
- [x] new builders which use docker as a driver
- [x] fix load persistent var from root node to leaf

## v0.7.1
- [x] accept cascade for wildcard version
- [x] Include tags and children on copy image method

## v0.7.0
- [x] use copy Image method
- [x] test for the wildcard version. even find methods
- [x] manage interruption
- [x] wildcard version
- [x] persistent vars
- [x] use context with interruptions --> v0.7.0

### issues
n/a

## icebox


