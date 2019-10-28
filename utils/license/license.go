package license

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/src-d/go-license-detector.v3/licensedb"

	"github.com/eclipse/codewind-installer/utils"
)

// ProduceInsights calls the appropriate crawling function for the provided language and then reports on licensing
func ProduceInsights(language string, projectDir string) {
	insightData := make(map[string]interface{})
	if language == "nodejs" {
		println("Node project")
		NodeCrawling(projectDir, insightData)
	} else {
		println("Not a node project")
	}
	PerformLicenseCheck(insightData)
	utils.PrettyPrintJSON(insightData)
}

// NodeCrawling recursively crawls through installed node packages to map dependencies
func NodeCrawling(projectDir string, insightData map[string]interface{}) {
	if utils.PathExists(projectDir + "/node_modules") {
		files, err := ioutil.ReadDir(projectDir + "/node_modules")
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			if file.IsDir() {
				path := projectDir + "node_modules/" + file.Name()
				if utils.PathExists(path + "/package.json") {
					jsonFile, err := os.Open(path + "/package.json")

					if err != nil {
						log.Fatal(err)
					}

					defer jsonFile.Close()

					byteValue, _ := ioutil.ReadAll(jsonFile)

					var result map[string]interface{}
					json.Unmarshal([]byte(byteValue), &result)

					newPackageData := make(map[string]string)
					TransferNodeData(result, newPackageData, path)
					packageID := newPackageData["name"] + "@" + newPackageData["version"]
					if _, ok := insightData[packageID]; !ok {
						insightData[packageID] = newPackageData
					}

					if utils.PathExists(path + "/node_modules") {
						NodeCrawling(path, insightData)
					}
				}
			}
		}
	} else {
		println("No node modules installed")
	}
}

// TransferNodeData takes existing package data from a package.json and loads it into a packageData struct
func TransferNodeData(packageJSON map[string]interface{}, packageData map[string]string, path string) {
	if str, ok := packageJSON["name"].(string); ok {
		packageData["name"] = str
	}
	if str, ok := packageJSON["version"].(string); ok {
		packageData["version"] = str
	}
	if str, ok := packageJSON["description"].(string); ok {
		packageData["description"] = str
	}
	if str, ok := packageJSON["license"].(string); ok {
		packageData["declaredLicenses"] = str
	}
	packageData["path"] = path
}

// PerformLicenseCheck takes an existing map of package data and performs a license check on each package
func PerformLicenseCheck(insightData map[string]interface{}) {
	for _, depI := range insightData {
		dep := depI.(map[string]interface{})
		results := licensedb.Analyse(dep["path"].(string))
		dep["license-analysis"] = results
	}
}
