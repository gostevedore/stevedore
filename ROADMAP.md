# ROADMAP

## v0.12.0
- [ ] Multi-platform builds
- [ ] Build an image as well as all its parent until the root image
- [ ] Execute the build plan to show the build intentions
- [ ] Cleanup compatibilities
- [ ] Create a variable mapping with the normalized parent image name

## v0.11.1
- [x] Fix on promote command, the flags --promote-image-namespace and --promote-image-registry are not marked as deprecated
- [x] Fix typos on promote command examples

## v0.11.0
- [x] Fix: promote images with subnamespaces fails example myregistry/mynamesapce/grafana/tempo:1.0.0
- [x] Fix: wildcard images could not use {{ .Name }} on images definition
- [x] Fix: wildcarded images do not inherit parent persistent_vars
- [ ] Fix: container builders must match to  [a-zA-Z0-9][a-zA-Z0-9_.-]*
- [x] Credentials from env-vars
- [x] Promote local and remote images
- [x] Testing promote package
- [x] Inject Dockerfiles
- [x] Git context should accept authentication
- [x] Generate images_tree from multiples files located in a folder
- [x] Define parents on image definitions
- [x] Clean: remove image CheckCompatibility
- [x] Rewrite tree package
- [x] Remove-promote-tag to remove-local-images-after-push

## v0.10.1
- [x] Stevedore init creates a configuration file with execution permissions
- [x] Promote does not use -S flag to generate semantic version tags
- [x] Tags defined on the image are ignored

## v0.10.0
- [x] Create an init command
- [x] Automate semver tagging for images
- [x] Inline builder definition on image definition
- [x] Config precedence: current folder, user home, global folder
- [x] Console message when no images to build are found
- [x] Define a depth to cascade builds
- [x] On buildWorker notify the user that semver won't be created
- [x] Fix set image from details to build images
- [x] Accept nil builder
- [x] Update image children to children

## v0.9.1
- [x] Log_path is always written on stdout when it is not overwritten.
- [x] Test promote
- [x] When a build fails return a status code != 0 --> v0.9.1
- [x] Functional tests using os.Exec --> v0.9.1 (tests moved to command test)
- [x] Override the image name given by argument with a flag --> v0.9.1
- [x] Improve get_test.go:118: Functional test to get builders --> v0.9.1 (tests moved to command test)
- [x] Improve test builder_test.go:35: Testing array generation from a builder conf. Not always return the array elements in the same order --> v0.9.1 (tests moved to command test)
- [x] Test command package

## v0.9.0
- [x] Define promote command
- [x] Fix application panics when vars or pesistent_vars are not strings on image_tree definition
- [x] Fix create a complete path when creating new credentials

## v0.8.1
- [x] Fix remove wildcarded version images when all images are listed  
- [x] Fix builders without options panics with a nil pointer exception
- [x] Prefixing docker builder output
- [x] Builders definition may accept remote repositories with definitions (git context)
- [x] Map build's command flags builder variables
- [x] Fix get builders to show non-string options

## v0.8.0
- [x] New builders which use docker as a driver
- [x] Fix load persistent var from root node to leaf

## v0.7.1
- [x] Accept cascade for wildcard version
- [x] Include tags and children on copy image method

## v0.7.0
- [x] Use the copy Image method
- [x] Test for the wildcard version. even find methods
- [x] Manage interruption
- [x] Wildcard version
- [x] Persistent vars
- [x] Use context with interruptions --> v0.7.0

## icebox
- [ ] ~~Promote images defined on the images tree~~
- [ ] Signing Docker images
- [ ] Apply Levenshtein distances on the image name to identify which image the user wants to build
- [ ] Enable builds over HTTP
