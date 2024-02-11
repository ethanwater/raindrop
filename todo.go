package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/urfave/cli/v2"
)

var (
	tasks []string
)

const todoFile string = "/Users/ethan/.todo"

func fetchTasks() error {
	file, err := os.Open(todoFile)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make([]byte, 1000)
	_, err = file.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	filter := func(arr []string, fn func(string) bool) []string {
		var result []string
		for _, v := range arr {
			switch {
			case len(v) < 1:
				continue
			default:
				result = append(result, v)
			}

		}

		return result[:len(result)-1]
	}

	tasks = filter(strings.Split(string(data), "\n"), func(task string) bool {
		return strings.TrimSpace(task) != ""
	})
	return nil
}

func updateTodo() error {
	var newContent strings.Builder
	if len(tasks) < 1 {
		fmt.Fprintf(&newContent, "\n")
	} else {
		for _, task := range tasks {
			fmt.Fprintf(&newContent, "%s\n", task)
		}
	}
	
	err := os.WriteFile(todoFile, []byte(newContent.String()), 0644)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func addTask(task string) error {
	fetchTasks()

	var ifTaskExists func([]string, string) bool
	ifTaskExists = func(arr []string, value string) bool {
		for _, v := range arr {
			if v == value {
				return true
			}
		}
		return false
	}

	if !ifTaskExists(tasks, task) {
		tasks = append(tasks, task)
	} else {
		fmt.Printf("notice: %s already exists\n", task)
	}

	updateTodo()
	fmt.Printf("[todo] added '%s' to tasks as [%d]\n", task, len(tasks))
	return nil
}

func removeTask(taskID int) {
	fetchTasks()
	switch {
	case len(tasks) == 1:
		fmt.Println("notice: to remove all elements use clear")
		return
	case taskID > len(tasks):
		fmt.Println("notice: taskID invalid")
		return
	}

	removedTask := tasks[taskID-1]
	tasks = append(tasks[:taskID-1], tasks[taskID:]...)

	updateTodo()
	fmt.Printf("[todo] removed '%s' from tasks\n", removedTask)
	return
}

func editTask(taskID int, newTask string) {
	fetchTasks()
	if taskID > len(tasks) {
		fmt.Println("notice: taskID invalid")
		return
	}
	originalTask := tasks[taskID-1]
	tasks[taskID-1] = newTask

	updateTodo()
	fmt.Printf("[todo] edited '%s' -> '%s'\n", originalTask, newTask)
	return
}

func doneTask(taskID int) {
	fetchTasks()
	if taskID > len(tasks) {
		fmt.Println("notice: taskID invalid")
		return
	}
	originalTask := tasks[taskID-1]
	completeTask := originalTask + "+"
	tasks[taskID-1] = completeTask
	fmt.Printf("[todo] '%s' done\n", originalTask)
	updateTodo()
	return
}

func clearTodo(sec string) error {
	fetchTasks()
	var filteredTasks []string
	switch sec {
	case "all":
		err := os.WriteFile(todoFile, []byte("\n"), 0644)
		if err != nil {
			return err
		}
		break
	case "done":
		for _, task := range tasks {
			if !strings.HasSuffix(task, "+") {
				filteredTasks = append(filteredTasks, task)
			}
		}
		break
	case "urgent":
		for _, task := range tasks {
			if !strings.HasSuffix(task, "!") {
				filteredTasks = append(filteredTasks, task)
			}
		}
		break
	case "misc":
		for _, task := range tasks {
			if strings.HasSuffix(task, "!") || strings.HasSuffix(task, "+") {
				filteredTasks = append(filteredTasks, task)
			}
		}
		break
	default:
		break
	}
	tasks = filteredTasks
	updateTodo()
	return nil
}

func displayTasks() {
	fetchTasks()
	if len(tasks) <= 0 {
		fmt.Printf("Nothing in Todo\n")
		return
	}

	urge	:= make(map[int]string) 
	misc	:= make(map[int]string)	
	done	:= make(map[int]string)

		for index, task := range tasks {
			if strings.HasSuffix(task, "!") {
				urge[index+1] = task[:len(task)-1]
				continue
			}
			if !strings.HasSuffix(task, "!") && !strings.HasSuffix(task, "+") {
				misc[index+1] = task
				continue
			}
			if strings.HasSuffix(task, "+") {
				done[index+1] = task[:len(task)-1]
				continue
			}
		}


	
	if len(urge) > 0 {
		fmt.Println(color.Ize(color.Red, "URGENT"))
		for id, task := range urge {
			fmt.Printf("[%d]: %s\n", id, task)
		}
	}

	if len(misc) > 0 {
		fmt.Println(color.Ize(color.Blue, "MISC:"))
		for id, task := range misc {
			fmt.Printf("[%d]: %s\n", id, task)
		}
	}

	if len(done) > 0 {
		fmt.Println(color.Ize(color.Green, "DONE:"))
		for id, task := range done {
			s := fmt.Sprintf("[%d]: %s\n", id, task)
			fmt.Print(color.Ize(color.Green, s))
		}
	}
}

func main() {
	app := &cli.App{
		Name:  "todo",
		Usage: "things to do",
		Action: func(*cli.Context) error {
			displayTasks()
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"+"},
				Usage:   "add a task",
				Action: func(cCtx *cli.Context) error {
					input := cCtx.Args().Get(0)
					addTask(input)
					return nil
				},
			},
			{
				Name:    "rm",
				Aliases: []string{"-", "remove"},
				Usage:   "remove a task",
				Action: func(cCtx *cli.Context) error {
					input, err := strconv.Atoi(cCtx.Args().Get(0))
					if err != nil {
						return err
					}
					removeTask(input)

					return nil
				},
			},
			{
				Name:  "done",
				Usage: "complete a task",
				Action: func(cCtx *cli.Context) error {
					input, err := strconv.Atoi(cCtx.Args().Get(0))
					if err != nil {
						return err
					}
					doneTask(input)

					return nil
				},
			},
			{
				Name:  "edit",
				Usage: "edit task",
				Action: func(cCtx *cli.Context) error {
					id, err := strconv.Atoi(cCtx.Args().Get(0))
					if err != nil {
						return err
					}
					newTask := cCtx.Args().Get(1)
					editTask(id, newTask)

					return nil
				},
			},
			{
				Name:  "clear",
				Usage: "clear todo list",
				Action: func(cCtx *cli.Context) error {
					sec := cCtx.Args().Get(0)
					clearTodo(sec)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
