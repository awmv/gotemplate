package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/kyokomi/emoji"
	"github.com/manifoldco/promptui"
)

// File contains properties that are needed to create a file
type File struct {
	Path    string
	Content string
}

func main() {
	// name, nameSpace, givenPath, errMsg := prompt()
	name, nameSpace, givenPath := prompt()
	if len(name) == 0 || len(nameSpace) == 0 {
		fmt.Println("Parameter Name and Namespace are required.")
		return
	}
	if len(givenPath) == 0 {
		dir, err := os.Getwd()
		exitOnError(err, "Failed to get working directory")
		givenPath = dir
	}
	projectPath := path.Join(givenPath, strings.ToLower(name))
	info, _ := os.Stat(projectPath)
	if info != nil {
		if info.IsDir() {
			fmt.Println("Project {" + name + "} already exists in " + projectPath + ".")
			return
		}
	}
	err := createDirectories(projectPath)
	exitOnError(err, "Failed to create directories")

	// cmd/service/main.go
	err = createFile(path.Join(projectPath, "cmd/service/main.go"))
	exitOnError(err, "Failed to create "+name+"cmd/service/main.go")
	const base64MainGo = "cGFja2FnZSBtYWluCgppbXBvcnQgImZtdCIKCmZ1bmMgbWFpbigpIHsKCWZtdC5QcmludGxuKCJIZWxsbyIpCn0K"
	mainGo, _ := decode(base64MainGo)

	err = writeFile(path.Join(projectPath, "cmd/service/main.go"), mainGo)
	exitOnError(err, "Failed to write in "+name+"cmd/service/main.go")

	createFilesInConfigs(name, nameSpace, projectPath)

	// docs/redoc.go
	err = createFile(path.Join(projectPath, "docs/redoc.go"))
	exitOnError(err, "Failed to create "+name+"docs/redoc.go")
	const base64Redoc = "cGFja2FnZSBkb2NzCgovLyBSZWRvY0RvY3VtZW50YXRpb24gaXMgYSBzdGF0aWMgSFRNTCB1c2VkIHRvIGRlbGl2ZXIgdGhlIHJlZG9jIEFQSS1Eb2N1bWVudGF0aW9uCmNvbnN0IFJlZG9jRG9jdW1lbnRhdGlvbiA9IGA8IURPQ1RZUEUgaHRtbD4KPGh0bWw+CiAgPGhlYWQ+CiAgICA8dGl0bGU+T25ib3JkaW5nIEFQSSBEb2N1bWVudGF0aW9uPC90aXRsZT4KICAgIDwhLS0gbmVlZGVkIGZvciBhZGFwdGl2ZSBkZXNpZ24gLS0+CiAgICA8bWV0YSBjaGFyc2V0PSJ1dGYtOCIvPgogICAgPG1ldGEgbmFtZT0idmlld3BvcnQiIGNvbnRlbnQ9IndpZHRoPWRldmljZS13aWR0aCwgaW5pdGlhbC1zY2FsZT0xIj4KICAgIDxsaW5rIGhyZWY9Imh0dHBzOi8vZm9udHMuZ29vZ2xlYXBpcy5jb20vY3NzP2ZhbWlseT1Nb250c2VycmF0OjMwMCw0MDAsNzAwfFJvYm90bzozMDAsNDAwLDcwMCIgcmVsPSJzdHlsZXNoZWV0Ij4KCiAgICA8IS0tCiAgICBSZURvYyBkb2Vzbid0IGNoYW5nZSBvdXRlciBwYWdlIHN0eWxlcwogICAgLS0+CiAgICA8c3R5bGU+CiAgICAgIGJvZHkgewogICAgICAgIG1hcmdpbjogMDsKICAgICAgICBwYWRkaW5nOiAwOwogICAgICB9CiAgICA8L3N0eWxlPgogIDwvaGVhZD4KICA8Ym9keT4KICAgIDxyZWRvYyBzcGVjLXVybD0naHR0cHM6Ly97bmFtZX0udGVzdC5maW5vLmNsb3VkL2FwaS92MS9kb2MvZG9jLmpzb24nPjwvcmVkb2M+CiAgICA8c2NyaXB0IHNyYz0iaHR0cHM6Ly9jZG4uanNkZWxpdnIubmV0L25wbS9yZWRvY0BuZXh0L2J1bmRsZXMvcmVkb2Muc3RhbmRhbG9uZS5qcyI+IDwvc2NyaXB0PgogIDwvYm9keT4KPC9odG1sPmAK"
	redoc, _ := decode(base64Redoc)
	redoc = replacePlaceholder(redoc, "{name}", nameSpace)

	err = writeFile(path.Join(projectPath, "docs/redoc.go"), redoc)
	exitOnError(err, "Failed to write in "+name+"docs/redoc.go")

	// docs/api/readme.md
	err = createFile(path.Join(projectPath, "docs/api/readme.md"))
	exitOnError(err, "Failed to create "+name+"docs/api/readme.md")

	err = writeFile(path.Join(projectPath, "docs/api/readme.md"), "Hallo Welt")
	exitOnError(err, "Failed to write in "+name+"docs/api/readme.md")

	// readme.md
	err = createFile(path.Join(projectPath, "readme.md"))
	exitOnError(err, "Failed to create "+name+"readme.md")

	err = writeFile(path.Join(projectPath, "readme.md"), "Hallo Welt")
	exitOnError(err, "Failed to write in "+name+"readme.md")

	createFilesInPath(name, nameSpace, projectPath)

	msg := emoji.Sprintf("Project {%s} has been created in %s :beer:", strings.ToLower(name), projectPath)
	fmt.Println(msg)
}

// createDirectories iterates over createDirectory
func createDirectories(projectPath string) error {
	directories := [5]string{
		"cmd/service", "configs", "docs/api", "pkg", "test",
	}
	for _, location := range directories {
		err := createDirectory(path.Join(projectPath, location))
		if err != nil {
			return err
		}
	}
	return nil
}

