package errcode

import "fmt"

const (
	Database_err = 1000001 + iota
	Faild
	ValidationFaild
	NotFoundUser
	TokenIsErr
	FileGetFail
	FileStoreFail
	NotFoundFile
	FileFastUploadFail
	FileMpInitErr
	FileMpCheckFail
	RetryErr
	DownloadFileNotValid
	ListUserFileErr
	DeleteUserFilesErr
	RenameUserFileErr
	GetSessionAllErr
	CreateSessionErr
	GetSessionInfoErr
	UserFileGetFail
	CreateUserFileFail
	AddFriendErr
	GetFriendsErr
)

type withCode struct {
	err   error
	code  int
	cause error
}

func (w *withCode) Error() string { return fmt.Sprintf("%v", w) }

func WithCode(code int, format string, args ...interface{}) error {
	return &withCode{
		err:  fmt.Errorf(format, args...),
		code: code,
	}
}
