package kube

type ActionType int

// TODO: use subcommands instead
const (
	ActionMerge ActionType = iota
	ActionPurge
	ActionPrint
)

func getAction(isPurge, printForward bool) ActionType {
	if printForward {
		return ActionPrint
	}

	if isPurge {
		return ActionPurge
	}

	return ActionMerge
}
