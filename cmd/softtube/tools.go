package main

//
//func getResourcePath(fileName string) (string, error) {
//	exePath, err := os.Executable()
//	if err != nil {
//		return "", err
//	}
//	exeDir := path.Dir(exePath)
//
//	fw := framework.NewFramework()
//
//	gladePath := path.Join(exeDir, fileName)
//	if fw.IO.FileExists(gladePath) {
//		return gladePath, nil
//	}
//	gladePath = path.Join(exeDir, "assets", fileName)
//	if fw.IO.FileExists(gladePath) {
//		return gladePath, nil
//	}
//	gladePath = path.Join(exeDir, "../assets", fileName)
//	if fw.IO.FileExists(gladePath) {
//		return gladePath, nil
//	}
//	return gladePath, nil
//}
