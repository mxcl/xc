package run

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joerdav/xc/models"
)

const maxDeps = 50

func runCmd(c *exec.Cmd) error {
	return c.Run()
}

// Runner is responsible for running Tasks.
type Runner struct {
	sep, cmdRunner, flag string
	runner               func(*exec.Cmd) error
	tasks                models.Tasks
}

// NewRunner takes Tasks and returns a Runner.
// If the OS is windows commands will be run using `cmd \C`
// and separated by `&&`.
// Otherwise, commands will be run using `bash -c`
// and separated by `;`.
//
// NewRunner will return an error in the case that Dependent tasks are cyclical,
// invalid or at a larger depth than 50.
func NewRunner(ts models.Tasks, runtime string) (runner Runner, err error) {
	runner = Runner{
		sep:       ";",
		cmdRunner: "bash",
		flag:      "-c",
		runner:    runCmd,
		tasks:     ts,
	}
	if runtime == "windows" {
		runner.sep = "&&"
		runner.cmdRunner = "cmd"
		runner.flag = "/C"
	}
	for _, t := range ts {
		err = runner.ValidateDependencies(t.Name, []string{})
		if err != nil {
			return
		}
	}
	return
}

// Run runs a task given a string name.
// Task dependencies will be run first, an error will return if any fail.
// Task commands are run next, in case of a non zero result an error will return.
func (r *Runner) Run(ctx context.Context, name string) error {
	task, ok := r.tasks.Get(name)
	if !ok {
		return fmt.Errorf("task %s not found", name)
	}
	for _, t := range task.DependsOn {
		err := r.Run(ctx, t)
		if err != nil {
			return err
		}
	}
	var cmdl []string
	for _, c := range task.Commands {
		if strings.TrimSpace(c) == "" {
			continue
		}
		cmdl = append(cmdl, fmt.Sprintf(`echo "%s"`, c), c)
	}
	if len(task.Commands) == 0 {
		return nil
	}
	cmds := strings.Join(cmdl, r.sep)
	cmd := exec.Command(r.cmdRunner, r.flag, cmds)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, task.Env...)
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	cmd.Dir = path
	if task.Dir != "" {
		cmd.Dir = task.Dir
	}
	err = r.runner(cmd)
	if err != nil {
		return err
	}
	return nil
}

// ValidateDependencies checks that task dependencies follow these rules:
// - No deeper dependency trees than maxDeps.
// - Dependencies must exist as tasks.
// - No cyclical dependencies.
func (r *Runner) ValidateDependencies(task string, prevTasks []string) error {
	if len(prevTasks) >= maxDeps {
		return fmt.Errorf("max dependency depth of %d reached", maxDeps)
	}
	// Check exists
	t, ok := r.tasks.Get(task)
	if !ok {
		return fmt.Errorf("task %s not found", task)
	}
	if t.ParsingError != "" {
		return fmt.Errorf("task %s has a parsing error: %s", task, t.ParsingError)
	}
	for _, t := range t.DependsOn {
		st, ok := r.tasks.Get(t)
		if !ok {
			return fmt.Errorf("task %s not found", t)
		}
		for _, pt := range prevTasks {
			if pt == st.Name {
				return fmt.Errorf("task %s contains a circular dependency", t)
			}
		}
		err := r.ValidateDependencies(st.Name, append([]string{st.Name}, prevTasks...))
		if err != nil {
			return err
		}
	}
	return nil
}
