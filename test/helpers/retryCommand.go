package helpers

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

const (
	// defautlRetry is the default number of times to retry a command.
	defaultRetry = 3
	// defaultWaitTimeTillNextRetrySeconds is the default number of seconds to sleep between retries.
	defaultWaitTimeTillNextRetrySeconds = 2
)

type RetryCommand struct {
	commander                    Commander
	postRetryCommand             []Commander
	retry                        int
	waitTimeTillNextRetrySeconds int
}

func NewRetryCommand(commander Commander) *RetryCommand {
	return &RetryCommand{
		commander:        commander,
		postRetryCommand: make([]Commander, 0, 10),
	}
}

func (r *RetryCommand) WithPostRetryCommand(cmd ...Commander) *RetryCommand {
	r.postRetryCommand = append(r.postRetryCommand, cmd...)
	return r
}

func (r *RetryCommand) WithRetry(retry int) *RetryCommand {
	r.retry = retry
	return r
}

func (r *RetryCommand) WithWaitTimeTillNextRetrySeconds(waitTimeTillNextRetrySeconds int) *RetryCommand {
	r.waitTimeTillNextRetrySeconds = waitTimeTillNextRetrySeconds
	return r
}

func (r *RetryCommand) Execute() (string, error) {
	var result string
	var err error
	var retry, waitTimeTillNextRetrySeconds int

	retry = r.retry
	waitTimeTillNextRetrySeconds = r.waitTimeTillNextRetrySeconds

	if retry == 0 {
		retry = defaultRetry
	}

	if waitTimeTillNextRetrySeconds == 0 {
		waitTimeTillNextRetrySeconds = defaultWaitTimeTillNextRetrySeconds
	}

	for i := 0; i < retry; i++ {
		result, err = r.commander.Execute()
		if err == nil {
			return result, nil
		}
		fmt.Println("Retry", i+1, "of", retry, "in", waitTimeTillNextRetrySeconds, "seconds")
		time.Sleep(time.Duration(waitTimeTillNextRetrySeconds) * time.Second)

		color.Green(" Post-retry actions")
		for _, command := range r.postRetryCommand {
			_, err := command.Execute()
			if err != nil {
				return "", err
			}
		}
		fmt.Println()
	}

	return result, err
}
