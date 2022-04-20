package version

var Version *version = new(version)

type version struct {
	LokiVersion string
	LokiType    string
	GoVersion   string
}
