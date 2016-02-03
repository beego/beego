package logs

import "time"

type filesLogWriter struct {
	writers [LevelDebug + 1]*fileLogWriter
}

func (f *filesLogWriter) Init(config string) error {
	writer := newFileWriter().(*fileLogWriter)
	err := writer.Init(config)
	if err != nil {
		return err
	}
	f.writers[0] = writer

	for i := LevelEmergency; i <= f.writers[0].Level; i++ {
		writer = newFileWriter().(*fileLogWriter)
		writer.Init(config)
		writer.Level = i
		writer.fileNameOnly += "." + levelNames[i]
		f.writers[i+1] = writer
	}

	return nil
}

func (f *filesLogWriter) Destroy() {
	for i := 0; i < len(f.writers); i++ {
		if f.writers[i] != nil {
			f.writers[i].Destroy()
		}
	}
}

func (f *filesLogWriter) WriteMsg(when time.Time, msg string, level int) error {
	if f.writers[0] != nil {
		f.writers[0].WriteMsg(when, msg, level)
	}
	for i := 1; i < len(f.writers); i++ {
		if f.writers[i] != nil {
			if level == f.writers[i].Level {
				f.writers[i].WriteMsg(when, msg, level)
			}
		}
	}
	return nil
}

func (f *filesLogWriter) Flush() {
	for i := 0; i < len(f.writers); i++ {
		if f.writers[i] != nil {
			f.writers[i].Flush()
		}
	}
}


// newFilesWriter create a FileLogWriter returning as LoggerInterface.
func newFilesWriter() Logger {
	return &filesLogWriter{}
}

func init() {
	Register("files", NewConn)
}
