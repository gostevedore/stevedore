package helpers

type Commander interface {
	Execute() (string, error)
}

type AssertAndCommander interface {
	Commander
	AssertExectedResult(expected, result string)
}

type CommandFactorier interface {
	Command(command string) *DockerComposeTerratestCommand
}
