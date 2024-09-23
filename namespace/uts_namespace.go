package namespace

import (
	"os"
	"os/exec"
	"syscall"
)

//ns api:
// clone()创建新进程包含到ns
//unshare()将进程移出某个ns
//share()将进程加入ns

func Uts() {
	//指定fork出来新进程的初始命令
	cmd := exec.Command("sh")

	/*
		CLONE_NEWPID: pstree -pl查看进程pid
		CLONE_NEWNS: mount -t proc proc /proc挂载到ns里,ps -ef少了很多进程
		CLONE_NEWUSER: id看用户id
		CLONE_NEWNET: ifconfig看网卡
	*/
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
	}
	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: uint32(1),
		Gid: uint32(1),
	}
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run() //使用clone()创建新进程,这个新进程的子进程也会被包含到这个namespace里
	if err != nil {
		return
	}
}
