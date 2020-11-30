package dyld

import (
	"encoding/binary"
	"fmt"
	"io"
	"sort"

	"github.com/blacktop/go-arm64"
	"github.com/blacktop/go-macho"
	"github.com/blacktop/ipsw/internal/demangle"
)

// Demangle a string just as the GNU c++filt program does.
func doDemangle(name string) string {
	var deStr string

	if len(name) == 0 {
		return name
	}

	skip := 0
	if name[0] == '.' || name[0] == '$' {
		skip++
	}
	if name[skip] == '_' {
		skip++
	}
	result := demangle.Filter(name[skip:])
	if result == name[skip:] {
		deStr += name
	} else {
		if name[0] == '.' {
			deStr += "."
		}
		deStr += result
	}
	return deStr
}

func (f *File) FunctionSize(starts []uint64, addr uint64) int64 {
	i := sort.Search(len(starts), func(i int) bool { return starts[i] >= addr })
	if i+1 == len(starts) && starts[i] == addr {
		return -1
	} else if i < len(starts) && starts[i] == addr {
		return int64(starts[i+1] - addr)
	}
	return 0
}

// IsFunctionStart checks if address is at a function start and returns symbol name
func (f *File) IsFunctionStart(starts []uint64, addr uint64, demangle bool) (bool, string) {
	if f.FunctionSize(starts, addr) != 0 {
		if symName, ok := f.AddressToSymbol[addr]; ok {
			if demangle {
				return ok, doDemangle(symName)
			}
			return ok, symName
		}
		return true, ""
	}
	return false, ""
}

// FindSymbol returns symbol from the addr2symbol map for a given virtual address
func (f *File) FindSymbol(addr uint64, demangle bool) string {
	if symName, ok := f.AddressToSymbol[addr]; ok {
		if demangle {
			return doDemangle(symName)
		}
		return symName
	}

	return ""
}

// ParseSymbolStubs parse symbol stubs in MachO
func (f *File) ParseSymbolStubs(m *macho.File) error {
	for _, sec := range m.Sections {
		if sec.Flags.IsSymbolStubs() {

			r := io.NewSectionReader(f.r, int64(sec.Offset), int64(sec.Size))

			var prevInstruction arm64.Instruction
			for i := range arm64.Disassemble(r, arm64.Options{StartAddress: int64(sec.Addr)}) {
				// TODO: remove duplicate code (refactor into IL)
				operation := i.Instruction.Operation().String()
				if (operation == "ldr" || operation == "add") && prevInstruction.Operation().String() == "adrp" {
					operands := i.Instruction.Operands()
					if operands != nil && prevInstruction.Operands() != nil {
						adrpRegister := prevInstruction.Operands()[0].Reg[0]
						adrpImm := prevInstruction.Operands()[1].Immediate
						if operation == "ldr" && adrpRegister == operands[1].Reg[0] {
							adrpImm += operands[1].Immediate
						} else if operation == "add" && adrpRegister == operands[0].Reg[0] {
							adrpImm += operands[2].Immediate
						}

						if len(f.AddressToSymbol[adrpImm]) > 0 {
							f.AddressToSymbol[prevInstruction.Address()] = f.AddressToSymbol[adrpImm]
						}
					}
				}
				// fmt.Printf("%#08x:  %s\t%s%s%s\n", i.Instruction.Address(), i.Instruction.OpCodes(), i.Instruction.Operation(), pad(10-len(i.Instruction.Operation().String())), i.Instruction.OpStr())
				prevInstruction = *i.Instruction
			}
		}
	}

	return nil
}

// ParseGOT parse global offset table in MachO
func (f *File) ParseGOT(m *macho.File) error {

	// authPtr := m.Section("__AUTH_CONST", "__auth_ptr")
	// data, err := authPtr.Data()
	// if err != nil {
	// 	return err
	// }
	// ptrs := make([]uint64, authPtr.Size/8)
	// if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &ptrs); err != nil {
	// 	return err
	// }
	// for _, ptr := range ptrs {
	// 	newPtr := convertToVMAddr(m, ptr)
	// 	fmt.Printf("ptr: %#x\n", ptr)
	// 	fmt.Printf("newPtr: %#x, %s\n", newPtr, symbolMap[newPtr])
	// }
	for _, sec := range m.Sections {
		if sec.Flags.IsNonLazySymbolPointers() {

			r := io.NewSectionReader(f.r, int64(sec.Offset), int64(sec.Size))

			ptrs := make([]uint64, sec.Size/8)

			if err := binary.Read(r, binary.LittleEndian, &ptrs); err != nil {
				return err
			}
			// imports, err := m.ImportedSymbolNames()
			// if err != nil {
			// 	return err
			// }
			// for name := range imports {
			// 	fmt.Println(name)
			// }
			for idx, ptr := range ptrs {
				gotPtr := sec.Addr + uint64(idx*8)
				// fmt.Printf("gotPtr: %#x\n", gotPtr)
				var targetValue uint64
				pointer := CacheSlidePointer3(ptr)
				if pointer.Authenticated() {
					targetValue = 0x180000000 + pointer.OffsetFromSharedCacheBase()
				} else {
					targetValue = pointer.SignExtend51()
				}
				// fmt.Printf("ptr: %#x\n", ptr)
				// fmt.Printf("newPtr: %#x, %s\n", targetValue, symbolMap[targetValue])
				// fmt.Println(lookupSymbol(m, targetValue))
				if _, ok := f.AddressToSymbol[gotPtr]; ok {
					// continue
					f.AddressToSymbol[gotPtr] = "__got." + f.AddressToSymbol[gotPtr]
				} else {
					if _, ok := f.AddressToSymbol[targetValue]; ok {
						f.AddressToSymbol[gotPtr] = "__got." + f.AddressToSymbol[targetValue]
					} else {
						f.AddressToSymbol[gotPtr] = fmt.Sprintf("__got_ptr_%#x", targetValue)
					}
				}
			}
		}
	}

	return nil
}
