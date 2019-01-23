package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	mavenMeta, err := DownloadString("https://maven.fabricmc.net/net/fabricmc/yarn/maven-metadata.xml")
	if err != nil {
		panic(err)
	}

	mavenMetadata := MavenMetadata{}
	xml.Unmarshal([]byte(mavenMeta), &mavenMetadata)

	jsonMap := make(map[string][]int)

	for _, version := range mavenMetadata.Versioning.Versions.Version {
		mcVer, build := parseVersion(version)
		jsonMap[mcVer] = append(jsonMap[mcVer], build)
	}

	json, err := json.Marshal(jsonMap)
	if err != nil {
		panic(err)
	}

	jsonStr := string(json)

	//TODO dont hard code this
	WriteStringToFile(jsonStr, "/home/webdata/maven/net/fabricmc/yarn/versions.json")

	fmt.Println(jsonStr)
}

func parseVersion(input string) (string, int) {
	split := strings.Split(input, ".")
	if strings.Contains(input, "-") {
		split = strings.Split(input, "-")
	}
	build, err := strconv.Atoi(split[1])
	if err != nil {
		panic(err)
	}
	return split[0], build
}

func DownloadString(url string) (string, error) {
	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 { // OK
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			return "", err2
		}
		bodyString := string(bodyBytes)
		return bodyString, nil
	}

	return "", errors.New("Failed to download file")
}

func WriteStringToFile(str string, file string) {
	ioutil.WriteFile(file, []byte(str), 0644)
}

type MavenMetadata struct {
	XMLName    xml.Name                `xml:"metadata"`
	GroupID    string                  `xml:"groupId"`
	ArtifactID string                  `xml:"artifactId"`
	Versioning MavenMetadataVersioning `xml:"versioning"`
}

type MavenMetadataVersioning struct {
	Latest      string                `xml:"latest"`
	Release     string                `xml:"release"`
	Versions    MavenMetadataVersions `xml:"versions"`
	LastUpdated string                `xml:"lastUpdated"`
}

type MavenMetadataVersions struct {
	Version []string `xml:"version"`
}
