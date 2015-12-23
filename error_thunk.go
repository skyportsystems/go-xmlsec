package xmlsec

// onError_cgo is a C function that can be passed to xmlSecErrorsSetCallback which
// in turn invokes the go function onError which captures the error.
// For reasons I do not completely understand, it must be defined in a different
// file from onError

// void onError_cgo(char *file, int line, char *funcName, char *errorObject, char *errorSubject, int reason, char *msg) {
//   onError(file, line, funcName, errorObject, errorSubject, reason, msg);
// }
import "C"
