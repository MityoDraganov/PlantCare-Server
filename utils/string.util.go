package utils

func CoalesceString(s *string) *string {
	if s == nil {
		empty := ""
		return &empty
	}
	return s
}
