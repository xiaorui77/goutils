package logx

type Hook interface {
	Fire(entry *Entry) error
	Levels() []Level
}

func (l *LogX) AddHook(hook Hook) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, level := range hook.Levels() {
		l.hooks[level] = append(l.hooks[level], hook)
	}
}

func (l *LogX) fireHooks(level Level, entry *Entry) error {
	for _, hs := range l.hooks[level] {
		if err := hs.Fire(entry); err != nil {
			return err
		}
	}
	return nil
}
