package bashColor

import "fmt"

type ColorCode string

const HEADER ColorCode = "\033[95m"
const OKBLUE ColorCode = "\033[94m"
const OKCYAN ColorCode = "\033[96m"
const OKGREEN ColorCode = "\033[92m"
const WARNING ColorCode = "\033[93m"
const FAIL ColorCode = "\033[91m"
const ENDC ColorCode = "\033[0m"
const BOLD ColorCode = "\033[1m"
const UNDERLINE ColorCode = "\033[4m"

/* ColorString returns a string with the color assigned at the front and ENDC at the end.
 */
func ColorString(color ColorCode, msg string) string {
	return ColorStringF(color, msg)
}

/*
	ColorStringF returns a string with the color assigned at the front and ENDC at the end.

It takes verbs and args as other F designated string functions
*/
func ColorStringF(color ColorCode, msg string, args ...interface{}) string {
	combinedMessage := fmt.Sprintf(msg, args...)

	return string(color + ColorCode(combinedMessage) + ENDC)
}
