package main

type ReplyManager struct {
	Register []int
}

func (rm *ReplyManager) addRegister(messageID int) {
	rm.Register = append(rm.Register, messageID)
}

func (rm *ReplyManager) contains(slice []int, val int) (bool, []int) {
	for i, item := range slice {
		if item == val {
			slice = append(slice[:i], slice[i+1:]...)
			return true, slice
		}
	}
	return false, slice
}

func (rm *ReplyManager) containsRegister(messageID int) bool {
	var contains bool
	contains, rm.Register = rm.contains(rm.Register, messageID)
	return contains
}
