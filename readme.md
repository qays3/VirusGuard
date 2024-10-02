### Project Directory Structure
```
/VirusGuard/
├── main.go
└── process/
    ├── block_process.sh
    └── docker_containment.sh
└── install_dependencies.sh
```

### 1. `main.go`
This Go program includes an animated progress bar during the YARA scan process.

```go
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
	ruleFiles := getAllYaraRules("./rules")
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
	containerName := fmt.Sprintf("container_%s", malwareName)
	cmd := exec.Command("bash", "-c", fmt.Sprintf("./process/docker_containment.sh %s %s", malwareName, containerName))
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing Docker containment: %s\n", err)
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
```

### 2. `process/block_process.sh`
This script remains unchanged.

```bash
#!/bin/bash

MALWARE_NAME=$1

if [ -z "$MALWARE_NAME" ]; then
    echo "No malware name provided."
    exit 1
fi

PIDS=$(pgrep -f "$MALWARE_NAME")

if [ -z "$PIDS" ]; then
    echo "No processes found for $MALWARE_NAME."
else
    echo "Terminating processes for $MALWARE_NAME with PIDs: $PIDS"
    kill $PIDS
    echo "Processes terminated."
fi
```

### 3. `process/docker_containment.sh`
This script also remains unchanged.

```bash
#!/bin/bash

MALWARE_NAME=$1
CONTAINER_NAME=$2

if [ -z "$MALWARE_NAME" ] || [ -z "$CONTAINER_NAME" ]; then
    echo "Usage: $0 <malware_name> <container_name>"
    exit 1
fi

mkdir -p "$CONTAINER_NAME"
docker run --name "$CONTAINER_NAME" -d -it --rm -v "$(pwd)/$CONTAINER_NAME:/malware" ubuntu:latest /bin/bash

docker cp "$MALWARE_NAME" "$CONTAINER_NAME:/malware/$MALWARE_NAME"
docker exec "$CONTAINER_NAME" /malware/"$MALWARE_NAME"
```

### 4. `install_dependencies.sh`
This script now includes a download animation bar for the installation process.

```bash
#!/bin/bash

function showProgressBar {
    echo -n "Installing dependencies... ["
    for i in {1..50}; do
        sleep 0.1
        echo -n "="
    done
    echo "] Done."
}

sudo apt update
showProgressBar

sudo apt install -y docker.io
sudo systemctl start docker
sudo systemctl enable docker

echo "Docker has been installed and started."
```

### 5. Usage Instructions
1. **Set up your directory structure**:
   ```
   mkdir -p VirusGuard/process
   ```

2. **Place the files accordingly**:
   - Place `main.go` in `VirusGuard/`.
   - Place `block_process.sh` and `docker_containment.sh` in `VirusGuard/process/`.
   - Place `install_dependencies.sh` in `VirusGuard/`.

3. **Make the scripts executable**:
   ```bash
   chmod +x process/block_process.sh
   chmod +x process/docker_containment.sh
   chmod +x install_dependencies.sh
   ```

4. **Run the installation script**:
   ```bash
   ./install_dependencies.sh
   ```

5. **Compile the Go program**:
   ```bash
   cd VirusGuard
   go build -o VirusGuard main.go
   ```

6. **Run the VirusGuard tool**:
   ```bash
   ./VirusGuard --malware <malware_name> --action <SignatureBlocking|ThreadInterruption|DockerContainment>
   ```

### Example Usage
To run the tool, you can execute:
```bash
./VirusGuard --malware mymalware.exe --action DockerContainment
```
 