package session

func (s *Session) Error(key string) string {
	v, ok := s.errors[key]
	if !ok {
		return ""
	}

	return v
}

func (s *Session) HasError(key string) bool {
	_, ok := s.errors[key]
	return ok
}

func (s *Session) Old(key string) string {
	v, ok := s.olds[key]
	if !ok {
		return ""
	}

	return v
}

func (s *Session) HasOld(key string) bool {
	_, ok := s.olds[key]
	return ok
}

func (s *Session) Errors(key string, val string) error {
	if s.HasError(key) {
		return nil
	}

	s.errors[key] = val
	return nil
}

func (s *Session) Olds(key string, val string) error {
	if s.HasOld(key) {
		return nil
	}

	s.olds[key] = val
	return nil
}
