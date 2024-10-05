package reader

import (
	"fmt"
	"log"
	"os"

	"io/fs"
	"path/filepath"
	"strings"

	"task-lineage-diagram/schema"

	"gopkg.in/yaml.v3"
)

var taskPaths []string

func getTask(filepath string) schema.Task {
	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	var c schema.Task
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("%s: Unmarshal: %v", filepath, err)
	}

	return c
}

func walk(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !d.IsDir() && (strings.HasSuffix(s, ".yml") || strings.HasSuffix(s, ".yaml")) {
		pathStrs := strings.Split(s, "/")
		fileStr := pathStrs[len(pathStrs)-1]
		taskLayer := pathStrs[len(pathStrs)-2]
		if !(strings.HasPrefix(fileStr, "_") || strings.HasPrefix(fileStr, ".") || strings.HasPrefix(taskLayer, "_")) {
			fmt.Printf("found yaml file %s\n", s)
			taskPaths = append(taskPaths, s)
		}
	}
	return nil
}

func ReadTasks(dirRoot string) (map[string]schema.Task, error) {
	err := filepath.WalkDir(dirRoot, walk)
	tasks := make(map[string]schema.Task)

	for _, path := range taskPaths {
		task := getTask(path)
		if task.TaskID == "" {
			fmt.Printf("WARNING: Cannot get process name for %s. Skip this file.\n", path)
		} else {
			tasks[task.TaskID] = task
		}
	}

	return tasks, err
}

func ReadConfig(filepath string) (*schema.Config, error) {
	yamlFile, err := os.ReadFile(filepath)
	if os.IsNotExist(err) {
		log.Printf("Cannot find config file, using default settings   #%v ", err)
		// If file does not exist, return default config with white as default color
		return &schema.Config{
			Colors: map[string]string{
				"default": "#FFFFFF", // Default to white color
			},
		}, nil
	} else if err != nil {
		return nil, err
	}
	var c schema.Config
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("%s: Unmarshal: %v", filepath, err)
	}
	return &c, nil
}
