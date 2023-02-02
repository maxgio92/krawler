package kernelrelease

import "regexp"

var kernelVersionPattern = regexp.MustCompile(`(?P<fullversion>^(?P<version>0|[1-9]\d*)\.(?P<patchlevel>0|[1-9]\d*)[.+]?(?P<sublevel>0|[1-9]\d*)?)(?P<fullextraversion>[-.+](?P<extraversion>0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)([\.+~](0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-_]*))*)?(\+[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*)?$`)
