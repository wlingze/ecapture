package config

import (
	"fmt"

	"golang.org/x/arch/arm64/arm64asm"
	"golang.org/x/arch/x86/x86asm"
)

func (gc *GoTLSConfig) decodeInstruction(instHex []byte) ([]int, error) {
	switch gc.goElfArch {
	case "amd64":
		return decodeInstructionAMD(instHex)
	case "arm64":
		return decodeInstructionARM(instHex)
	default:
		return nil, fmt.Errorf("unsupport CPU arch :%s", gc.goElfArch)
	}
}

// decodeInstruction Decode into assembly instructions and identify the RET instruction to return the offset.
func decodeInstructionAMD(instHex []byte) ([]int, error) {
	var offsets []int
	for i := 0; i < len(instHex); {
		inst, err := x86asm.Decode(instHex[i:], 64)
		if err != nil {
			return nil, err
		}
		if inst.Op == x86asm.RET {
			offsets = append(offsets, i)
		}
		i += inst.Len
	}
	return offsets, nil
}

const (
	// Arm64armInstSize via :  arm64/arm64asm/decode.go:Decode() size = 4
	Arm64armInstSize = 4
)

// decodeInstruction Decode into assembly instructions and identify the RET instruction to return the offset.
func decodeInstructionARM(instHex []byte) ([]int, error) {
	var offsets []int
	for i := 0; i < len(instHex); {
		inst, _ := arm64asm.Decode(instHex[i:]) // Why ignore error: https://github.com/gojue/ecapture/pull/506
		if inst.Op == arm64asm.RET {
			offsets = append(offsets, i)
		}
		i += Arm64armInstSize
	}
	return offsets, nil
}
