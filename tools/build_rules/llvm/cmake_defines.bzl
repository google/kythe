LLVM_TARGETS = ["X86", "PowerPC", "ARM", "AArch64", "Mips"]
_CMAKE_DEFINES = {
    "BACKTRACE_HEADER": "execinfo.h",
    "BUG_REPORT_URL": "https://bugs.llvm.org/",
    "ENABLE_BACKTRACES": "1",
    "ENABLE_CRASH_OVERRIDES": "1",
    "HAVE_BACKTRACE": "TRUE",
    "HAVE_DECL_FE_ALL_EXCEPT": "1",
    "HAVE_DECL_FE_INEXACT": "1",
    "HAVE_DLFCN_H": "1",
    "HAVE_DLOPEN": "1",
    "HAVE_ERRNO_H": "1",
    "HAVE_FCNTL_H": "1",
    "HAVE_FENV_H": "1",
    "HAVE_FUTIMENS": "1",
    "HAVE_FUTIMES": "1",
    "HAVE_GETPAGESIZE": "1",
    "HAVE_GETRLIMIT": "1",
    "HAVE_GETRUSAGE": "1",
    "HAVE_ISATTY": "1",
    "HAVE_LIBPTHREAD": "1",
    "HAVE_LIBZ": "1",
    "HAVE_POSIX_SPAWN": "1",
    "HAVE_PREAD": "1",
    "HAVE_PTHREAD_GETNAME_NP": "1",
    "HAVE_PTHREAD_GETSPECIFIC": "1",
    "HAVE_PTHREAD_H": "1",
    "HAVE_PTHREAD_MUTEX_LOCK": "1",
    "HAVE_PTHREAD_RWLOCK_INIT": "1",
    "HAVE_PTHREAD_SETNAME_NP": "1",
    "HAVE_REALPATH": "1",
    "HAVE_SBRK": "1",
    "HAVE_SETENV": "1",
    "HAVE_SETRLIMIT": "1",
    "HAVE_SIGNAL_H": "1",
    "HAVE_STRERROR": "1",
    "HAVE_STRERROR_R": "1",
    "HAVE_SYSCONF": "1",
    "HAVE_SYS_IOCTL_H": "1",
    "HAVE_SYS_MMAN_H": "1",
    "HAVE_SYS_PARAM_H": "1",
    "HAVE_SYS_RESOURCE_H": "1",
    "HAVE_SYS_STAT_H": "1",
    "HAVE_SYS_TIME_H": "1",
    "HAVE_SYS_TYPES_H": "1",
    "HAVE_TERMINFO": "1",
    "HAVE_TERMIOS_H": "1",
    "HAVE_UNISTD_H": "1",
    "HAVE__UNWIND_BACKTRACE": "1",
    "HAVE_ZLIB_H": "1",
    "LLVM_ENABLE_THREADS": "1",
    "LLVM_ENABLE_ZLIB": "1",
    "LLVM_HAS_ATOMICS": "1",
    "LLVM_NATIVE_ARCH": "X86",
    "LLVM_ON_UNIX": "1",
    "LLVM_VERSION_MAJOR": "8",
    "LLVM_VERSION_MINOR": "0",
    "LLVM_VERSION_PATCH": "0",
    "LLVM_VERSION_PRINTER_SHOW_HOST_TARGET_INFO": "1",
    "LLVM_VERSION_STRING": "8.0.0svn",
    "CLANG_VERSION": "8.0.0",
    "CLANG_VERSION_MAJOR": "8",
    "CLANG_VERSION_MINOR": "0",
    "CLANG_VERSION_PATCHLEVEL": "0",
    "CLANG_VERSION_STRING": "8.0.0svn",
    "PACKAGE_BUGREPORT": "https://bugs.llvm.org/",
    "PACKAGE_NAME": "LLVM",
    "PACKAGE_STRING": "LLVM 8.0.0svn",
    "PACKAGE_VERSION": "8.0.0svn",
    "RETSIGTYPE": "void",
    "LLVM_ENUM_TARGETS": "".join(["LLVM_TARGET(%s)\n" % t for t in LLVM_TARGETS]),
    "LLVM_ENUM_ASM_PRINTERS": "".join(["LLVM_ASM_PRINTER(%s)\n" % t for t in LLVM_TARGETS]),
    "LLVM_ENUM_ASM_PARSERS": "".join(["LLVM_ASM_PARSER(%s)\n" % t for t in LLVM_TARGETS]),
    "LLVM_ENUM_DISASSEMBLERS": "".join(["LLVM_DISASSEMBLER(%s)\n" % t for t in LLVM_TARGETS]),
}

def cmake_defines():
    return struct(
        default = _CMAKE_DEFINES + {
            "HAVE_MALLOC_H": "1",
            "HAVE_MALLINFO": "1",
            "HAVE_LINK_H": "1",
            "HAVE_LSEEK64": "1",
            "HAVE_POSIX_FALLOCATE": "1",
            "HAVE_SCHED_GETAFFINITY": "1",
            "HAVE_CPU_COUNT": "1",
            "HAVE_SIGALTSTACK": "1",
            "LLVM_DEFAULT_TARGET_TRIPLE": "x86_64-unknown-linux-gnu",
            "LLVM_HOST_TRIPLE": "x86_64-unknown-linux-gnu",
            "LTDL_SHLIB_EXT": ".so",
        },
        darwin = _CMAKE_DEFINES + {
            "HAVE_DECL_ARC4RANDOM": "1",
            "HAVE_MALLOC_MALLOC_H": "1",
            "HAVE_MALLOC_ZONE_STATISTICS": "1",
            "HAVE_DLADDR": "1",
            "HAVE_MACH_MACH_H": "1",
            "HAVE_LIBXAR": "1",
            "LLVM_DEFAULT_TARGET_TRIPLE": "x86_64-apple-darwin17.7.0",
            "LLVM_HOST_TRIPLE": "x86_64-apple-darwin17.7.0",
            "LTDL_SHLIB_EXT": ".dylib",
        },
    )
