package handlers

import "github.com/pkg/errors"

var OpenFileError = errors.New("error opening file")
var ReadingFileError = errors.New("error reading file")
var JobIDIsNotFound = errors.New("jobID param does not exist for the existing job")
var StoreIsNotFoundInContext = errors.New("store is not found in context")
var SessionIDDoesNotMatch = errors.New("sessionID does not match the sessionID in the context")
var SessionIsNotFoundInContext = errors.New("session is not found in context")
var FileTypeIsNotSupported = errors.New("file type is not supported")
