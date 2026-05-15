# PeritiaGo Security Guarantees

PeritiaGo is designed for high-integrity digital forensic acquisitions. The following technical guarantees are enforced by the core architecture:

## 1. Immutability (Forensic Mode)
When `ForensicMode` is enabled, the application enforces a strict **Read-Only** policy for the target system. 
- **Guarantee**: No files will be modified or created outside of the designated output directory.
- **Enforcement**: Centralized `internal/guard` component validates every write attempt via `filepath.EvalSymlinks` and deep path canonicalization.

## 2. Chain of Custody (Integrity)
- **Guarantee**: Every byte collected is accounted for and its integrity is mathematically provable.
- **Enforcement**: 
    - Real-time SHA-256 hashing of every exported artifact.
    - Automated creation of a `manifesto.txt` containing hashes of all collected files.
    - Generation of a **Master Hash** representing the entire acquisition state.

## 3. Execution Audit
- **Guarantee**: All security-critical decisions (allowed/blocked writes) are logged.
- **Enforcement**: The `filesystem` wrappers log each operation with a unique `ExecutionID`, providing a full operational audit trail for the investigator.

## 4. Anti-Bypass Protection
- **Guarantee**: Common path-based bypasses are blocked.
- **Enforcement**: 
    - **Path Traversal**: Blocked via `filepath.Abs` and relative path verification.
    - **Symlinks/Junctions**: Resolved to their final targets before permission checks.
    - **UNC Paths**: Canonicalized and checked against the output directory whitelist.
