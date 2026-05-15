# PeritiaGo Security Limitations

While PeritiaGo provides robust user-mode protection, investigators should be aware of the following technical limitations:

## 1. User-Mode Scope
The Guard system operates entirely in **User Mode**. 
- **Limitation**: It cannot block writes initiated by kernel drivers or other processes with SYSTEM-level privileges if they operate outside the Go runtime.
- **Context**: PeritiaGo protects itself and the system from its own code; it is not a system-wide EDR/Antivirus.

## 2. Low-Level Disk Access
- **Limitation**: Direct sector-level writes to physical drives (e.g., `\\.\PhysicalDrive0`) bypass the filesystem-level Guard.
- **Status**: PeritiaGo does not currently implement raw disk write capabilities, but investigators should be aware that raw access bypasses path-based security.

## 3. Pre-Initialization Window
- **Limitation**: During the very first milliseconds of execution (configuration loading and logger initialization), the Guard may not be fully active.
- **Mitigation**: Configuration and logging are restricted to the application's own directory.

## 4. Environment Dependency
- **Limitation**: The effectiveness of symlink and junction protection on Windows depends on the OS version and the permissions of the user running the tool.
- **Warning**: On older versions of Windows where `EvalSymlinks` might behave differently, some edge cases might exist. Always run as Administrator for maximum accuracy.

## 5. Race Conditions (TOCTOU)
- **Limitation**: Like all path-based checks, there is a theoretical Time-of-Check to Time-of-Use window.
- **Context**: In a forensic acquisition scenario where the investigator has exclusive access to the machine, this risk is negligible.
