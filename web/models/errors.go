package models

import "errors"

var ErrFileNotFound error = errors.New("Unable to find requested file")
var ErrParsingCookies error = errors.New("An error occured while parsing cookies. Please try contacting the staff.")
