package core


import (
	"errors"
	"os"
	"path"
	"path/filepath"
)

// Resources : Handles SoftTeam resources
type Resources struct {
}

func NewResources() *Resources {
	return new(Resources)
}

// GetExecutablePath : Returns the path of the executable
func (r *Resources) GetExecutablePath() string {
	ex, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(ex)
}

// GetResourcesPath : Returns the resources path
func (r *Resources) GetResourcesPath() string {
	executablePath:=r.GetExecutablePath()

	var pathsToCheck []string
	pathsToCheck = append(pathsToCheck,path.Join(executablePath, "assets"))
	pathsToCheck = append(pathsToCheck,path.Join(executablePath, "../assets"))

	dir, err := r.checkPathsExists(pathsToCheck)
	if err!=nil {
		return executablePath
	}
	return dir
}

func (r *Resources) checkPathsExists(pathsToCheck []string) (string, error) {
	for _,path := range pathsToCheck {
		if _, err := os.Stat(path); os.IsNotExist(err) == false {
			return path, nil
		}
	}
	return "", errors.New("paths do not exist")
}

// GetResourcePath : Gets the path for a single resource file
func (r *Resources) GetResourcePath(fileName string) string {
	resourcesPath:=r.GetResourcesPath()
	resourcePath:=path.Join(resourcesPath, fileName)
	return resourcePath
}
