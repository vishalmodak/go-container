package main

import (
    "syscall"
    "fmt"
    "os"
    "os/exec"
)

func main() {
    switch os.Args[1] {
    case "run":
        run()
    case "child":
        child()
    default:
        panic("help")
    }
}

func run() {
    fmt.Printf("Running %v\n", os.Args[2:])
    
    cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.SysProcAttr = &syscall.SysProcAttr {
        //Clone creates a new process
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,   //sets the host namespace (UNIX TIME SHARING)
    }

    must(cmd.Run())
}

func child() {
    fmt.Printf("Running %v \n", os.Args[2:])

    cmd := exec.Command(os.Args[2], os.Args[3:]...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS,
    }
    must(syscall.Sethostname([]byte("container")))
    must(syscall.Chroot("/home/vagrant/ubuntufs"))
    must(os.Chdir("/"))
    must(syscall.Mount("proc", "proc", "proc", 0, ""))
    
    must(cmd.Run())

    must(syscall.Unmount("proc", 0))
}

func must(err error) {
    if err != nil {
        panic(err)
    }
}
