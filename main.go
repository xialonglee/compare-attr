package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"reflect"
	"strings"
)

type config_paths struct {
	flannel_conf_path        string
	master_conf_path         string
	node_conf_path           string
	validity_check_conf_path string
}

type DefaultOpts struct {
	Flannel_opts         []string `yaml:"flannel_options"`
	Kubelet_opts         []string `yaml:"kubelet_options"`
	Kube_proxy_opts      []string `yaml:"kube_proxy_options"`
	Kube_apiserver_opts  []string `yaml:"kube_apiserver_options"`
	Kube_controller_opts []string `yaml:"kube_controller_manager_options"`
	Kube_scheduller_opts []string `yaml:"kube_scheduler_options"`
}

type ValidityOpts struct {
	Args_check struct {
			   Flannel_opts         []string `yaml:"flannel_options"`
			   Kubelet_opts         []string `yaml:"kubelet_options"`
			   Kube_proxy_opts      []string `yaml:"kube_proxy_options"`
			   Kube_apiserver_opts  []string `yaml:"kube_apiserver_options"`
			   Kube_controller_opts []string `yaml:"kube_controller_manager_options"`
			   Kube_scheduller_opts []string `yasml:"kube_scheduler_options"`
		   } `yaml: args_check`
}

//const relative_prefix_path = "/home/k8s_kdt/kdt/"
const relative_prefix_path = "D:\\go_src\\contrib\\ansible\\"

var Conf_paths = config_paths{
	flannel_conf_path:        "roles/flannel/defaults/main.yaml",
	master_conf_path:         "roles/master/defaults/main.yml",
	node_conf_path:           "roles/node/defaults/main.yml",
	validity_check_conf_path: "roles/validity-check/defaults/main.yml",
}

func main() {
	paths := [2]string{relative_prefix_path + Conf_paths.validity_check_conf_path, relative_prefix_path + Conf_paths.flannel_conf_path}
	fmt.Println(CompareArgs(paths, ValidityOpts{}, DefaultOpts{}, "flannel"))

}

func CompareArgs(file_path [2]string, validity_opts ValidityOpts, default_opts DefaultOpts, component_name string) (bool bool, err error) {
	fmt.Printf("Check %s args in path %s and path %s.\n", component_name, file_path[0], file_path[1])
	field_name := strings.Title(component_name + "_opts")
	validity_value, default_value, err := yaml_unmarshal(file_path, validity_opts, default_opts, field_name)
	if err != nil {
		return false, err
	}

	fmt.Printf("validity_value is %v\n", validity_value)
	fmt.Printf("default_value is %v\n", default_value)
	if isInclude(validity_value, default_value) && isInclude(default_value, validity_value) {
		return true, nil
	}
	return false, nil
}

func yaml_unmarshal(file_path [2]string, validity_opts ValidityOpts, default_opts DefaultOpts, field_name string) ([]string, []string, error) {
	data, err := ioutil.ReadFile(file_path[0])
	if err != nil {
		return []string{}, []string{}, fmt.Errorf("Reading config file failed: %v", err)
	}
	yaml.Unmarshal(data, &validity_opts)

	data, err = ioutil.ReadFile(file_path[1])
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
