package xmlsec

import (
	"C"
	"fmt"
	"runtime"

	"github.com/crewjam/errset"
)

var globalErrors = map[uintptr]errset.ErrSet{}

type libraryError struct {
	FileName string
	Line     int
	FuncName string
	Object   string
	Subject  string
	Reason   int
	Message  string
}

func (e libraryError) Error() string {
	return fmt.Sprintf(
		"func=%s:file=%s:line=%d:obj=%s:subj=%s:error=%d:%s",
		e.FuncName,
		e.FileName,
		e.Line,
		e.Object,
		e.Subject,
		e.Reason,
		e.Message)
}

//export onError
func onError(file *C.char, line C.int, funcName *C.char, errorObject *C.char, errorSubject *C.char, reason C.int, msg *C.char) {
	err := libraryError{
		FuncName: C.GoString(funcName),
		FileName: C.GoString(file),
		Line:     int(line),
		Object:   C.GoString(errorObject),
		Subject:  C.GoString(errorSubject),
		Reason:   int(reason),
		Message:  C.GoString(msg)}
	threadID := getThreadID()
	globalErrors[threadID] = append(globalErrors[threadID], err)
}

// startProcessingXML is called whenever we enter a function exported by this package.
// It locks the current goroutine to the current thread and establishes a thread-local
// error object. If the library later calls onError then the error will be appended
// to the error object associated with the current thread.
func startProcessingXML() {
	runtime.LockOSThread()
	globalErrors[getThreadID()] = errset.ErrSet{}
}

// stopProcessingXML unlocks the goroutine-thread lock and deletes the current
// error stack.
func stopProcessingXML() {
	runtime.UnlockOSThread()
	delete(globalErrors, getThreadID())
}

// popError returns the global error for the current thread and resets it to
// an empty error. Returns nil if no errors have occurred. This function must be
// called after startProcessingXML() and before stopProcessingXML(). All three
// functions must be called on the same goroutine.
func popError() error {
	threadID := getThreadID()
	rv := globalErrors[threadID].ReturnValue()
	globalErrors[threadID] = errset.ErrSet{}
	return rv
}

// mustPopError is like popError except that if there is no error on the stack
// it returns a generic error.
func mustPopError() error {
	err := popError()
	if err == nil {
		err = fmt.Errorf("libxmlsec: call failed")
	}
	return err
}
