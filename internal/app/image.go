package app

import (
	"encoding/base64"
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
	images = make(map[string][]imageVersion)
)

// Defines an image name and its json-path in a helm values file
type imageLookup struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type imageVersion struct {
	Name      string    `json:"name"`
	CreatedOn time.Time `json:"createdOn"`
	Version   string    `json:"version"`
	Status    string    `json:"status"`
}

type configMap struct {
	Data map[string]string `json:"data"`
}

func readConfigMap(name string, namespace string) (configMap, error) {
	var configmap configMap
	cmd := kubectl([]string{"get", "configmap", name, "-n", namespace, "-o", "json"}, "Getting Image history from secret storage")

	result := cmd.exec()
	if result.code != 0 {
		log.Fatal(result.errors)
		return configmap, fmt.Errorf("cannot read configmap")
	}

	error := json.Unmarshal([]byte(result.output), &configmap)
	if error != nil {
		log.Fatal(error.Error())
		return configmap, fmt.Errorf("cannot parse json of configmap")
	}

	return configmap, nil
}

// LoadImageVersions extracts image versions from secret storage
func LoadImageVersions(namespace string) {
	configmap, err := readConfigMap(imageConfigmapName, namespace)
	if err != nil {
		log.Fatal(err.Error())
	}
	decoded, err := base64.StdEncoding.DecodeString(configmap.Data["images"])
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := json.Unmarshal(decoded, &images); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("images = %+v\n", images)
}

// getImageVersions extracts image versions from secret storage
func addImageVersion(namespace string, imageName string, versionStr string) {
	found := false
	version := imageVersion{
		Name:      imageName,
		CreatedOn: time.Now(),
		Status:    "New",
		Version:   versionStr,
	}

	if images[imageName] == nil {
		images[imageName] = []imageVersion{}
	}

	// append version if not found already
	// todo: should remove the old and add the new date instead in case ppl reuse their tags for some stupid reason!
	// or maybe we can even use image hash for uniquness
	for _, image := range images[imageName] {
		if image.Name == version.Name {
			found = true
		}

	}
	if !found {
		images[imageName] = append(images[imageName], version)
	}

	saveVersionHistory(namespace, images)
}

// saveVersionHistory saves all image history
func saveVersionHistory(namespace string, images map[string][]imageVersion) {

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
	fmt.Print("preparing Version history string json..\n")
	fmt.Print(string(d))
	encoded := base64.StdEncoding.EncodeToString([]byte(d))

	if err != nil {
		log.Fatal(err.Error())
	}

	definition = strings.ReplaceAll(definition, "NAME", imageConfigmapName)
	definition = definition + Indent(encoded, strings.Repeat(" ", 8))
	targetFile := path.Join(createTempDir(tempFilesDir, "tmp"), "temp-imageVersions.yaml")
	fmt.Print("saving file")
	fmt.Print(definition)

	if err := ioutil.WriteFile(targetFile, []byte(definition), 0666); err != nil {
		log.Fatal(err.Error())
	}

	cmd := kubectl([]string{"apply", "-f", targetFile, "-n", namespace, flags.getKubeDryRunFlag("apply")}, "Creating Image versions configmap in namespace [ "+namespace+" ]")
	result := cmd.exec()

	if result.code != 0 {
		log.Fatal("Failed to create imageVersions in namespace [ " + namespace + " ] with error: " + result.errors)
	}

	deleteFile(targetFile)

}
