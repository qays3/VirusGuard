### Project Directory Structure
```
/VirusGuard/
├── main.go
└── process/
    ├── block_process.sh
    └── docker_containment.sh
└── setup.sh
```

### 1. `main.go`
This Go program includes an animated progress bar during the YARA scan process.

[main](main.go)

### 2. `process/block_process.sh`
This script remains unchanged.

[process/block_process](process/block_process.sh)

### 3. `process/docker_containment.sh`
This script also remains unchanged.

[process/docker_containment](process/docker_containment.sh)


### 4. `setup.sh`
This script now includes a download animation bar for the installation process.

[install_requirements](setup.sh)

### 5. Usage Instructions


1. **Make the scripts executable**:
   ```bash
   chmod +x process/block_process.sh
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
   sudo ./VirusGuard --malware <malware_name> --action <SignatureBlocking|ThreadInterruption|DockerContainment>
   ```

### Example Usage
To run the tool, you can execute:
```bash
sudo ./VirusGuard --malware mymalware.exe --action DockerContainment
```
 