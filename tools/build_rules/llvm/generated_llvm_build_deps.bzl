LLVM_BUILD_DEPS = {
    "AArch64AsmParser": [
        "AArch64Desc",
        "AArch64Info",
        "AArch64Utils",
        "MC",
        "MCParser",
        "Support",
    ],
    "AArch64CodeGen": [
        "AArch64Desc",
        "AArch64Info",
        "AArch64Utils",
        "Analysis",
        "AsmPrinter",
        "CodeGen",
        "Core",
        "MC",
        "Scalar",
        "SelectionDAG",
        "Support",
        "Target",
        "GlobalISel",
    ],
    "AArch64Desc": ["AArch64Info", "AArch64Utils", "MC", "Support"],
    "AArch64Disassembler": [
        "AArch64Desc",
        "AArch64Info",
        "AArch64Utils",
        "MC",
        "MCDisassembler",
        "Support",
    ],
    "AArch64Info": ["Support"],
    "AArch64Utils": ["Support"],
    "AMDGPUAsmParser": [
        "MC",
        "MCParser",
        "AMDGPUDesc",
        "AMDGPUInfo",
        "AMDGPUUtils",
        "Support",
    ],
    "AMDGPUCodeGen": [
        "Analysis",
        "AsmPrinter",
        "CodeGen",
        "Core",
        "IPO",
        "MC",
        "AMDGPUDesc",
        "AMDGPUInfo",
        "AMDGPUUtils",
        "Scalar",
        "SelectionDAG",
        "Support",
        "Target",
        "TransformUtils",
        "Vectorize",
        "GlobalISel",
        "BinaryFormat",
        "MIRParser",
    ],
    "AMDGPUDesc": [
        "Core",
        "MC",
        "AMDGPUInfo",
        "AMDGPUUtils",
        "Support",
        "BinaryFormat",
    ],
    "AMDGPUDisassembler": [
        "AMDGPUDesc",
        "AMDGPUInfo",
        "AMDGPUUtils",
        "MC",
        "MCDisassembler",
        "Support",
    ],
    "AMDGPUInfo": ["Support"],
    "AMDGPUUtils": ["Core", "MC", "BinaryFormat", "Support"],
    "ARCCodeGen": [
        "Analysis",
        "AsmPrinter",
        "CodeGen",
        "Core",
        "MC",
        "SelectionDAG",
        "Support",
        "Target",
        "TransformUtils",
        "ARCDesc",
        "ARCInfo",
    ],
    "ARCDesc": ["MC", "Support", "ARCInfo"],
    "ARCDisassembler": ["MCDisassembler", "Support", "ARCInfo"],
    "ARCInfo": ["Support"],
    "ARMAsmParser": [
        "ARMDesc",
        "ARMInfo",
        "MC",
        "MCParser",
        "Support",
        "ARMUtils",
    ],
    "ARMCodeGen": [
        "ARMDesc",
        "ARMInfo",
        "Analysis",
        "AsmPrinter",
        "CodeGen",
        "Core",
        "MC",
        "Scalar",
        "SelectionDAG",
        "Support",
        "Target",
        "GlobalISel",
        "ARMUtils",
        "TransformUtils",
    ],
    "ARMDesc": ["ARMInfo", "ARMUtils", "MC", "MCDisassembler", "Support"],
    "ARMDisassembler": [
        "ARMDesc",
        "ARMInfo",
        "MCDisassembler",
        "Support",
        "ARMUtils",
    ],
    "ARMInfo": ["Support"],
    "ARMUtils": ["Support"],
    "AVRAsmParser": ["MC", "MCParser", "AVRDesc", "AVRInfo", "Support"],
    "AVRCodeGen": [
        "AsmPrinter",
        "CodeGen",
        "Core",
        "MC",
        "AVRDesc",
        "AVRInfo",
        "SelectionDAG",
        "Support",
        "Target",
    ],
    "AVRDesc": ["MC", "AVRInfo", "Support"],
    "AVRDisassembler": ["MCDisassembler", "AVRInfo", "Support"],
    "AVRInfo": ["MC", "Support"],
    "AggressiveInstCombine": ["Analysis", "Core", "Support", "TransformUtils"],
    "Analysis": ["BinaryFormat", "Core", "Object", "ProfileData", "Support"],
    "AsmParser": ["BinaryFormat", "Core", "Support"],
    "AsmPrinter": [
        "Analysis",
        "BinaryFormat",
        "CodeGen",
        "Core",
        "DebugInfoCodeView",
        "DebugInfoDWARF",
        "DebugInfoMSF",
        "MC",
        "MCParser",
        "Remarks",
        "Support",
        "Target",
    ],
    "BPFAsmParser": ["MC", "MCParser", "BPFDesc", "BPFInfo", "Support"],
    "BPFCodeGen": [
        "AsmPrinter",
        "CodeGen",
        "Core",
        "MC",
        "BPFDesc",
        "BPFInfo",
        "SelectionDAG",
        "Support",
        "Target",
    ],
    "BPFDesc": ["MC", "BPFInfo", "Support"],
    "BPFDisassembler": ["MCDisassembler", "BPFInfo", "Support"],
    "BPFInfo": ["Support"],
    "BinaryFormat": ["Support"],
    "BitReader": ["Core", "Support"],
    "BitWriter": ["Analysis", "Core", "MC", "Object", "Support"],
    "CFIVerify": [
        "DebugInfoDWARF",
        "MC",
        "MCDisassembler",
        "MCParser",
        "Support",
        "Symbolize",
    ],
    "CodeGen": [
        "Analysis",
        "BitReader",
        "BitWriter",
        "Core",
        "MC",
        "ProfileData",
        "Scalar",
        "Support",
        "Target",
        "TransformUtils",
    ],
    "Core": ["BinaryFormat", "Remarks", "Support"],
    "Coroutines": [
        "Analysis",
        "Core",
        "IPO",
        "Scalar",
        "Support",
        "TransformUtils",
    ],
    "Coverage": ["Core", "Object", "ProfileData", "Support"],
    "DebugInfoCodeView": ["Support", "DebugInfoMSF"],
    "DebugInfoDWARF": ["BinaryFormat", "Object", "MC", "Support"],
    "DebugInfoMSF": ["Support"],
    "DebugInfoPDB": ["Object", "Support", "DebugInfoCodeView", "DebugInfoMSF"],
    "DlltoolDriver": ["Object", "Option", "Support"],
    "ExecutionEngine": [
        "Core",
        "MC",
        "Object",
        "RuntimeDyld",
        "Support",
        "Target",
    ],
    "Exegesis": [
        "CodeGen",
        "ExecutionEngine",
        "MC",
        "MCDisassembler",
        "MCJIT",
        "Object",
        "ObjectYAML",
        "Support",
    ],
    "ExegesisAArch64": ["AArch64"],
    "ExegesisPowerPC": ["PowerPC"],
    "ExegesisX86": ["X86"],
    "FuzzMutate": [
        "Analysis",
        "BitReader",
        "BitWriter",
        "Core",
        "Scalar",
        "Support",
        "Target",
    ],
    "GlobalISel": [
        "Analysis",
        "CodeGen",
        "Core",
        "MC",
        "Support",
        "Target",
        "TransformUtils",
    ],
    "HexagonAsmParser": [
        "MC",
        "MCParser",
        "Support",
        "HexagonDesc",
        "HexagonInfo",
    ],
    "HexagonCodeGen": [
        "Analysis",
        "AsmPrinter",
        "CodeGen",
        "Core",
        "HexagonAsmParser",
        "HexagonDesc",
        "HexagonInfo",
        "IPO",
        "MC",
        "Scalar",
        "SelectionDAG",
        "Support",
        "Target",
        "TransformUtils",
    ],
    "HexagonDesc": ["HexagonInfo", "MC", "Support"],
    "HexagonDisassembler": [
        "HexagonDesc",
        "HexagonInfo",
        "MC",
        "MCDisassembler",
        "Support",
    ],
    "HexagonInfo": ["Support"],
    "IPO": [
        "AggressiveInstCombine",
        "Analysis",
        "BitReader",
        "BitWriter",
        "Core",
        "InstCombine",
        "IRReader",
        "Linker",
        "Object",
        "ProfileData",
        "Scalar",
        "Support",
        "TransformUtils",
        "Vectorize",
        "Instrumentation",
    ],
    "IRReader": ["AsmParser", "BitReader", "Core", "Support"],
    "InstCombine": ["Analysis", "Core", "Support", "TransformUtils"],
    "Instrumentation": [
        "Analysis",
        "Core",
        "MC",
        "Support",
        "TransformUtils",
        "ProfileData",
    ],
    "Interpreter": ["CodeGen", "Core", "ExecutionEngine", "Support"],
    "JITLink": ["BinaryFormat", "Object", "Support"],
    "LTO": [
        "AggressiveInstCombine",
        "Analysis",
        "BitReader",
        "BitWriter",
        "CodeGen",
        "Core",
        "IPO",
        "InstCombine",
        "Linker",
        "MC",
        "ObjCARC",
        "Object",
        "Passes",
        "Remarks",
        "Scalar",
        "Support",
        "Target",
        "TransformUtils",
    ],
    "LanaiAsmParser": ["MC", "MCParser", "Support", "LanaiDesc", "LanaiInfo"],
    "LanaiCodeGen": [
        "Analysis",
        "AsmPrinter",
        "CodeGen",
        "Core",
        "LanaiAsmParser",
        "LanaiDesc",
        "LanaiInfo",
        "MC",
        "SelectionDAG",
        "Support",
        "Target",
        "TransformUtils",
    ],
    "LanaiDesc": ["LanaiInfo", "MC", "MCDisassembler", "Support"],
    "LanaiDisassembler": [
        "LanaiDesc",
        "LanaiInfo",
        "MC",
        "MCDisassembler",
        "Support",
    ],
    "LanaiInfo": ["Support"],
    "LibDriver": ["BinaryFormat", "BitReader", "Object", "Option", "Support"],
    "LineEditor": ["Support"],
    "Linker": ["Core", "Support", "TransformUtils"],
    "MC": ["Support", "BinaryFormat", "DebugInfoCodeView"],
    "MCA": ["MC", "Support"],
    "MCDisassembler": ["MC", "Support"],
    "MCJIT": [
        "Core",
        "ExecutionEngine",
        "Object",
        "RuntimeDyld",
        "Support",
        "Target",
    ],
    "MCParser": ["MC", "Support"],
    "MIRParser": [
        "AsmParser",
        "BinaryFormat",
        "CodeGen",
        "Core",
        "MC",
        "Support",
        "Target",
    ],
    "MSP430AsmParser": ["MC", "MCParser", "MSP430Desc", "MSP430Info", "Support"],
    "MSP430CodeGen": [
        "AsmPrinter",
        "CodeGen",
        "Core",
        "MC",
        "MSP430Desc",
        "MSP430Info",
        "SelectionDAG",
        "Support",
        "Target",
    ],
    "MSP430Desc": ["MC", "MSP430Info", "Support"],
    "MSP430Disassembler": ["MCDisassembler", "MSP430Info", "Support"],
    "MSP430Info": ["Support"],
    "MipsAsmParser": ["MC", "MCParser", "MipsDesc", "MipsInfo", "Support"],
    "MipsCodeGen": [
        "Analysis",
        "AsmPrinter",
        "CodeGen",
        "Core",
        "MC",
        "MipsDesc",
        "MipsInfo",
        "SelectionDAG",
        "Support",
        "Target",
        "GlobalISel",
    ],
    "MipsDesc": ["MC", "MipsInfo", "Support"],
    "MipsDisassembler": ["MCDisassembler", "MipsInfo", "Support"],
    "MipsInfo": ["Support"],
    "NVPTXCodeGen": [
        "Analysis",
        "AsmPrinter",
        "CodeGen",
        "Core",
        "IPO",
        "MC",
        "NVPTXDesc",
        "NVPTXInfo",
        "Scalar",
        "SelectionDAG",
        "Support",
        "Target",
        "TransformUtils",
        "Vectorize",
    ],
    "NVPTXDesc": ["MC", "NVPTXInfo", "Support"],
    "NVPTXInfo": ["Support"],
    "ObjCARC": ["Analysis", "Core", "Support", "TransformUtils"],
    "Object": ["BitReader", "Core", "MC", "BinaryFormat", "MCParser", "Support"],
    "ObjectYAML": ["Object", "Support", "DebugInfoCodeView"],
    "Option": ["Support"],
    "OrcJIT": [
        "Core",
        "ExecutionEngine",
        "JITLink",
        "Object",
        "MC",
        "RuntimeDyld",
        "Support",
        "Target",
        "TransformUtils",
    ],
    "Passes": [
        "AggressiveInstCombine",
        "Analysis",
        "CodeGen",
        "Core",
        "IPO",
        "InstCombine",
        "Scalar",
        "Support",
        "Target",
        "TransformUtils",
        "Vectorize",
        "Instrumentation",
    ],
    "PowerPCAsmParser": [
        "MC",
        "MCParser",
        "PowerPCDesc",
        "PowerPCInfo",
        "Support",
    ],
    "PowerPCCodeGen": [
        "Analysis",
        "AsmPrinter",
        "CodeGen",
        "Core",
        "MC",
        "PowerPCDesc",
        "PowerPCInfo",
        "Scalar",
        "SelectionDAG",
        "Support",
        "Target",
        "TransformUtils",
    ],
    "PowerPCDesc": ["MC", "PowerPCInfo", "Support"],
    "PowerPCDisassembler": ["MCDisassembler", "PowerPCInfo", "Support"],
    "PowerPCInfo": ["Support"],
    "ProfileData": ["Core", "Support"],
    "RISCVAsmParser": [
        "MC",
        "MCParser",
        "RISCVDesc",
        "RISCVInfo",
        "RISCVUtils",
        "Support",
    ],
    "RISCVCodeGen": [
        "AsmPrinter",
        "Core",
        "CodeGen",
        "MC",
        "RISCVDesc",
        "RISCVInfo",
        "RISCVUtils",
        "SelectionDAG",
        "Support",
        "Target",
    ],
    "RISCVDesc": ["MC", "RISCVInfo", "RISCVUtils", "Support"],
    "RISCVDisassembler": ["MCDisassembler", "RISCVInfo", "Support"],
    "RISCVInfo": ["Support"],
    "RISCVUtils": ["Support"],
    "Remarks": ["Support"],
    "RuntimeDyld": ["MC", "Object", "Support"],
    "Scalar": [
        "AggressiveInstCombine",
        "Analysis",
        "Core",
        "InstCombine",
        "Support",
        "TransformUtils",
    ],
    "SelectionDAG": [
        "Analysis",
        "CodeGen",
        "Core",
        "MC",
        "Support",
        "Target",
        "TransformUtils",
    ],
    "SparcAsmParser": ["MC", "MCParser", "SparcDesc", "SparcInfo", "Support"],
    "SparcCodeGen": [
        "AsmPrinter",
        "CodeGen",
        "Core",
        "MC",
        "SelectionDAG",
        "SparcDesc",
        "SparcInfo",
        "Support",
        "Target",
    ],
    "SparcDesc": ["MC", "SparcInfo", "Support"],
    "SparcDisassembler": ["MCDisassembler", "SparcInfo", "Support"],
    "SparcInfo": ["Support"],
    "Support": ["Demangle"],
    "Symbolize": [
        "DebugInfoDWARF",
        "DebugInfoPDB",
        "Object",
        "Support",
        "Demangle",
    ],
    "SystemZAsmParser": [
        "MC",
        "MCParser",
        "Support",
        "SystemZDesc",
        "SystemZInfo",
    ],
    "SystemZCodeGen": [
        "Analysis",
        "AsmPrinter",
        "CodeGen",
        "Core",
        "MC",
        "Scalar",
        "SelectionDAG",
        "Support",
        "SystemZDesc",
        "SystemZInfo",
        "Target",
    ],
    "SystemZDesc": ["MC", "Support", "SystemZInfo"],
    "SystemZDisassembler": [
        "MC",
        "MCDisassembler",
        "Support",
        "SystemZDesc",
        "SystemZInfo",
    ],
    "SystemZInfo": ["Support"],
    "TableGen": ["Support"],
    "Target": ["Analysis", "Core", "MC", "Support"],
    "TestingSupport": ["Support"],
    "TextAPI": ["Support", "BinaryFormat"],
    "TransformUtils": ["Analysis", "Core", "Support"],
    "Vectorize": ["Analysis", "Core", "Support", "TransformUtils"],
    "WebAssemblyAsmParser": ["MC", "MCParser", "WebAssemblyInfo", "Support"],
    "WebAssemblyCodeGen": [
        "Analysis",
        "AsmPrinter",
        "BinaryFormat",
        "CodeGen",
        "Core",
        "MC",
        "Scalar",
        "SelectionDAG",
        "Support",
        "Target",
        "TransformUtils",
        "WebAssemblyDesc",
        "WebAssemblyInfo",
    ],
    "WebAssemblyDesc": ["MC", "Support", "WebAssemblyInfo"],
    "WebAssemblyDisassembler": [
        "WebAssemblyDesc",
        "MCDisassembler",
        "WebAssemblyInfo",
        "Support",
    ],
    "WebAssemblyInfo": ["Support"],
    "WindowsManifest": ["Support"],
    "X86AsmParser": ["MC", "MCParser", "Support", "X86Desc", "X86Info"],
    "X86CodeGen": [
        "Analysis",
        "AsmPrinter",
        "CodeGen",
        "Core",
        "MC",
        "SelectionDAG",
        "Support",
        "Target",
        "X86Desc",
        "X86Info",
        "X86Utils",
        "GlobalISel",
        "ProfileData",
    ],
    "X86Desc": [
        "MC",
        "MCDisassembler",
        "Object",
        "Support",
        "X86Info",
        "X86Utils",
    ],
    "X86Disassembler": ["MCDisassembler", "Support", "X86Info"],
    "X86Info": ["Support"],
    "X86Utils": ["Support"],
    "XCoreCodeGen": [
        "Analysis",
        "AsmPrinter",
        "CodeGen",
        "Core",
        "MC",
        "SelectionDAG",
        "Support",
        "Target",
        "TransformUtils",
        "XCoreDesc",
        "XCoreInfo",
    ],
    "XCoreDesc": ["MC", "Support", "XCoreInfo"],
    "XCoreDisassembler": ["MCDisassembler", "Support", "XCoreInfo"],
    "XCoreInfo": ["Support"],
    "XRay": ["Support", "Object"],
    "gtest": ["Support"],
    "gtest_main": ["gtest"],
}
