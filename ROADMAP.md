# ROADMAP

## v0.10.1
- [] stevedore init creates configuration file with execution permissions
- [] promote does not uses -S flag to generate semantic version tags
- [] tags defined on the image are ignored

## v0.10.0
- [x] create an init command
- [x] automate semver tagging for images
- [x] inline builder definition on image definition
- [x] config precedence: current folder, user home, global folder
- [x] console message when no images to build found
- [x] define a depth to on cascade builds
- [x] on buildWorker notify user that semver won't be created
- [x] fix set image from details to build images
- [x] accept nil builder
- [x] update image childs to children

## v0.9.1
- [x] log_path is always writing on stdout when it is not overwrite.
- [x] test promte
- [x] when a build fails return an status code != 0 --> v0.9.1
- [x] functional tests using os.Exec --> v0.9.1 (tests moved to command test)
- [x] override the image name gave by argument with a flag --> v0.9.1
- [x] improve get_test.go:118: Functional test to get builders --> v0.9.1 (tests moved to command test)
- [x] imporve test builder_test.go:35: Testing array generation from a builder conf. Not always return the array elements in same order --> v0.9.1 (tests moved to command test)
- [x] test command package

## v0.9.0
- [x] define promote command
- [x] fix application panics when vars or pesistent_vars are not strings on image_tree definition
- [x] fix create complete path when creates a new credentials

## v0.8.1
- [x] fix remove wildcarded version images when all images are listed  
- [x] fix builders without options panics with a nil pointer exception
- [x] prefixing docker builder output
- [x] builders definition may accept remote repositories with definitions (git context)
- [x] map build's command flags builder variables
- [x] fix get builders to show non string options

## v0.8.0
- [x] new builders which uses docker as a driver
- [x] fix load persistent var from root node to leaf

## v0.7.1
- [x] accept cascade for wildcard version
- [x] Include tags and childs on copy image method

## v0.7.0
- [x] use copy Image method
- [x] test for wildcard version. even find methods
- [x] manage interruption
- [x] wildcard version
- [x] persistent vars
- [x] use context with interruptions --> v0.7.0

### issues
n/a

## scheduled
- [ ] inject dockerfiles --> v0.11.0
- [ ] define parents on image definitions --> v0.11.0
- [ ] remove image CheckCompatibility --> v0.11.0
- [ ] documentation

## icebox
- [ ] enable builds over http server
- [ ] copy image from one regitry to another one (it could be a promote flag)
- [ ] example on usage message
- [ ] apply levenshtein distances on image name to identify which image the user want to build 
- [ ] rewrite tree package
  - engine::imagesEngine: findNodes has dependencies
  - tree::renderizeGraphRec include to wildcard index nodes
  - tree::imageIndex global
