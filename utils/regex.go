package utils

import "regexp"

func RegexpSemanticVersion(value *string) bool {
	res, _ := regexp.MatchString(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`, *value)
	return res
}

func RegexpIsNumber(value *string) bool {
	res, _ := regexp.MatchString(`^\d+(?:[.]\d+)?$`, *value)
	return res
}
