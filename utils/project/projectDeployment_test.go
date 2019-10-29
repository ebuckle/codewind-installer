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
	"testing"

	"github.com/stretchr/testify/assert"
)

const testProjectID = "a9384430-f177-11e9-b862-edc28aca827a"
const testConnectionID = "local"

// Test_ProjectConnection :  Tests
func Test_ProjectConnection(t *testing.T) {
	ResetTargetFile(testProjectID)

	t.Run("Asserts there are no target connections", func(t *testing.T) {
		connectionTargets, projError := ListTargetConnections(testProjectID)
		if projError != nil {
			t.Fail()
		}
		assert.Len(t, connectionTargets.ConnectionTargets, 0)
	})

	t.Run("Asserts getting connection URL fails", func(t *testing.T) {
		_, projError := GetConnectionURL(testProjectID)
		if projError == nil {
			t.Fail()
		}
		assert.Equal(t, errOpNotFound, projError.Op)
	})

	t.Run("Add project to local connection", func(t *testing.T) {
		projError := AddConnectionTarget(testProjectID, testConnectionID)
		if projError != nil {
			t.Fail()
		}
	})

	t.Run("Asserts re-adding the same connection fails", func(t *testing.T) {
		projError := AddConnectionTarget(testProjectID, testConnectionID)
		if projError == nil {
			t.Fail()
		}
		assert.Equal(t, errOpConflict, projError.Op)
	})

	t.Run("Asserts there is just 1 target connection added", func(t *testing.T) {
		connectionTargets, projError := ListTargetConnections(testProjectID)
		if projError != nil {
			t.Fail()
		}
		assert.Len(t, connectionTargets.ConnectionTargets, 1)
	})

	t.Run("Asserts an unknown connection can not be removed", func(t *testing.T) {
		projError := RemoveConnectionTarget(testProjectID, "test-AnUnknownConnectionID")
		if projError == nil {
			t.Fail()
		}
		assert.Equal(t, "con_not_found", projError.Op)
	})

	t.Run("Asserts removing a known connection is successful", func(t *testing.T) {
		projError := RemoveConnectionTarget(testProjectID, "local")
		if projError != nil {
			t.Fail()
		}
	})

	t.Run("Asserts there are no targets left for this project", func(t *testing.T) {
		connectionTargets, projError := ListTargetConnections(testProjectID)
		if projError != nil {
			t.Fail()
		}
		assert.Len(t, connectionTargets.ConnectionTargets, 0)
	})

	t.Run("Asserts attempting to manage an invalid project ID fails", func(t *testing.T) {
		projError := AddConnectionTarget("bad-project-ID", testConnectionID)
		if projError == nil {
			t.Fail()
		}
		assert.Equal(t, errOpInvalidID, projError.Op)
	})

}