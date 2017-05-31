package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"reflect"
	"strings"
	"os"
)

type config_paths struct {
	flannel_conf_path        string
	master_conf_path         string
	node_conf_path           string
	validity_check_conf_path string
}

type DefaultOpts struct {
	Flannel_opts                 []string `yaml:"flannel_options"`
	Kubelet_opts                 []string `yaml:"kubelet_options"`
	Kube_proxy_opts              []string `yaml:"kube_proxy_options"`
	Kube_apiserver_opts          []string `yaml:"kube_apiserver_options"`
	Kube_controller_manager_opts []string `yaml:"kube_controller_manager_options"`
	Kube_scheduler_opts          []string `yaml:"kube_scheduler_options"`
}

type ValidityOpts struct {
	Args_check struct {
			   Flannel_opts                 []string `yaml:"flannel_options"`
			   Kubelet_opts                 []string `yaml:"kubelet_options"`
			   Kube_proxy_opts              []string `yaml:"kube_proxy_options"`
			   Kube_apiserver_opts          []string `yaml:"kube_apiserver_options"`
			   Kube_controller_manager_opts []string `yaml:"kube_controller_manager_options"`
			   Kube_scheduler_opts          []string `yaml:"kube_scheduler_options"`
		   } `yaml: args_check`
}

const relative_prefix_path = "../"

func main() {
	fmt.Println("++++ This tool check wheather default flags are same as in the role 'validity-check' in building process. ++++")
	fmt.Println("++++ Repository lies in https://github.com/xialonglee/compare-attr ++++")
	fmt.Println("++++ Put this tool in 'build' directory ++++")
	var conf_paths = config_paths{
		flannel_conf_path:        "roles/flannel/defaults/main.yaml",
		master_conf_path:         "roles/master/defaults/main.yml",
		node_conf_path:           "roles/node/defaults/main.yml",
		validity_check_conf_path: "roles/validity-check/defaults/main.yml",
	}
	mapper := map[string]string{
		"flannel": conf_paths.flannel_conf_path,
		"kube_apiserver" : conf_paths.master_conf_path,
		"kube_controller_manager":conf_paths.master_conf_path,
		"kube_scheduler" : conf_paths.master_conf_path,
		"kubelet": conf_paths.node_conf_path,
		"kube_proxy": conf_paths.node_conf_path,
	}

	for k, v := range mapper {
		paths := [2]string{conf_paths.validity_check_conf_path, v}
		isEqual, _ := CompareArgs(paths, ValidityOpts{}, DefaultOpts{}, k)
		fmt.Printf("%s flags keep consistency? %t\n", k, isEqual)
		if !isEqual {
			fmt.Printf("++++ Please check %s args in \"%s\" and \"%s\"! ++++\n", k, relative_prefix_path + conf_paths.validity_check_conf_path, relative_prefix_path + v)
			os.Exit(1)
		}
	}
}

func CompareArgs(file_path [2]string, validity_opts ValidityOpts, default_opts DefaultOpts, component_name string) (bool bool, err error) {
	fmt.Printf("\n++++ Check %s args in \"%s\" and \"%s\". ++++\n", component_name, file_path[0], file_path[1])
	field_name := strings.Title(component_name + "_opts")
	validity_value, default_value, err := yaml_unmarshal(file_path, validity_opts, default_opts, field_name)
	if err != nil {
		return false, err
	}

	fmt.Printf("%s validity args are: %v\n", component_name, validity_value)
	fmt.Printf("%s default args are: %v\n", component_name, default_value)
	if isInclude(validity_value, default_value) && isInclude(default_value, validity_value) {
		return true, nil
	}
	return false, nil
}

func yaml_unmarshal(file_path [2]string, validity_opts ValidityOpts, default_opts DefaultOpts, field_name string) ([]string, []string, error) {
	data, err := ioutil.ReadFile(relative_prefix_path + file_path[0])
	if err != nil {
		return []string{}, []string{}, fmt.Errorf("Reading config file failed: %v", err)
	}
	yaml.Unmarshal(data, &validity_opts)

	data, err = ioutil.ReadFile(relative_prefix_path + file_path[1])
	if err != nil {
		return []string{}, []string{}, fmt.Errorf("Reading config file failed: %v", err)
	}
	yaml.Unmarshal(data, &default_opts)

	validity_value := reflect.ValueOf(validity_opts).FieldByName("Args_check").FieldByName(field_name).Interface()
	default_value := reflect.ValueOf(default_opts).FieldByName(field_name).Interface()
	val, _ := validity_value.([]string)
	def, _ := default_value.([]string)
	return val, def, nil
}

func isInclude(args0, args1 []string) bool {
	regx := regexp.MustCompile(`s*-+([^s=]+).*`)
	for _, arg0 := range args0 {
		include := false
		for _, arg1 := range args1 {
			if regx.FindStringSubmatch(arg0)[1] == regx.FindStringSubmatch(arg1)[1] {
				include = true
				break
			}
		}
		if !include {
			return include
		}
	}
	return true
}
