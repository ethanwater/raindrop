package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

var tasks []string

const todoFile string = "/Users/ethan/.todo"

func fetchTasks() error {
	file, err := os.Open(todoFile)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make([]byte, 100)
	_, err = file.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	filter := func(arr []string, fn func(string) bool) []string {
		var result []string
		for _, v := range arr {
			if fn(v) {
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
	for _, task := range tasks {
		fmt.Fprintf(&newContent, "%s\n", task)
	}
	err := os.WriteFile(todoFile, []byte(newContent.String()), 0644)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func addTask(task string) error {
	fetchTasks()
	contains := func(arr []string, value string) bool {
		for _, v := range arr {
			if v == value {
				return true
			}
		}
		return false
	}

	if !contains(tasks, task) {
		tasks = append(tasks, task)
	}

	updateTodo()
	return nil
}

func removeTask(taskID int) {
	fetchTasks()
	if len(tasks) == 1 {
		fmt.Println("notice: to remove last element use clear")
		return
	}
	if taskID > len(tasks) {
		fmt.Println("notice: taskID invalid")
		return
	}
	tasks = append(tasks[:taskID-1], tasks[taskID:]...)
	updateTodo()
}

func clearTodo() error {
	err := os.WriteFile(todoFile, []byte("\n"), 0644)
	if err != nil {
		return err
	}

	return nil
}

func displayTasks() {
	fetchTasks()
	fmt.Println("todo:")
	for index, task := range tasks {
		fmt.Printf("[%d]: %s\n", index+1, task)
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
					addTask(cCtx.Args().Get(0))
					return nil
				},
			},
			{
				Name:    "rm",
				Aliases: []string{"-"},
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
				Name:  "clear",
				Usage: "clear todo list",
				Action: func(*cli.Context) error {
					clearTodo()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
