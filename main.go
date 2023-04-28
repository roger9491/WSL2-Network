package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {
	// Execute "wsl hostname -I" to get the IP address of WSL2.
	cmd := exec.Command("wsl", "hostname", "-I")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		return
	}

	// Parse the output and extract the IP address.
	outputStr := strings.TrimSpace(string(output))
	ipRegEx := regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)
	ip := ipRegEx.FindString(outputStr)

	fmt.Printf("WSL2 IP address: %s\n", ip)

	// 清除所有的端口轉發
	cmd = exec.Command("netsh", "interface", "portproxy", "reset")

	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Command failed with error:", err)
	}
	fmt.Println(string(output))

	// 獲取要映射的port號
	// 讀取 YAML 檔案
	data, err := ioutil.ReadFile(".yaml")
	if err != nil {
		panic(err)
	}

	// 解析 YAML 檔案
	var config map[string]interface{}
	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		panic(err)
	}

	listenAddress := "0.0.0.0"

	// 使用 config 變數
	for k, v := range config {
		fmt.Println(k, v)
		v = strconv.Itoa(v.(int))
		// Command to execute.
		cmd := exec.Command("netsh", "interface", "portproxy", "add",
			"v4tov4", "listenaddress="+listenAddress, "listenport="+v.(string),
			"connectaddress="+ip, "connectport="+v.(string))
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Command failed with error:", err)
		}
		fmt.Println(string(output))
	}

	cmd = exec.Command("netsh", "interface", "portproxy", "show", "v4tov4")

	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Command failed with error:", err)
		return
	}
	fmt.Println(string(output))
}
