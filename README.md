
## VirusGuard

![VirusGuard](img/VirusGuard.png)

**VirusGuard** is part of the **MalWatch Suite**, a tool specifically designed for malware scanning and providing a safe environment for protection. VirusGuard can terminate malicious processes running in the background and track signatures to identify potential threats. If you want to execute malware in a virtual machine without the safe environment, you can still monitor all the threads and stop them by tracking the signature.

---

### Project Directory Structure

```
/VirusGuard/
├── main.go
├── .gitignore
├── README.md
├── setup.sh
├── img/
├── YaraRules/
├── process/
│   ├── terminate_process.sh
│   ├── signature_control.sh
│   └── docker_containment.sh

```

### 1. `main.go`

The core of VirusGuard, containing the processes that handle malware scanning, YARA scanning, process termination, Docker containment, and signature blocking/unblocking. The logger tracks relevant events and errors, making it easier to monitor activities and threats.

[main.go](main.go)

---

### 2. `process/terminate_process.sh`

This script terminates processes identified as malware.

[process/terminate_process.sh](process/terminate_process.sh)

---

### 3. `process/docker_containment.sh`

This script isolates suspicious processes in Docker containers for safe analysis.

[process/docker_containment.sh](process/docker_containment.sh)

---

### 4. `setup.sh`

Sets up the required environment and dependencies for VirusGuard.

[setup.sh](setup.sh)

---

### 5. Usage Instructions

1. **Make the scripts executable**:
   ```bash
   chmod +x process/terminate_process.sh
   chmod +x process/docker_containment.sh
   chmod +x setup.sh
   ```

2. **Run the installation script**:
   ```bash
   ./setup.sh
   ```

3. **Compile the Go program**:
   ```bash
   cd VirusGuard
   go build -o VirusGuard main.go
   ```

4. **Run the VirusGuard tool**:
   ```bash
   sudo ./VirusGuard --malware <malware_name> --action <YaraScan|SignatureBlocking|TerminateProcess|DockerContainment>
   ```

---

### Example Usage

To contain malware in Docker:
```bash
sudo ./VirusGuard --malware mymalware.exe --action DockerContainment
```

---

### New Features

- **Signature Blocking**: Enables the blocking of malware based on specific signatures.
- **Background Process Management**: Use the `--block` and `--unblock` options to control malware threads in the background.

---

### Additional Commands

- To block a malware signature:
  ```bash
  sudo ./VirusGuard --malware <malware_name> --action SignatureBlocking --block
  ```

- To unblock a malware signature:
  ```bash
  sudo ./VirusGuard --malware <malware_name> --action SignatureBlocking --unblock
  ```
