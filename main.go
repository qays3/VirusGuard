package main

import (
	"crypto/sha256"
	"encoding/hex"
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
	action := flag.String("action", "", "Choose an action: YaraScan, ThreadInterruption, DockerContainment, BlockSignature")
	block := flag.Bool("block", false, "Block the malware")
	unblock := flag.Bool("unblock", false, "Unblock the malware")
	help := flag.Bool("help", false, "Display help for VirusGuard")

	flag.Parse()

	if *help {
		fmt.Println("Usage: VirusGuard [OPTIONS]")
		fmt.Println("Options:")
		fmt.Println("  --malware <name>         Name of the malware file to handle")
		fmt.Println("  --action <type>          Type of action to perform:")
		fmt.Println("                           YaraScan: Scan the malware with YARA rules.")
		fmt.Println("                           ThreadInterruption: Terminate all running processes associated with the malware.")
		fmt.Println("                           DockerContainment: Run the malware in a Docker container to isolate its execution.")
		fmt.Println("                           BlockSignature: Block the malware while under analysis.")
		fmt.Println("  --block                  Block the specified malware.")
		fmt.Println("  --unblock                Unblock the specified malware.")
		fmt.Println("  --help                   Show this help message")
		return
	}

	if *malwareName == "" || *action == "" {
		fmt.Println("Error: both --malware and --action options must be provided")
		return
	}

	switch strings.ToLower(*action) {
	case "yarascan":
		yaraScan(*malwareName)
	case "threadinterruption":
		threadInterruption(*malwareName)
	case "dockercontainment":
		dockerContainment(*malwareName)
	case "blocksignature":
		if *block {
			signature := calculateSignature(*malwareName)
			blockSignature(*malwareName, signature)
		} else if *unblock {
			unblockSignature(*malwareName)
		} else {
			fmt.Println("Error: Specify --block or --unblock with --action BlockSignature")
		}
	default:
		fmt.Println("Invalid action specified. Use --help to see the available options.")
	}
}

func yaraScan(malwareName string) {
	ruleFiles := getAllYaraRules("./YaraRules")
	for _, ruleFile := range ruleFiles {
		fmt.Printf("Scanning with YARA rule: %s\n", ruleFile)
		showProgressBar("YARA scan in progress", 5)
		cmd := exec.Command("bash", "-c", fmt.Sprintf("yara %s %s", ruleFile, malwareName))
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("Error executing YARA scan for %s: %s\n", ruleFile, err)
		}
		fmt.Printf("YARA scan result for %s: %s\n", ruleFile, string(output))
	}
}

func threadInterruption(malwareName string) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("./process/terminate_process.sh %s", malwareName))
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

func calculateSignature(malwareName string) string {
	fileInfo, err := os.Stat(malwareName)
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return ""
	}

	hash := sha256.New()
	hash.Write([]byte(malwareName))
	hash.Write([]byte(fmt.Sprintf("%d", fileInfo.Size())))
	return hex.EncodeToString(hash.Sum(nil))
}

func blockSignature(malwareName, signature string) {
	blockCmd := exec.Command("bash", "-c", fmt.Sprintf("nohup ./process/signature_control.sh %s %s %s &", malwareName, signature, "block"))
	if err := blockCmd.Start(); err != nil {
		fmt.Printf("Error starting block signature: %s\n", err)
		return
	}
	fmt.Printf("Blocking malware: %s is under analysis with signature: %s. Use --unblock to stop.\n", malwareName, signature)
}

func unblockSignature(malwareName string) {
	unblockCmd := exec.Command("bash", "-c", fmt.Sprintf("./process/signature_control.sh %s unblock", malwareName))
	output, err := unblockCmd.Output()
	if err != nil {
		fmt.Printf("Error executing unblock signature: %s\n", err)
		return
	}
	fmt.Printf("Unblock result: %s\n", string(output))
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
