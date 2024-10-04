

# VirusGuard

## Project Directory Structure
```
/VirusGuard/
├── main.go
├── process/
│   ├── terminate_process.sh
│   └── docker_containment.sh
└── setup.sh
```

### 1. `main.go`
This Go program includes an animated progress bar during the YARA scan process.

[main](main.go)

### 2. `process/terminate_process.sh`
 

[process/block_process](process/terminate_process.sh)

### 3. `process/docker_containment.sh`
 

[process/docker_containment](process/docker_containment.sh)

### 4. `setup.sh`
 

[install_requirements](setup.sh)

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
   sudo ./VirusGuard --malware <malware_name> --action <SignatureBlocking|TerminateProcess|DockerContainment>
   ```

### Example Usage
To run the tool, you can execute:
```bash
sudo ./VirusGuard --malware mymalware.exe --action DockerContainment
```

### New Features
- **Signature Blocking**: Added the ability to block specific malware signatures while analyzing the malware.
- **Background Process Management**: The `--block` and `--unblock` commands allow managing malware processes as background threads.

### Additional Commands
- To block a malware signature:
  ```bash
  sudo ./VirusGuard --malware <malware_name> --action SignatureBlocking --block
  ```

- To unblock a malware signature:
  ```bash
  sudo ./VirusGuard --malware <malware_name> --action SignatureBlocking --unblock
  ```
 