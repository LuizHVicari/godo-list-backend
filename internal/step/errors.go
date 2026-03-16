package step

import "errors"

var ErrorStepNotFound = errors.New("step not found")
var ErrorInvalidFilterParams = errors.New("invalid filter parameters")
var ErrorStepNotBelongsToProject = errors.New("step does not belong to project")
var ErrorStepPositionTaken = errors.New("position already taken in this project")
