/*******************************************************************************
 * Copyright (c) 2019 IBM Corporation and others.
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v2.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v20.html
 *
 * Contributors:
 *     IBM Corporation - initial API and implementation
 *******************************************************************************/

package project

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

var goTestFolder, goPackagePath, goShortPackagePath, nodeTestFolder string
var nodeOutput = &map[string]interface{}{}
var goOutput = &map[string]interface{}{}

func before() {
	goTestFolder = "insights_test_folder_delete_me"
	goVendorPath := path.Join(goTestFolder, "vendor")
	goSourcePath := path.Join(goVendorPath, "source")
	goAuthorPath := path.Join(goSourcePath, "author")
	goPackagePath = path.Join(goAuthorPath, "package")

	goShortPackagePath = path.Join(goSourcePath, "package")

	os.Mkdir(goTestFolder, 0777)
	os.Mkdir(goVendorPath, 0777)
	os.Mkdir(goSourcePath, 0777)
	os.Mkdir(goAuthorPath, 0777)
	os.Mkdir(goPackagePath, 0777)

	os.Mkdir(goShortPackagePath, 0777)
	os.Create(goShortPackagePath + "/test.txt")

	goPackageFullPath := map[string]interface{}{
		"license-analysis": "no license file was found",
		"path":             goPackagePath,
	}
	goPackageShortPath := map[string]interface{}{
		"license-analysis": "no license file was found",
		"path":             goShortPackagePath,
	}
	(*goOutput)["source/author/package"] = goPackageFullPath
	(*goOutput)["source/package"] = goPackageShortPath

	fullNodeModule := map[string]interface{}{
		"declaredLicenses": "MIT",
		"description":      "a test module",
		"license-analysis": "no license file was found",
		"name":             "full-data",
		"path":             "../../resources/insights-testing/node-project/node_modules/full-data",
		"version":          "1.0.0",
	}
	missingDataNodeModule := map[string]interface{}{
		"declaredLicenses": "No Declared License",
		"license-analysis": "no license file was found",
		"name":             "missing-data",
		"path":             "../../resources/insights-testing/node-project/node_modules/missing-data",
		"version":          "2.0.0",
	}
	nestedNodeModule := map[string]interface{}{
		"declaredLicenses": "(ISC OR GPL-3.0)",
		"license-analysis": "no license file was found",
		"name":             "nested-module",
		"path":             "../../resources/insights-testing/node-project/node_modules/full-data/node_modules/nested-module",
		"version":          "1.0.0",
	}
	(*nodeOutput)["full-data@1.0.0"] = fullNodeModule
	(*nodeOutput)["missing-data@2.0.0"] = missingDataNodeModule
	(*nodeOutput)["nested-module@1.0.0"] = nestedNodeModule

	nodeTestFolder = "../../resources/insights-testing/node-project"
}

func after() {
	os.RemoveAll(goTestFolder)
}

func TestProjectInsights(t *testing.T) {
	before()

	tests := map[string]struct {
		path           string
		language       string
		expectedOutput *map[string]interface{}
		expectsError   bool
	}{
		"success case: invalid language error returned for an unknown language": {
			path:           nodeTestFolder,
			language:       "unknown",
			expectedOutput: nil,
			expectsError:   true,
		},
		"success case: error returned for an invalid project path": {
			path:           "this-path-does-not-exist",
			language:       "unknown",
			expectedOutput: nil,
			expectsError:   true,
		},
		"success case: empty response returned when no node dependencies are installed": {
			path:           "../../resources/insights-testing",
			language:       "nodejs",
			expectedOutput: &map[string]interface{}{},
			expectsError:   false,
		},
		"success case: empty response returned when no go dependencies are installed": {
			path:           "../../resources/insights-testing",
			language:       "go",
			expectedOutput: &map[string]interface{}{},
			expectsError:   false,
		},
		"success case: returns correct response for a node project": {
			path:           nodeTestFolder,
			language:       "nodejs",
			expectedOutput: nodeOutput,
			expectsError:   false,
		},
		"success case: returns correct response for a go project": {
			path:           goTestFolder,
			language:       "go",
			expectedOutput: goOutput,
			expectsError:   false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			response, err := ProduceInsights(test.language, test.path)

			assert.Exactly(t, test.expectedOutput, response, "insights gave incorrect response")

			if test.expectsError {
				assert.Error(t, err, "insights did not return an error when one was expected")
			}
		})
	}
	after()
}
