package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	malwareName := flag.String("malware", "", "Name of the malware to scan and block")
	action := flag.String("action", "", "Choose an action: SignatureBlocking, ThreadInterruption, DockerContainment")
	help := flag.Bool("help", false, "Display help for VirusGuard")

	flag.Parse()

	if *help {
		fmt.Println("Usage: VirusGuard [OPTIONS]")
		fmt.Println("Options:")
		fmt.Println("  --malware <name>         Name of the malware file to handle")
		fmt.Println("  --action <type>          Type of action to perform:")
		fmt.Println("                           SignatureBlocking: Block malware using YARA signatures.")
		fmt.Println("                           ThreadInterruption: Terminate all running processes associated with the malware.")
		fmt.Println("                           DockerContainment: Run the malware in a Docker container to isolate its execution.")
		fmt.Println("  --help                   Show this help message")
		return
	}

	if *malwareName == "" || *action == "" {
		fmt.Println("Error: both --malware and --action options must be provided")
		return
	}

	switch strings.ToLower(*action) {
	case "signatureblocking":
		signatureBlocking(*malwareName)
	case "threadinterruption":
		threadInterruption(*malwareName)
	case "dockercontainment":
		dockerContainment(*malwareName)
	default:
		fmt.Println("Invalid action specified. Use --help to see the available options.")
	}
}

func signatureBlocking(malwareName string) {
	ruleFiles := getAllYaraRules("./YaraRules")
	for _, ruleFile := range ruleFiles {
		fmt.Printf("Scanning with YARA rule: %s\n", ruleFile)
		showProgressBar("YARA scan in progress", 5)
		cmd := exec.Command("bash", "-c", fmt.Sprintf("yara %s %s", ruleFile, malwareName))
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("Error executing YARA signature blocking for %s: %s\n", ruleFile, err)
		}
		fmt.Printf("YARA scan result for %s: %s\n", ruleFile, string(output))
	}
}

func threadInterruption(malwareName string) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("./process/block_process.sh %s", malwareName))
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing thread interruption: %s\n", err)
		return
	}
	fmt.Printf("Thread interruption result: %s\n", string(output))
}

func dockerContainment(malwareName string) {

	err := exec.Command("chmod", "+x", malwareName).Run()
	if err != nil {
		fmt.Printf("Error making malware executable: %s\n", err)
		return
	}

	containerName := fmt.Sprintf("container_%s", malwareName)
	cmd := exec.Command("bash", "-c", fmt.Sprintf("./process/docker_containment.sh %s %s", malwareName, containerName))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing Docker containment: %s\nOutput: %s\n", err, string(output))
		return
	}
	fmt.Printf("Docker containment result: %s\n", string(output))
}

func getAllYaraRules(ruleDir string) []string {
	var ruleFiles []string
	err := filepath.Walk(ruleDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), ".yar") {
			ruleFiles = append(ruleFiles, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error scanning rule directory: %s\n", err)
		return nil
	}
	return ruleFiles
}

func showProgressBar(message string, totalSteps int) {
	fmt.Print(message + " [")
	for i := 0; i <= totalSteps; i++ {
		time.Sleep(1 * time.Second)
		fmt.Print("=")
	}
	fmt.Println("] Done.")
}
