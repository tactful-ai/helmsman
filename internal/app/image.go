package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"
)

const (
	imageConfigmapName = "deck-images-config"
)

var (
	images = make(map[string][]ImageVersion)
)

type ImageVersion struct {
	Name      string    `json:"name"`
	CreatedOn time.Time `json:"createdOn"`
	Version   string    `json:"version"`
	Status    string    `json:"status"`
}

// getImageVersions extracts image versions from secret storage
func createImageVersions(namespace string) {
	cmd := kubectl([]string{"create", "configmap", imageConfigmapName, "-n", namespace}, "Creating Image history from secret storage")

	result := cmd.exec()
	if result.code != 0 {
		log.Fatal(result.errors)
	}
}

// getImageVersions extracts image versions from secret storage
func LoadImageVersions(namespace string) string {
	cmd := kubectl([]string{"get", "configmap", imageConfigmapName, "-n", namespace, "-o", "json"}, "Getting Image history from secret storage")

	result := cmd.exec()
	if result.code != 0 {
		log.Fatal(result.errors)
		createImageVersions(namespace)
	}
	fmt.Print(result)
	rctx := strings.Trim(result.output, `"' `)
	return rctx
}

// getImageVersions extracts image versions from secret storage
func addImageVersion(namespace string, imageName string, versionStr string) {

	version := ImageVersion{
		Name:      imageName,
		CreatedOn: time.Now(),
		Status:    "New",
		Version:   versionStr,
	}

	if images[imageName] == nil {
		images[imageName] = []ImageVersion{}
	}
	images[imageName] = append(images[imageName], version)

	saveVersionHistory(namespace, images)
}

// saveVersionHistory saves all image history
func saveVersionHistory(namespace string, images map[string][]ImageVersion) {

	if len(images) == 0 {
		return
	}

	definition := `
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: NAME
data:
  images: |
`
	d, err := json.Marshal(&images)
	fmt.Print("Yaml..\n")
	fmt.Print(string(d))
	if err != nil {
		log.Fatal(err.Error())
	}

	definition = strings.ReplaceAll(definition, "NAME", imageConfigmapName)
	definition = definition + Indent(string(d), strings.Repeat(" ", 8))
	targetFile := path.Join(createTempDir(tempFilesDir, "tmp"), "temp-ImageVersions.yaml")
	fmt.Print("Yaml222..\n")
	fmt.Print(definition)

	if err := ioutil.WriteFile(targetFile, []byte(definition), 0666); err != nil {
		log.Fatal(err.Error())
	}

	cmd := kubectl([]string{"apply", "-f", targetFile, "-n", namespace, flags.getKubeDryRunFlag("apply")}, "Creating Image versions configmap in namespace [ "+namespace+" ]")
	result := cmd.exec()

	if result.code != 0 {
		log.Fatal("Failed to create ImageVersions in namespace [ " + namespace + " ] with error: " + result.errors)
	}

	deleteFile(targetFile)

}
