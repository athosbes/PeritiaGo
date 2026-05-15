package verify

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modwintrust        = windows.NewLazySystemDLL("wintrust.dll")
	procWinVerifyTrust = modwintrust.NewProc("WinVerifyTrust")
)

const (
	WTD_UI_NONE            = 2
	WTD_REVOKE_NONE        = 0
	WTD_CHOICE_FILE        = 1
	WTD_STATEACTION_IGNORE = 0
	WTD_SAFER_FLAG         = 0x100
)

var (
	// WINTRUST_ACTION_GENERIC_VERIFY_V2
	WINTRUST_ACTION_GENERIC_VERIFY_V2 = windows.GUID{
		Data1: 0xaac56b,
		Data2: 0xcd44,
		Data3: 0x11d0,
		Data4: [8]byte{0x8c, 0xc2, 0x0, 0xc0, 0x4f, 0xc2, 0x95, 0xee},
	}
)

type WINTRUST_FILE_INFO struct {
	Size     uint32
	FilePath *uint16
	File     windows.Handle
	PgPolicy windows.Handle
}

type WINTRUST_DATA struct {
	Size               uint32
	PolicyCallbackData windows.Handle
	SIPClientData      windows.Handle
	UIChoice           uint32
	RevocationChecks   uint32
	UnionChoice        uint32
	File               unsafe.Pointer
	StateAction        uint32
	StateData          windows.Handle
	URLReference       *uint16
	ProvFlags          uint32
	UIContext          uint32
	SignatureSettings  unsafe.Pointer
}

// IsWindowsSigned checks if a file has a valid digital signature via Windows WinVerifyTrust.
func IsWindowsSigned(path string) error {
	pathPtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	fileInfo := WINTRUST_FILE_INFO{
		Size:     uint32(unsafe.Sizeof(WINTRUST_FILE_INFO{})),
		FilePath: pathPtr,
	}

	wintrustData := WINTRUST_DATA{
		Size:             uint32(unsafe.Sizeof(WINTRUST_DATA{})),
		UIChoice:         WTD_UI_NONE,
		RevocationChecks: WTD_REVOKE_NONE,
		UnionChoice:      WTD_CHOICE_FILE,
		File:             unsafe.Pointer(&fileInfo),
		StateAction:      WTD_STATEACTION_IGNORE,
		ProvFlags:        WTD_SAFER_FLAG,
	}

	r1, _, _ := procWinVerifyTrust.Call(
		0,
		uintptr(unsafe.Pointer(&WINTRUST_ACTION_GENERIC_VERIFY_V2)),
		uintptr(unsafe.Pointer(&wintrustData)),
	)

	if uint32(r1) != 0 {
		return fmt.Errorf("WinVerifyTrust failed with error: 0x%X (file might be unsigned or signature invalid)", uint32(r1))
	}

	return nil
}

// IsSigned checks if a file has a valid digital signature.
// It prioritizes Sigstore verification if a .sigstore.json bundle is present,
// otherwise fallbacks to Windows Authenticode verification.
func IsSigned(path string) error {
	// Try Sigstore first
	err := IsSigstoreSigned(path)
	if err == nil {
		return nil
	}

	// If Sigstore failed because the bundle was missing, fallback to Windows
	// Otherwise, report the Sigstore error
	bundlePath := path + ".sigstore.json"
	if _, statErr := os.Stat(bundlePath); os.IsNotExist(statErr) {
		return IsWindowsSigned(path)
	}

	return fmt.Errorf("sigstore verification failed: %v (note: windows signature check skipped)", err)
}
