package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"gitlab.com/c0b/go-ordered-json"
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

	om := ordered.NewOrderedMap()

	for _, version := range mavenMetadata.Versioning.Versions.Version {
		mcVer, build := parseVersion(version)
		var values []int
		if om.Has(mcVer) {
			values = om.Get(mcVer).([]int)
		}
		values = append(values, build)
		om.Set(mcVer, values)
	}

	json, err := json.Marshal(om)
	if err != nil {
		panic(err)
	}

	jsonStr := string(json)

	//TODO dont hard code this
	WriteStringToFile(jsonStr, "/home/webdata/maven/net/fabricmc/yarn/versions.json")

	fmt.Println(jsonStr)

}

type Member struct {
	Key   string
	Value interface{}
}

type OrderedObject []Member

func parseVersion(input string) (string, int) {
	if strings.Contains(input, "+build.") {
		return parseVersionNew(input)
	}
	splitpos := strings.LastIndex(input, ".")
	if strings.Contains(input, "-") {
		splitpos = strings.LastIndex(input, "-")
	}
	mcver := input[:splitpos]
	build, err := strconv.Atoi(input[splitpos+1:])
	if err != nil {
		panic(err)
	}
	return mcver, build
}

func parseVersionNew(input string) (string, int) {
	build, err := strconv.Atoi(input[strings.LastIndex(input, ".")+1:])
	if err != nil {
		panic(err)
	}
	return input[:strings.LastIndex(input, "+")], build
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
