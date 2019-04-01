package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

func mustYamlFileToGoObject(path string, obj interface{}) {
	contents := mustRead(path)
	if err := yaml.Unmarshal(contents, obj); err != nil {
		panic(err.Error())
	}
}

func mustRead(path string) []byte {
	read, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}
	return read
}

func mustWrite(path, out string) {
	if err := ioutil.WriteFile(path, []byte(out), 0644); err != nil {
		panic(err.Error())
	}
	exec.Command("gofmt", "-w", path).Output()
}

func mustExecuteTemplateFile(templateFilepath string, data interface{}, funcMaps template.FuncMap, delims ...[]string) string {
	ret, err := executeTemplateSource(string(mustRead(templateFilepath)), data, funcMaps, delims...)
	if err != nil {
		log.Fatal(err.Error())
	}
	return ret
}

func executeTemplateSource(templateString string, data interface{}, funcMaps template.FuncMap, delims ...[]string) (string, error) {
	t := template.New(reflect.TypeOf(data).Name()).Funcs(funcMaps)
	if len(delims) > 0 {
		t.Delims(delims[0][0], delims[0][1]).Parse(templateString)
	}

	t, err := t.Parse(templateString)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func executer(cmdType, cmdName string, cmdArgs []string, cmdWait bool, executeDir string, addEnv []string) (string, error) {
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, addEnv...)
	if executeDir != "" {
		cmd.Dir = executeDir
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
