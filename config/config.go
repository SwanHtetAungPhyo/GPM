package config

type Config struct {
	ModuleName           string
	AppType              string
	Framework            string
	UseDocker            bool
	UseAir               bool
	UseMakefile          bool
	ProjectDir           string
	SelectedDependencies []string
}
