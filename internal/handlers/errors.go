package handlers

import "github.com/pkg/errors"

var OpenFileError = errors.New("error opening file")
var ReadingFileError = errors.New("error reading file")
var JobIDIsNotFound = errors.New("jobID param does not exist for the existing job")
