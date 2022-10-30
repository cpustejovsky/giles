package giles

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
)

func (w *watcher) read(filepath string) error {
	match, err := regexp.MatchString(".*\\.yaml$", filepath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if !match {
		return fmt.Errorf("config file was not yaml, was insead:\t %v", filepath)
	}

	file, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	data := make(map[string]any)
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		return err
	}
	for k, v := range data {
		if k == "paths" {
			paths, ok := v.([]any)
			if !ok {
				fmt.Println(paths)
				return errors.New("incorrectly formatted paths key and value")
			}
			for _, path := range paths {
				fmt.Println("PATH", path)
				err = w.addPath(path.(string))
				if err != nil {
					return err
				}
			}
		}
		if k == "services" {
			services, ok := v.([]any)
			if !ok {
				return errors.New("incorrectly formatted services key and value")
			}
			for _, s := range services {
				service, ok := s.(map[string]any)
				if !ok {
					return errors.New("incorrectly formatted services key and value")
				}
				name := service["name"].(string)
				path := service["path"].(string)
				w.services = append(w.services, Service{
					Name: name,
					Path: path,
				})
			}
		}
	}
	return nil
}
