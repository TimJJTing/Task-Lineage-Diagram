package schema

type Task struct {
	Task       string       `yaml:"task"`
	TaskID     string       `yaml:"task_id"`
	StartDate  string       `yaml:"start_date"`
	EndDate    string       `yaml:"end_date"`
	Frequency  int          `yaml:"frequency"`
	Unit       string       `yaml:"unit"`
	Queue      string       `yaml:"queue"`
	Level      string       `yaml:"level"`
	Runtime    Runtime      `yaml:"runtime"`
	Dependency []Dependency `yaml:"dependency"`
}

type Runtime struct {
	Directory  string `yaml:"directory"`
	Executable string `yaml:"executable"`
	File       string `yaml:"file"`
	Extension  string `yaml:"extension"`
}

type Dependency struct {
	TaskID    string `yaml:"task_id"`
	Storage   string `yaml:"storage"`
	Unit      string `yaml:"unit"`
	Frequency int    `yaml:"frequency"`
	StartDate string `yaml:"start_date"`
	EndDate   string `yaml:"end_date"`
}

type Config struct {
	Colors map[string]string `yaml:"colors"`
}
