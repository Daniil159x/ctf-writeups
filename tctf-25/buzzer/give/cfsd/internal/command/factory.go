package command

import (
	"os/exec"
	"path"
)

const (
	ChrootPath  = "/usr/sbin/chroot"
	SetPrivPath = "/bin/setpriv"

	UserName = "u214528"
	UserId   = "1001"
	GroupId  = "1001"

	CatPath = "/bin/cat"

	FindPath     = "/bin/find"
	UserFilesDir = "/files"

	PosixShell = "/bin/sh"
)

func chrootDir() string {
	return path.Join("/home", UserName)
}

func chrootCmd() []string {
	// chroot /home/<user> /bin/setpriv --nnp --regid <gid> --reuid <uid> --clear-groups --reset-env
	return []string{
		ChrootPath,
		chrootDir(),
		SetPrivPath,
		"--nnp",
		"--regid", GroupId,
		"--reuid", UserId,
		"--clear-groups",
		"--reset-env",
	}
}

func MakeListFiles() *exec.Cmd {
	line := chrootCmd()
	line = append(line, FindPath, UserFilesDir, "-type", "f")
	return exec.Command(line[0], line[1:]...)
}

func MakeOpenFiles(chrootedPath string) *exec.Cmd {
	line := chrootCmd()
	line = append(line, CatPath, chrootedPath)
	return exec.Command(line[0], line[1:]...)
}

func MakeShell() *exec.Cmd {
	line := chrootCmd()
	line = append(line, PosixShell)
	return exec.Command(line[0], line[1:]...)
}
