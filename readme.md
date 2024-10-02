### Project Directory Structure
```
/VirusGuard/
├── main.go
└── process/
    ├── block_process.sh
    └── docker_containment.sh
└── install_requirements.sh
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


### 4. `install_requirements.sh`
This script now includes a download animation bar for the installation process.

[install_requirements](install_requirements.sh)

### 5. Usage Instructions
1. **Set up your directory structure**:
   ```
   mkdir -p VirusGuard/process
   ```

2. **Place the files accordingly**:
   - Place `main.go` in `VirusGuard/`.
   - Place `block_process.sh` and `docker_containment.sh` in `VirusGuard/process/`.
   - Place `install_requirements.sh` in `VirusGuard/`.

3. **Make the scripts executable**:
   ```bash
   chmod +x process/block_process.sh
   chmod +x process/docker_containment.sh
   chmod +x install_requirements.sh
   ```

4. **Run the installation script**:
   ```bash
   ./install_requirements.sh
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
 
