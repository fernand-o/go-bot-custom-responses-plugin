package customresponses

import (
	"fmt"
	"strings"
)

var matchExample = strings.Join([]string{
	"!responses match set \"Is someone there?\" \"Hello\"",
	"!responses match set \"Rick Sanchez\" \"@rick\" #fun",
	"!responses match showall",
	"!responses match unset 1",
}, "\n")
var argumentsMatchExample = formatMessage(matchExample)

var listExample = strings.Join([]string{
	"!responses list add #fun \"Hey %s, looks like you're guilty\"",
	"!responses list add #fun \"Is that you %s?\"",
	"!responses list show #fun",
	"!responses list showall",
	"!responses list remove #fun \"Is that you %s?\"",
	"!responses list delete #fun",
}, "\n")
var argumentsListExample = formatMessage(listExample)

var argumentsGeneralExample = strings.Join([]string{
	"```",
	"Defining matches and responses:",
	matchExample,
	"",
	"Defining lists to be used within the responses:",
	listExample,
	"```",
}, "\n")

func formatMessage(msg string) string {
	return "```" + msg + "```"
}

func matchesKeyFmt(sufix string) string {
	return "matches:" + sufix
}

func userMessageSetResponse(match, response string) string {
	return fmt.Sprintf("Ok! I will send a message with `%s` when i found any occurences of `%s`", response, match)
}

func userMessageUnsetResponse(match string) string {
	return fmt.Sprintf("Done, i'll not say anything more related to `%s`", match)
}

func userMessageNoResposesDefined() string {
	return fmt.Sprintf("There are no responses defined yet. \n %s", argumentsMatchExample)
}

func userMessageResponsesDeleted() string {
	return "All responses were deleted."
}

func userMessageListMessageAdded(list, message string) string {
	return fmt.Sprintf("The message `%s` was added to the list `%s`.", message, list)
}

func userMessageListMessageRemoved(list, message string) string {
	return fmt.Sprintf("The message `%s` was removed of the list `%s`.", message, list)
}

func userMessageListDeleted(list string) string {
	return fmt.Sprintf("The list %s was deleted.", list)
}

func userMessageNoListsDefined() string {
	return fmt.Sprintf("There are no lists defined yet. \n %s", argumentsListExample)
}

func userMessageListInvalidName() string {
	return "The list name must starts with #"
}

func userMessageDBErased() string {
	return "Database erased."
}
