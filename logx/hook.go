package logx

type Hook interface {
	SetLogger(logger *LogX)
	Fire(entry *Entry) error
	Levels() []Level
}

func (l *LogX) AddHook(hook Hook) {
	l.mu.Lock()
	defer l.mu.Unlock()

	hook.SetLogger(l)

	for _, level := range hook.Levels() {
		if l.hooks[level] == nil {
			l.hooks[level] = []Hook{}
		}
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
