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
	"sync"
	"time"
)

type Logger struct {
	mu      sync.Mutex
	logFile *os.File
}

func AddLogger(filePath string) (*Logger, error) {
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		return nil, err
	}
	logFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return &Logger{logFile: logFile}, nil
}

func (l *Logger) Log(message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("%s: %s\n", timestamp, message)
	if _, err := l.logFile.WriteString(logMessage); err != nil {
		fmt.Println("Error writing to log file:", err)
	}
}

func (l *Logger) Close() {
	l.logFile.Close()
}

func main() {
	malwareName := flag.String("malware", "", "Name of the malware to scan and block")
	action := flag.String("action", "", "Choose an action: YaraScan, TerminateProcess, DockerContainment, BlockSignature")
	block := flag.Bool("block", false, "Block the malware")
	unblock := flag.Bool("unblock", false, "Unblock the malware")
	help := flag.Bool("help", false, "Display help for VirusGuard")

	flag.Parse()

	logger, err := AddLogger("logs/application.log")
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}
	defer logger.Close()

	if *help {
		logger.Log("Displayed help")
		fmt.Println("Usage: VirusGuard [OPTIONS]")
		fmt.Println("Options:")
		fmt.Println("  --malware <name>         Name of the malware file to handle")
		fmt.Println("  --action <type>          Type of action to perform:")
		fmt.Println("                           YaraScan: Scan the malware with YARA rules.")
		fmt.Println("                           TerminateProcess: Terminate all running processes associated with the malware.")
		fmt.Println("                           DockerContainment: Run the malware in a Docker container to isolate its execution.")
		fmt.Println("                           BlockSignature: Block the malware while under analysis.")
		fmt.Println("  --block                  Block the specified malware.")
		fmt.Println("  --unblock                Unblock the specified malware.")
		fmt.Println("  --help                   Show this help message")
		return
	}

	if *malwareName == "" || *action == "" {
		logger.Log("Error: both --malware and --action options must be provided")
		fmt.Println("Error: both --malware and --action options must be provided")
		return
	}

	switch strings.ToLower(*action) {
	case "yarascan":
		yaraScan(*malwareName, logger)
	case "terminateprocess":
		TerminateProcess(*malwareName, logger)
	case "dockercontainment":
		dockerContainment(*malwareName, logger)
	case "blocksignature":
		if *block {
			signature := calculateSignature(*malwareName)
			blockSignature(*malwareName, signature, logger)
		} else if *unblock {
			unblockSignature(*malwareName, logger)
		} else {
			logger.Log("Error: Specify --block or --unblock with --action BlockSignature")
			fmt.Println("Error: Specify --block or --unblock with --action BlockSignature")
		}
	default:
		logger.Log("Invalid action specified. Use --help to see the available options.")
		fmt.Println("Invalid action specified. Use --help to see the available options.")
	}
}

func yaraScan(malwareName string, logger *Logger) {
	ruleFiles := getAllYaraRules("./YaraRules")
	for _, ruleFile := range ruleFiles {
		fmt.Printf("Scanning with YARA rule: %s\n", ruleFile)
		showProgressBar("YARA scan in progress", 5)
		cmd := exec.Command("bash", "-c", fmt.Sprintf("yara %s %s", ruleFile, malwareName))
		output, err := cmd.Output()
		if err != nil {
			logger.Log(fmt.Sprintf("Error executing YARA scan for %s: %s", ruleFile, err))
			fmt.Printf("Error executing YARA scan for %s: %s\n", ruleFile, err)
		}
		logger.Log(fmt.Sprintf("YARA scan result for %s: %s", ruleFile, string(output)))
		fmt.Printf("YARA scan result for %s: %s\n", ruleFile, string(output))
	}
}

func TerminateProcess(malwareName string, logger *Logger) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("./process/terminate_process.sh %s", malwareName))
	output, err := cmd.Output()
	if err != nil {
		logger.Log(fmt.Sprintf("Error executing thread interruption: %s", err))
		fmt.Printf("Error executing thread interruption: %s\n", err)
		return
	}
	logger.Log(fmt.Sprintf("Thread interruption result: %s", string(output)))
	fmt.Printf("Thread interruption result: %s\n", string(output))
}

func dockerContainment(malwareName string, logger *Logger) {
	err := exec.Command("chmod", "+x", malwareName).Run()
	if err != nil {
		logger.Log(fmt.Sprintf("Error making malware executable: %s", err))
		fmt.Printf("Error making malware executable: %s\n", err)
		return
	}

	containerName := fmt.Sprintf("container_%s", malwareName)
	cmd := exec.Command("bash", "-c", fmt.Sprintf("./process/docker_containment.sh %s %s", malwareName, containerName))
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log(fmt.Sprintf("Error executing Docker containment: %s\nOutput: %s", err, string(output)))
		fmt.Printf("Error executing Docker containment: %s\nOutput: %s\n", err, string(output))
		return
	}
	logger.Log(fmt.Sprintf("Docker containment result: %s", string(output)))
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

func blockSignature(malwareName, signature string, logger *Logger) {
	blockCmd := exec.Command("bash", "-c", fmt.Sprintf("nohup ./process/signature_control.sh %s %s %s &", malwareName, signature, "block"))
	if err := blockCmd.Start(); err != nil {
		logger.Log(fmt.Sprintf("Error starting block signature: %s", err))
		fmt.Printf("Error starting block signature: %s\n", err)
		return
	}
	logger.Log(fmt.Sprintf("Blocking malware: %s is under analysis with signature: %s. Use --unblock to stop.", malwareName, signature))
	fmt.Printf("Blocking malware: %s is under analysis with signature: %s. Use --unblock to stop.\n", malwareName, signature)
}

func unblockSignature(malwareName string, logger *Logger) {
	unblockCmd := exec.Command("bash", "-c", fmt.Sprintf("./process/signature_control.sh %s unblock", malwareName))
	output, err := unblockCmd.Output()
	if err != nil {
		logger.Log(fmt.Sprintf("Error executing unblock signature: %s", err))
		fmt.Printf("Error executing unblock signature: %s\n", err)
		return
	}
	logger.Log(fmt.Sprintf("Unblock result: %s", string(output)))
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