// createDirectory creates a directory
func createDirectory(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	return err
}

// createFile creates a file
func createFile(path string) error {
	var _, err = os.Stat(path)

	if os.IsNotExist(err) {
		file, err := os.Create(path)

		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

// writeFile writes in a file
func writeFile(path string, content string) error {
	var file, err = os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	err = file.Sync()
	if err != nil {
		return err
	}
	return nil
}

// decode decodes a string
func decode(str string) (string, error) {
	if len(str) == 0 {
		return str, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Println("Failed to decode: ", err)
		return "", err
	}
	return string(decoded), nil
}

// replacePlaceholder replaces a placeholders in a string
func replacePlaceholder(content string, placeholder string, replacement string) string {
	return strings.ReplaceAll(content, placeholder, replacement)
}

// createFilesInConfigs creates/writes in ./configs
func createFilesInConfigs(name string, nameSpace string, projectPath string) {
	// configs/deployment.yml
	err := createFile(path.Join(projectPath, "configs/deployment.yml"))
	exitOnError(err, "Failed to create "+name+"configs/deployment.yml")
	const base64Deployment = "YXBpVmVyc2lvbjogYXBwcy92MQpraW5kOiBEZXBsb3ltZW50Cm1ldGFkYXRhOgogIG5hbWU6IHtuYW1lfQogIG5hbWVzcGFjZToge25hbWVzcGFjZX0KICBsYWJlbHM6CiAgICBhcHA6IHtuYW1lc3BhY2V9CiAgICBjb21wb25lbnQ6IHtuYW1lfQpzcGVjOgogIHJlcGxpY2FzOiAxCiAgc2VsZWN0b3I6CiAgICBtYXRjaExhYmVsczoKICAgICAgYXBwOiB7bmFtZXNwYWNlfQogICAgICBjb21wb25lbnQ6IHtuYW1lfQogIHRlbXBsYXRlOgogICAgbWV0YWRhdGE6CiAgICAgIGxhYmVsczoKICAgICAgICBhcHA6IHtuYW1lc3BhY2V9CiAgICAgICAgY29tcG9uZW50OiB7bmFtZX0KICAgIHNwZWM6CiAgICAgIGNvbnRhaW5lcnM6CiAgICAgIC0gbmFtZToge25hbWV9CiAgICAgICAgZW52OgogICAgICAgICAgICB2YWx1ZTogJHtGSU5PX0JVSUxEX1RBR30KICAgICAgICAgIC0gbmFtZTogU0VSVklDRV9OQU1FCiAgICAgICAgICAgIHZhbHVlOiB7bmFtZXNwYWNlfS17bmFtZX0KICAgICAgICAgIC0gbmFtZTogRU5WSVJPTk1FTlQKICAgICAgICAgICAgdmFsdWU6ICR7RklOT19FTlZJUk9OTUVOVH0KICAgICAgICAgIC0gbmFtZTogRU5WSVJPTk1FTlRfUFJFRklYCiAgICAgICAgICAgIHZhbHVlOiAke0ZJTk9fRU5WX1BSRUZJWH0KICAgICAgICAgIC0gbmFtZTogTE9HR0lOR19VUkwKICAgICAgICAgICAgdmFsdWVGcm9tOgogICAgICAgICAgICAgIHNlY3JldEtleVJlZjoKICAgICAgICAgICAgICAgIG5hbWU6IGxvZ2dpbmcKICAgICAgICAgICAgICAgIGtleTogdXJsCiAgICAgICAgICAtIG5hbWU6IE1PTkdPREJfU0VSVkVSCiAgICAgICAgICAgIHZhbHVlOiBtZGIxLnNoYXJlZC1zZXJ2aWNlcywgbWRiMi5zaGFyZWQtc2VydmljZXMKICAgICAgICAgIC0gbmFtZTogTU9OR09EQl9SRVBMSUNBU0VUX05BTUUKICAgICAgICAgICAgdmFsdWU6ICdyczAnCiAgICAgICAgICAtIG5hbWU6IE1PTkdPREJfREFUQUJBU0UKICAgICAgICAgICAgdmFsdWU6IHtuYW1lc3BhY2V9CiAgICAgICAgICAtIG5hbWU6IE1PTkdPREJfVVNFUk5BTUUKICAgICAgICAgICAgdmFsdWVGcm9tOgogICAgICAgICAgICAgIHNlY3JldEtleVJlZjoKICAgICAgICAgICAgICAgIG5hbWU6IG1vbmdvZGIKICAgICAgICAgICAgICAgIGtleTogdXNlcgogICAgICAgICAgLSBuYW1lOiBNT05HT0RCX1BBU1NXT1JECiAgICAgICAgICAgIHZhbHVlRnJvbToKICAgICAgICAgICAgICBzZWNyZXRLZXlSZWY6CiAgICAgICAgICAgICAgICBuYW1lOiBtb25nb2RiCiAgICAgICAgICAgICAgICBrZXk6IHBhc3N3b3JkICAgICAgICAKICAgICAgICBpbWFnZTogZG9tYWluLmNvbS97bmFtZXNwYWNlfS97bmFtZX06JHtGSU5PX0JVSUxEX1RBR30KICAgICAgICBpbWFnZVB1bGxQb2xpY3k6IEFsd2F5cwogICAgICAgIHBvcnRzOgogICAgICAgICAgLSBjb250YWluZXJQb3J0OiA4MDgwCiAgICAgICAgbGl2ZW5lc3NQcm9iZToKICAgICAgICAgIGh0dHBHZXQ6CiAgICAgICAgICAgIHBhdGg6IC9saXZlCiAgICAgICAgICAgIHBvcnQ6IDgwODYKICAgICAgICAgIGluaXRpYWxEZWxheVNlY29uZHM6IDE1CiAgICAgICAgICBwZXJpb2RTZWNvbmRzOiAzCiAgICAgICAgICB0aW1lb3V0U2Vjb25kczogNQogICAgICAgIHJlYWRpbmVzc1Byb2JlOgogICAgICAgICAgaHR0cEdldDoKICAgICAgICAgICAgcGF0aDogL3JlYWR5CiAgICAgICAgICAgIHBvcnQ6IDgwODYKICAgICAgICAgIGluaXRpYWxEZWxheVNlY29uZHM6IDMKICAgICAgICAgIHBlcmlvZFNlY29uZHM6IDMKICAgICAgICAgIHRpbWVvdXRTZWNvbmRzOiA1CiAgICAgICAgcmVzb3VyY2VzOgogICAgICAgICAgbGltaXRzOgogICAgICAgICAgICBjcHU6IDQwbQogICAgICAgICAgICBtZW1vcnk6IDQwTWkKICAgICAgICAgIHJlcXVlc3RzOgogICAgICAgICAgICBjcHU6IDEwbQogICAgICAgICAgICBtZW1vcnk6IDEwTWkKICAgICAgaW1hZ2VQdWxsU2VjcmV0czoKICAgICAgICAtIG5hbWU6IGF6dXJlY3I="
	deployment, _ := decode(base64Deployment)
	deployment = replacePlaceholder(deployment, "{name}", name)
	deployment = replacePlaceholder(deployment, "{namespace}", nameSpace)

	err = writeFile(path.Join(projectPath, "configs/deployment.yml"), deployment)
	exitOnError(err, "Failed to write in "+name+"configs/deployment.yml")

	// configs/ingress.yml
	err = createFile(path.Join(projectPath, "configs/ingress.yml"))
	exitOnError(err, "Failed to create "+name+"configs/ingress.yml")
	const base64Ingress = "YXBpVmVyc2lvbjogZXh0ZW5zaW9ucy92MWJldGExCmtpbmQ6IEluZ3Jlc3MKbWV0YWRhdGE6CiAgbmFtZToge25hbWV9CiAgbmFtZXNwYWNlOiB7bmFtZXNwYWNlfQogIGFubm90YXRpb25zOgogICAga3ViZXJuZXRlcy5pby9pbmdyZXNzLmNsYXNzOiB0cmFlZmlrICAgIAogICAgdHJhZWZpay5pbmdyZXNzLmt1YmVybmV0ZXMuaW8vcmVkaXJlY3QtZW50cnktcG9pbnQ6IGh0dHBzCiAgICB0cmFlZmlrLmluZ3Jlc3Mua3ViZXJuZXRlcy5pby9yZWRpcmVjdC1wZXJtYW5lbnQ6ICJ0cnVlIgpzcGVjOgogIHJ1bGVzOgogIC0gaG9zdDoge25hbWVzcGFjZX0uJHtGSU5PX0VOVl9QUkVGSVh9Zmluby5jbG91ZAogICAgaHR0cDoKICAgICAgcGF0aHM6CiAgICAgIC0gcGF0aDogL2FwaQogICAgICAgIGJhY2tlbmQ6CiAgICAgICAgICBzZXJ2aWNlTmFtZToge25hbWV9CiAgICAgICAgICBzZXJ2aWNlUG9ydDogODAK"
	ingress, _ := decode(base64Deployment)
	ingress = replacePlaceholder(ingress, "{name}", name)
	ingress = replacePlaceholder(ingress, "{namespace}", nameSpace)

	err = writeFile(path.Join(projectPath, "configs/ingress.yml"), ingress)
	exitOnError(err, "Failed to write in "+name+"configs/ingress.yml")

	// configs/service.yml
	err = createFile(path.Join(projectPath, "configs/service.yml"))
	exitOnError(err, "Failed to create "+name+"configs/service.yml")
	const base64Service = "YXBpVmVyc2lvbjogdjEKa2luZDogU2VydmljZQptZXRhZGF0YToKICBuYW1lOiB7bmFtZX0KICBuYW1lc3BhY2U6IHtuYW1lc3BhY2V9CiAgbGFiZWxzOgogICAgYXBwOiB7bmFtZXNwYWNlfQogICAgY29tcG9uZW50OiB7bmFtZX0Kc3BlYzoKICBwb3J0czoKICAtIHBvcnQ6IDgwCiAgICB0YXJnZXRQb3J0OiA4MDgwCiAgICBwcm90b2NvbDogVENQCiAgc2VsZWN0b3I6CiAgICBhcHA6IHtuYW1lc3BhY2V9CiAgICBjb21wb25lbnQ6IHtuYW1lfQ=="
	service, _ := decode(base64Service)
	service = replacePlaceholder(service, "{name}", name)
	service = replacePlaceholder(service, "{namespace}", nameSpace)

	err = writeFile(path.Join(projectPath, "configs/service.yml"), service)
	exitOnError(err, "Failed to write in "+name+"configs/service.yml")
}

// exitOnError exists the the app on an error
func exitOnError(err error, msg string) {
	if err != nil {
		fmt.Println(msg+": ", err)
		os.Exit(1)
	}
}

// getParameters gets parameters off the prompt
func getParameters() (string, string, string) {
	prompt := promptui.Prompt{
		Label: "Name",
	}
	name, err := prompt.Run()
	exitOnError(err, "Prompt failed at name")
	prompt = promptui.Prompt{
		Label: "Namespace",
	}
	nameSpace, err := prompt.Run()
	exitOnError(err, "Prompt failed at namespace")
	dir, err := os.Getwd()
	exitOnError(err, "Failed to get working directory")
	prompt = promptui.Prompt{
		Label:   "Path",
		Default: dir,
	}
	givenPath, err := prompt.Run()
	exitOnError(err, "Prompt failed at path")
	return name, nameSpace, givenPath
}

// createFilesInPath creates/writes files in ./
func createFilesInPath(name string, nameSpace string, projectPath string) {
	// .gitignore
	err := createFile(path.Join(projectPath, ".gitignore"))
	exitOnError(err, "Failed to create "+name+".gitignore")

	const base64Gitignore = "IC5pZGVhLyoKCmRlYnVnCmNtZAoKI2tkaWZmMyBhbmQgb3RoZXIgbWVyZ2V0b29scwoqLm9yaWcKKi5CQUNLVVAuKgoqLkJBU0UuKgoqLkxPQ0FMLioKKi5SRU1PVEUuKgoKIyBCaW5hcmllcyBmb3IgcHJvZ3JhbXMgYW5kIHBsdWdpbnMKKi5leGUKKi5leGV+CiouZGxsCiouc28KKi5keWxpYgoKIyBUZXN0IGJpbmFyeSwgYnVpbGQgd2l0aCBgZ28gdGVzdCAtY2AKKi50ZXN0CgojIE91dHB1dCBvZiB0aGUgZ28gY292ZXJhZ2UgdG9vbCwgc3BlY2lmaWNhbGx5IHdoZW4gdXNlZCB3aXRoIExpdGVJREUKKi5vdXQK"
	gitignore, _ := decode(base64Gitignore)
	err = writeFile(path.Join(projectPath, ".gitignore"), gitignore)
	exitOnError(err, "Failed to write in "+name+".gitignore")

	// Dockerfile
	err = createFile(path.Join(projectPath, "Dockerfile"))
	exitOnError(err, "Failed to create "+name+"Dockerfile")
	const base64Dockerfile = "RlJPTSBnb2xhbmc6MS4xMi1hbHBpbmUgQVMgYnVpbGQKClJVTiBhcGsgLS1uby1jYWNoZSAtLXVwZGF0ZSBhZGQgZ2l0IG9wZW5zc2gtY2xpZW50IGNhLWNlcnRpZmljYXRlcyAmJiB1cGRhdGUtY2EtY2VydGlmaWNhdGVzCgpBREQgLiB7cGF0aH0KV09SS0RJUiB7cGF0aH0vY21kL3NlcnZpY2UKClJVTiBnaXQgY29uZmlnIC0tZ2xvYmFsIHVybC4iZ2l0QGdpdGxhYi5jb206Ii5pbnN0ZWFkT2YgImh0dHBzOi8vZ2l0bGFiLmNvbS8iClJVTiBta2RpciAvcm9vdC8uc3NoICYmIGVjaG8gIlN0cmljdEhvc3RLZXlDaGVja2luZyBubyAiID4gL3Jvb3QvLnNzaC9jb25maWcgCgpBUkcgU1NIX0tFWQpSVU4gc2V0ICt4ICYmIGVjaG8gIiRTU0hfS0VZIiA+IC9yb290Ly5zc2gvaWRfcnNhICYmIGNobW9kIDQwMCAvcm9vdC8uc3NoL2lkX3JzYQoKIyBFTlYgR08xMTFNT0RVTEUgb24KClJVTiBnbyBnZXQgLWQgLXY7IENHT19FTkFCTEVEPTAgR09PUz1saW51eCBnbyBidWlsZCAtYSAtaW5zdGFsbHN1ZmZpeCBjZ28gLW8gYXBwCgojIFJ1biBsaW50IGFuZCBnbyB2ZXQKRlJPTSBnb2xhbmc6MS4xMi1hbHBpbmUgQVMgbGludAoKUlVOIGFwayAtLW5vLWNhY2hlIC0tdXBkYXRlIGFkZCBnaXQKUlVOIGdvIGdldCAtdSBnb2xhbmcub3JnL3gvbGludC9nb2xpbnQKQ09QWSAtLWZyb209YnVpbGQgL2dvLyAvZ28vCldPUktESVIge3BhdGh9ClJVTiBnbyBnZXQgLXUgLXYgZ2l0aHViLmNvbS9wb3kvb25wYXIKUlVOIGdvIGdldCAtdSAtdiBnaXRodWIuY29tL3N0cmV0Y2hyL3Rlc3RpZnkvYXNzZXJ0ClJVTiBjaG1vZCAreCBsaW50LnNoIApSVU4gLi9saW50LnNoCgojIFJ1biB0ZXN0cwpGUk9NIGdvbGFuZzoxLjEyLWFscGluZSBBUyB0ZXN0ClJVTiBhcGsgLS1uby1jYWNoZSAtLXVwZGF0ZSBhZGQgZ2l0CkNPUFkgLS1mcm9tPWJ1aWxkIC9nby9zcmMvIC9nby9zcmMvCldPUktESVIge3BhdGh9ClJVTiBnbyBnZXQgLXUgLXYgZ2l0aHViLmNvbS9wb3kvb25wYXIKUlVOIGdvIGdldCAtdSAtdiBnaXRodWIuY29tL3N0cmV0Y2hyL3Rlc3RpZnkvYXNzZXJ0ClJVTiBjaG1vZCAreCBydW5UZXN0cy5zaCAKUlVOIC4vcnVuVGVzdHMuc2ggJHtTRVJWSUNFX05BTUV9ICR7U09OQVJfS0VZfQoKIyBSdW4gc29uYXJxdWJlCkZST00gamF2YTphbHBpbmUgYXMgc29uYXJxdWJlCkFSRyBTRVJWSUNFX05BTUUKQVJHIEJSQU5DSF9OQU1FCkFSRyBTT05BUl9LRVkKRU5WIFNPTkFSX1NDQU5ORVJfVkVSU0lPTiAzLjMuMC4xNDkyClJVTiBhcGsgYWRkIC0tbm8tY2FjaGUgd2dldCAmJiBcCiAgICB3Z2V0IGh0dHBzOi8vYmluYXJpZXMuc29uYXJzb3VyY2UuY29tL0Rpc3RyaWJ1dGlvbi9zb25hci1zY2FubmVyLWNsaS9zb25hci1zY2FubmVyLWNsaS0ke1NPTkFSX1NDQU5ORVJfVkVSU0lPTn0uemlwICYmIFwKICAgIHVuemlwIHNvbmFyLXNjYW5uZXItY2xpLSR7U09OQVJfU0NBTk5FUl9WRVJTSU9OfSAmJiBcCiAgICBjZCAvdXNyL2JpbiAmJiBsbiAtcyAvc29uYXItc2Nhbm5lci0ke1NPTkFSX1NDQU5ORVJfVkVSU0lPTn0vYmluL3NvbmFyLXNjYW5uZXIgc29uYXItc2Nhbm5lciAmJiBcCiAgICBhcGsgZGVsIHdnZXQgIApDT1BZIC0tZnJvbT1idWlsZCAvZ28vc3JjLyAvZ28vc3JjLwpDT1BZIC0tZnJvbT10ZXN0IHtwYXRofS9jb3ZlcmFnZS5vdXQgIHtwYXRofS9jb3ZlcmFnZS5vdXQgCldPUktESVIge3BhdGh9ClJVTiBjaG1vZCAreCBwdXNoU29uYXJxdWJlLnNoIApSVU4gLi9wdXNoU29uYXJxdWJlLnNoICR7U0VSVklDRV9OQU1FfSAke1NPTkFSX0tFWX0gJHtCUkFOQ0hfTkFNRX0KCkZST00gc2NyYXRjaCBBUyBzZXJ2aWNlCgpDT1BZICAtLWZyb209YnVpbGQgL2V0Yy9zc2wvY2VydHMvY2EtY2VydGlmaWNhdGVzLmNydCAvZXRjL3NzbC9jZXJ0cy8KQ09QWSAgLS1mcm9tPWJ1aWxkIHtwYXRofS9jbWQvc2VydmljZS9hcHAgL2Jpbi9hcHAKCkVYUE9TRSA4MDgwCgpDTUQgWyIvYmluL2FwcCJdCg=="
	dockerfile, _ := decode(base64Dockerfile)
	dockerfile = replacePlaceholder(dockerfile, "{path}", projectPath)

	err = writeFile(path.Join(projectPath, "Dockerfile"), dockerfile)
	exitOnError(err, "Failed to write in "+name+"Dockerfile")

	// Jenkinsfile
	err = createFile(path.Join(projectPath, "Jenkinsfile"))
	exitOnError(err, "Failed to create "+name+"Jenkinsfile")
	const base64Jenkinsfile = "IyFncm9vdnkKQExpYnJhcnkoJ2RvbWFpbi1kZWxpdmVyeS1waXBlbGluZScpIF8KCi8vIFNldCB0aGUgcHJvcGVyIHByb2plY3QgbmFtZSAoYWthIE5hbWVzcGFjZSkKU3RyaW5nIHByb2plY3QgPSAib25ib2FyZGluZyIKCi8vIEJ5IGRlZmF1bHQgdGhlIHNlcnZpY2UgbmFtZSBpcyBzZXQgdG8gdGhlIGplbmtpbnMgam9iX25hbWUKU3RyaW5nIHNlcnZpY2UgPSAicGxhdGZvcm0iCgovLyBEZWZpbmVzIHRoZSBkb2NrZXJpbWFnZSBuYW1lClN0cmluZyBpbWFnZUJhc2UgPSAiZG9tYWluLmNvbS8ke3Byb2plY3R9LyR7c2VydmljZX0iCgovLyBPcHRpb25hbCAoZm9yIGdpdGxhYiBidWlsZCBzdGF0dXMgb25seSkKcHJvcGVydGllcyhbCiAgICBnaXRMYWJDb25uZWN0aW9uKCdHaXRsYWIgQ29ubmVjdGlvbicpCl0pCgpub2RlIHsKICAgIFN0cmluZyBjb21taXRIYXNoCgogICAgdHJ5IHsKICAgICAgICAvLyBUaGlzIHN0YWdlIGNhbm5vdCBsaXZlIGluc2lkZSBnaXRsYWJCdWlsZHMtU3RhdHVzCiAgICAgICAgc3RhZ2UoJ0NoZWNrb3V0JykgewogICAgICAgICAgICBjb21taXRIYXNoID0gY2hlY2tvdXQoc2NtKS5HSVRfQ09NTUlUCiAgICAgICAgfQoKICAgICAgICBnaXRsYWJCdWlsZHMoYnVpbGRzOiBbIjFfYnVpbGQiLCAiMl9saW50IiAsIjNfdGVzdCIsICI0X3NvbmFycXViZSIsICI1X2ZpbmFsaXplIiwgIjZfcHVzaCIsICI3X2RlcGxveSJdKSB7CiAgICAgICAgICAgIHN0YWdlKCdCdWlsZCcpIHsKICAgICAgICAgICAgICAgIGdpdGxhYkNvbW1pdFN0YXR1cyhuYW1lOiAiMV9idWlsZCIpIHsKICAgICAgICAgICAgICAgICAgICB3aXRoQ3JlZGVudGlhbHMoW3NzaFVzZXJQcml2YXRlS2V5KGNyZWRlbnRpYWxzSWQ6ICdlbXB0eScsIGtleUZpbGVWYXJpYWJsZTogJ3NzaEtleScpXSkgewogICAgICAgICAgICAgICAgICAgICAgICAgYW5zaUNvbG9yKCd4dGVybScpIHsKICAgICAgICAgICAgICAgICAgICAgICAgICAgIHNoKCJkb2NrZXIgcHVsbCAke2ltYWdlQmFzZX06JHtjb21taXRIYXNofSB8fCAoc2V0ICt4ICYmIERPQ0tFUl9CVUlMREtJVD0xIGRvY2tlciBidWlsZCAtdCAke2ltYWdlQmFzZX06JHtjb21taXRIYXNofSAtLXRhcmdldCBidWlsZCAgLS1idWlsZC1hcmcgU1NIX0tFWT1cIlwkKGNhdCAke3NzaEtleX0pXCIgLikiKQogICAgICAgICAgICAgICAgICAgICAgICB9IAogICAgICAgICAgICAgICAgICAgIH0KICAgICAgICAgICAgICAgIH0KICAgICAgICAgICAgfQogICAgCiAgICAgICAgICAgIHN0YWdlKCdMaW50JykgewogICAgICAgICAgICAgICAgZ2l0bGFiQ29tbWl0U3RhdHVzKG5hbWU6ICIyX2xpbnQiKSB7CiAgICAgICAgICAgICAgICAgICAgd2l0aENyZWRlbnRpYWxzKFtzc2hVc2VyUHJpdmF0ZUtleShjcmVkZW50aWFsc0lkOiAnZW1wdHknLCBrZXlGaWxlVmFyaWFibGU6ICdzc2hLZXknKV0pIHsKICAgICAgICAgICAgICAgICAgICAgICAgYW5zaUNvbG9yKCd4dGVybScpIHsKICAgICAgICAgICAgICAgICAgICAgICAgICAgIHNoKCJzZXQgK3ggJiYgRE9DS0VSX0JVSUxES0lUPTEgZG9ja2VyIGJ1aWxkIC10ICR7aW1hZ2VCYXNlfTpsaW50LSR7Y29tbWl0SGFzaH0gLS10YXJnZXQgbGludCAtLWJ1aWxkLWFyZyBTU0hfS0VZPVwiXCQoY2F0ICR7c3NoS2V5fSlcIiAuIikKICAgICAgICAgICAgICAgICAgICAgICAgfQogICAgICAgICAgICAgICAgICAgIH0KICAgICAgICAgICAgICAgIH0KICAgICAgICAgICAgfQoKICAgICAgICAgICAgc3RhZ2UoJ1Rlc3QnKSB7CiAgICAgICAgICAgICAgICB3aXRoQ3JlZGVudGlhbHMoW3NzaFVzZXJQcml2YXRlS2V5KGNyZWRlbnRpYWxzSWQ6ICdlbXB0eScsIGtleUZpbGVWYXJpYWJsZTogJ3NzaEtleScpXSkgewogICAgICAgICAgICAgICAgICAgIGdpdGxhYkNvbW1pdFN0YXR1cyhuYW1lOiAiM190ZXN0IikgewogICAgICAgICAgICAgICAgICAgICAgICBhbnNpQ29sb3IoJ3h0ZXJtJykgewogICAgICAgICAgICAgICAgICAgICAgICAgICAgc2goInNldCAreCAmJiBET0NLRVJfQlVJTERLSVQ9MSBkb2NrZXIgYnVpbGQgLXQgJHtpbWFnZUJhc2V9OnRlc3QtJHtjb21taXRIYXNofSAtLXRhcmdldCB0ZXN0IC0tYnVpbGQtYXJnIFNTSF9LRVk9XCJcJChjYXQgJHtzc2hLZXl9KVwiIC4iKQogICAgICAgICAgICAgICAgICAgICAgICB9CiAgICAgICAgICAgICAgICAgICAgfQogICAgICAgICAgICAgICAgfQogICAgICAgICAgICB9CgogICAgICAgICAgICBzdGFnZSgnU29uYXJxdWJlJykgewogICAgICAgICAgICAgICAgd2l0aENyZWRlbnRpYWxzKFtzc2hVc2VyUHJpdmF0ZUtleShjcmVkZW50aWFsc0lkOiAnZW1wdHknLCBrZXlGaWxlVmFyaWFibGU6ICdzc2hLZXknKV0pIHsKICAgICAgICAgICAgICAgICAgICBnaXRsYWJDb21taXRTdGF0dXMobmFtZTogIjRfc29uYXJxdWJlIikgewogICAgICAgICAgICAgICAgICAgICAgICBhbnNpQ29sb3IoJ3h0ZXJtJykgewogICAgICAgICAgICAgICAgICAgICAgICAgICAgc2goInNldCAreCAmJiBET0NLRVJfQlVJTERLSVQ9MSBkb2NrZXIgYnVpbGQgLXQgJHtpbWFnZUJhc2V9OnRlc3QtJHtjb21taXRIYXNofSAtLXRhcmdldCBzb25hcnF1YmUgLS1idWlsZC1hcmcgU0VSVklDRV9OQU1FPW9uYm9hcmRpbmctcGxhdGZvcm0gLS1idWlsZC1hcmcgQlJBTkNIX05BTUU9JHtCUkFOQ0hfTkFNRX0gLS1idWlsZC1hcmcgU09OQVJfS0VZPTAyNjg0NTEyZTAxNjBmZjkyMjIyNmJjOGFkNGQ2NDNiZjZhYWNiZTAgLS1idWlsZC1hcmcgU1NIX0tFWT1cIlwkKGNhdCAke3NzaEtleX0pXCIgLiIpCiAgICAgICAgICAgICAgICAgICAgICAgIH0KICAgICAgICAgICAgICAgICAgICB9CiAgICAgICAgICAgICAgICB9CiAgICAgICAgICAgIH0KICAgICAgICAKCiAgICAgICAgICAgIHN0YWdlKCdGaW5hbGl6ZScpewogICAgICAgICAgICAgICAgd2l0aENyZWRlbnRpYWxzKFtzc2hVc2VyUHJpdmF0ZUtleShjcmVkZW50aWFsc0lkOiAnZW1wdHknLCBrZXlGaWxlVmFyaWFibGU6ICdzc2hLZXknKV0pIHsKICAgICAgICAgICAgICAgICAgICBnaXRsYWJDb21taXRTdGF0dXMobmFtZTogIjVfZmluYWxpemUiKSB7CiAgICAgICAgICAgICAgICAgICAgICAgIGFuc2lDb2xvcigneHRlcm0nKSB7CiAgICAgICAgICAgICAgICAgICAgICAgICAgICBzaCgic2V0ICt4ICYmIERPQ0tFUl9CVUlMREtJVD0xIGRvY2tlciBidWlsZCAtdCAke2ltYWdlQmFzZX06JHtjb21taXRIYXNofSAtLXRhcmdldCBzZXJ2aWNlIC0tYnVpbGQtYXJnIFNTSF9LRVk9XCJcJChjYXQgJHtzc2hLZXl9KVwiIC4iKQogICAgICAgICAgICAgICAgICAgICAgICB9CiAgICAgICAgICAgICAgICAgICAgfQogICAgICAgICAgICAgICAgfQogICAgICAgICAgICB9ICAgCgogICAgICAgICAgICBzdGFnZSgnUHVzaCcpIHsKICAgICAgICAgICAgICAgIGdpdGxhYkNvbW1pdFN0YXR1cyhuYW1lOiAiNl9wdXNoIikgewogICAgICAgICAgICAgICAgICAgIGZpbm9QdXNoKAogICAgICAgICAgICAgICAgICAgICAgICBjb21taXQ6IGNvbW1pdEhhc2gsCiAgICAgICAgICAgICAgICAgICAgICAgIGltYWdlOiBpbWFnZUJhc2UKICAgICAgICAgICAgICAgICAgICApCiAgICAgICAgICAgICAgICB9CiAgICAgICAgICAgIH0KCiAgICAgICAgICAgIHN0YWdlKCdEZXBsb3knKSB7CiAgICAgICAgICAgICAgICBnaXRsYWJDb21taXRTdGF0dXMobmFtZTogIjdfZGVwbG95IikgewogICAgICAgICAgICAgICAgICAgIGZpbm9EZXBsb3koCiAgICAgICAgICAgICAgICAgICAgICAgIGNvbW1pdDogY29tbWl0SGFzaCwKICAgICAgICAgICAgICAgICAgICAgICAgcHJvamVjdDogcHJvamVjdCwKICAgICAgICAgICAgICAgICAgICAgICAgc2VydmljZTogc2VydmljZSwKICAgICAgICAgICAgICAgICAgICAgICAgY29uZmlnczogImNvbmZpZ3MvIgogICAgICAgICAgICAgICAgICAgICkKICAgICAgICAgICAgICAgIH0KICAgICAgICAgICAgfQogICAgICAgIH0KICAgIH0KICAgIC8vIE9wdGlvbmFsbHkgY2F0Y2ggYSBwb3NzaWJsZSBlcnJvcgogICAgY2F0Y2ggKGUpIHsKICAgICAgICBjdXJyZW50QnVpbGQucmVzdWx0ID0gJ0ZBSUxVUkUnCiAgICB9CiAgICBmaW5hbGx5IHsKICAgICAgICBmaW5vQ2xlYW51cChpbWFnZUJhc2UpCiAgICB9Cn0K"
	jenkinsfile, _ := decode(base64Jenkinsfile)
	jenkinsfile = replacePlaceholder(jenkinsfile, "{name}", name)
	jenkinsfile = replacePlaceholder(jenkinsfile, "{namespace}", nameSpace)

	err = writeFile(path.Join(projectPath, "Jenkinsfile"), jenkinsfile)
	exitOnError(err, "Failed to write in "+name+"Jenkinsfile")

	// lint.sh
	err = createFile(path.Join(projectPath, "lint.sh"))
	exitOnError(err, "Failed to create "+name+"lint.sh")
	const base64LintSh = "IyEvYmluL3NoCmVjaG8gIiMjIyMjIyMjIyMjIyMjIyMjIyMjIGdvIGxpbnQgIyMjIyMjIyMjIyMjIyMjIyMjIyMiCiMgRmFpbCBidWlsZCBvbiBlcnJvcgovZ28vYmluL2dvbGludCAtc2V0X2V4aXRfc3RhdHVzIC4vLi4uCiMgRG9uJ3QgZmFpbCBidWlsZCBvbiBlcnJvcgojL2dvL2Jpbi9nb2xpbnQgLi8uLi4KCmlmIFsgJD8gLWVxIDEgXTsgdGhlbgogICAgZWNobyAiIyMjIyMjIyMjIyMjIyMjIyMjIyMgZ28gbGludCBGYWlsZWQgIyMjIyMjIyMjIyMjIyMjIyMjIyMiICAgIAogICAgZXhpdCAxCmZpCgoKCmVjaG8gIiMjIyMjIyMjIyMjIyMjIyMjIyMjIGdvIHZldCAjIyMjIyMjIyMjIyMjIyMjIyMjIyIKZ28gdmV0IC4vLi4uCgppZiBbICQ/IC1lcSAxIF07IHRoZW4KICAgIGVjaG8gIiMjIyMjIyMjIyMjIyMjIyMjIyMjIGdvIHZldCBGYWlsZWQgIyMjIyMjIyMjIyMjIyMjIyMjIyMiICAgIAogICAgZXhpdCAxCmZpCg=="
	lintSh, _ := decode(base64LintSh)

	err = writeFile(path.Join(projectPath, "lint.sh"), lintSh)
	exitOnError(err, "Failed to write in "+name+"lint.sh")

	// pushSonarqube.sh
	err = createFile(path.Join(projectPath, "pushSonarqube.sh"))
	exitOnError(err, "Failed to create "+name+"pushSonarqube.sh")
	const base64PushSonarqubeSh = "IyEvYmluL3NoCmVjaG8gIiMjIyMjIyMjIyMjIyMjIyMjIyMjIFJlcG9ydCBUbyBTb25hcnF1YmUgIyMjIyMjIyMjIyMjIyMjIyMjIyMiCmVjaG8gIiMjIyMjIyMjIyMjIyMjIyMjIyMjIEJyYW5jaCBpcyAkMyAjIyMjIyMjIyMjIyMjIyMjIyMjIyIgCmlmIFsgIiQzIiAhPSAnbWFzdGVyJyBdOyB0aGVuCiAgICBlY2hvICIjIyMjIyMjIyMjIyMjIyMjIyMjIyBCcmFuY2ggaXMgbm90IG1hc3Rlci4gU2tpcHBpbmcuLi4gIyMjIyMjIyMjIyMjIyMjIyMjIyMiICAgIAogICAgZXhpdCAwCmZpCi91c3IvYmluL3NvbmFyLXNjYW5uZXIgLURzb25hci5wcm9qZWN0S2V5PSQxIC1Ec29uYXIuc291cmNlcz0uIC1Ec29uYXIuaG9zdC51cmw9aHR0cHM6Ly9zb25hcmNsb3VkLmlvIC1Ec29uYXIubG9naW49IiQyIiAtRHNvbmFyLm9yZ2FuaXphdGlvbj1maW5vIC1Ec29uYXIuZ28uY292ZXJhZ2UucmVwb3J0UGF0aHM9ImNvdmVyYWdlLm91dCIgIC1Ec29uYXIuZXhjbHVzaW9ucz0iKiovdmVuZG9yLyoqLCoqLypfdGVzdC5nbyxjb25maWcvKiosZG9jcy8qKixyZXNvdXJjZXMvKioiIC1Ec29uYXIudGVzdHM9Ii4iIC1Ec29uYXIudGVzdC5pbmNsdXNpb25zPSIqKi8qX3Rlc3QuZ28iICAtRHNvbmFyLnRlc3QuZXhjbHVzaW9ucz0iKiovdmVuZG9yLyoqLGNvbmZpZy8qKixkb2NzLyoqLHJlc291cmNlcy8qKiIKIyBVc2UgdGhlIGZvbGxvd2luZyBsaW5lIHRvIGRlYnVnCiMvdXNyL2Jpbi9zb25hci1zY2FubmVyIC1YIC1Ec29uYXIucHJvamVjdEtleT0kMSAtRHNvbmFyLnNvdXJjZXM9LiAtRHNvbmFyLmhvc3QudXJsPWh0dHBzOi8vc29uYXJjbG91ZC5pbyAtRHNvbmFyLmxvZ2luPSIkMiIgLURzb25hci5vcmdhbml6YXRpb249ZmlubyAtRHNvbmFyLmdvLmNvdmVyYWdlLnJlcG9ydFBhdGhzPSJjb3ZlcmFnZS5vdXQiICAtRHNvbmFyLmV4Y2x1c2lvbnM9IioqL3ZlbmRvci8qKiwqKi8qX3Rlc3QuZ28sY29uZmlnLyoqLGRvY3MvKioscmVzb3VyY2VzLyoqIiAtRHNvbmFyLnRlc3RzPSIuIiAtRHNvbmFyLnRlc3QuaW5jbHVzaW9ucz0iKiovKl90ZXN0LmdvIiAgLURzb25hci50ZXN0LmV4Y2x1c2lvbnM9IioqL3ZlbmRvci8qKixjb25maWcvKiosZG9jcy8qKixyZXNvdXJjZXMvKioiCmlmIFsgJD8gLWVxIDEgXTsgdGhlbgogICAgZWNobyAiIyMjIyMjIyMjIyMjIyMjIyMjIyMgUmVwb3J0IFRvIEZhaWxlZCAjIyMjIyMjIyMjIyMjIyMjIyMjIyIgICAgCiAgICBleGl0IDEKZmkKZWNobyAiIyMjIyMjIyMjIyMjIyMjIyMjIyMgRmluaXNoZWQgUmVwb3J0aW5nIFRvIFNvbmFycXViZSAjIyMjIyMjIyMjIyMjIyMjIyMjIyI="
	pushSonarqubeSh, _ := decode(base64PushSonarqubeSh)

	err = writeFile(path.Join(projectPath, "pushSonarqube.sh"), pushSonarqubeSh)
	exitOnError(err, "Failed to write in "+name+"pushSonarqube.sh")

	// runTests.sh
	err = createFile(path.Join(projectPath, "runTests.sh"))
	exitOnError(err, "Failed to create "+name+"runTests.sh")
	const base64RunTestsSh = "IyEvYmluL3NoCgplY2hvICIjIyMjIyMjIyMjIyMjIyMjIyMjIyBSdW5uaW5nIFRlc3RzICMjIyMjIyMjIyMjIyMjIyMjIyMjIgpDR09fRU5BQkxFRD0wIGdvIHRlc3QgLXNob3J0IC1jb3ZlciAtdiAtY292ZXJwcm9maWxlPWNvdmVyYWdlLm91dCAtY292ZXJtb2RlPWF0b21pYyAuLy4uLgppZiBbICQ/IC1lcSAxIF07IHRoZW4KICAgIGVjaG8gIiMjIyMjIyMjIyMjIyMjIyMjIyMjIFJ1bm5pbmcgVGVzdHMgRmFpbGVkICMjIyMjIyMjIyMjIyMjIyMjIyMjIiAgICAKICAgIGV4aXQgMQpmaQplY2hvICIjIyMjIyMjIyMjIyMjIyMjIyMjIyBGaW5pc2hlZCBUZXN0cyAjIyMjIyMjIyMjIyMjIyMjIyMjIyIKCmV4aXQgMA=="
	runTestsSh, _ := decode(base64RunTestsSh)

	err = writeFile(path.Join(projectPath, "runTests.sh"), runTestsSh)
	exitOnError(err, "Failed to write in "+name+"runTests.sh")
}

// prompt returns parameters
// func prompt() (string, string, string, string) {
func prompt() (string, string, string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Name: ")
	scanner.Scan()
	name := scanner.Text()
	// fmt.Print("Name: ")
	// var name string
	// fmt.Scanf("%s", &name)
	// fmt.Println(name)
	fmt.Print("Namespace: ")
	var nameSpace string
	fmt.Scanln(&nameSpace)
	dir, err := os.Getwd()
	exitOnError(err, "Failed to get working directory")
	fmt.Printf("Path: [%s] (hit enter/change it)", dir)
	var givenPath string
	fmt.Scanln(&givenPath)
	return name, nameSpace, givenPath
}
