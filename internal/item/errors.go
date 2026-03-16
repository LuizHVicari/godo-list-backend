package item

import "errors"

var ErrorItemNotFound = errors.New("item not found")
var ErrorInvalidFilterParams = errors.New("invalid filter parameters")
var ErrorItemNotBelongsToStep = errors.New("item does not belong to step")
var ErrorItemPositionTaken = errors.New("position already taken in this step")
