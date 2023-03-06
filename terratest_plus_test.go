package terratestPlus

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassNilForOptionsCreatesDefault(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	expectedPath := "vars/local.tfvars"

	testStruct.SetupTerraform(t, nil)

	assert.Equal(expectedPath, testStruct.VarFilePath, "Nill Options passed to SetupTerraform did not result in default Options")
}

func TestSetupTerraformSetsLocalVarsWithNoInformationToStart(t *testing.T) {
	assert := assert.New(t)

	testStruct := new(Deployment)
	expectedPath := "vars/local.tfvars"
	options := new(SetupTerraformOptions)

	testStruct.SetupTerraform(t, options)

	assert.Equalf(expectedPath, testStruct.VarFilePath, "VarFilePath of %s did not default to %s", testStruct.VarFilePath, expectedPath)

}

func TestTFvarfileEnvVariableHasPriorityOverOtherOptions(t *testing.T) {
	assert := assert.New(t)
	expectedPath := "vars/not_local.tfvars"
	os.Setenv("TF_var_file", expectedPath)
	defer os.Unsetenv("TF_var_file")

	testStruct := new(Deployment)
	options := new(SetupTerraformOptions)
	options.VarFileDirectoryPath = "var2/"

	testStruct.SetupTerraform(t, options)

	assert.Equalf(expectedPath, testStruct.VarFilePath, "VarFilePath of %s did not default to %s", testStruct.VarFilePath, expectedPath)

}

func TestStringVarInVarFile(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	variableName := "test_string"
	expectedValue := "a_string_value"

	testStruct.SetupTerraform(t, nil)

	testStringVar := testStruct.VarFileValues[variableName].(string)
	assert.Equalf(expectedValue, testStringVar, "Var File variable %s does not equal %s", variableName, expectedValue)
}

func TestDefaultVarFileOfLocalForcesInitTrue(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	testStruct.SetupTerraform(t, nil)

	assert.True(testStruct.RunInit)
	if testStruct.RunInit {
		fmt.Print("Weeee")
	}
}

func TestDefaultVarFileOfLocalForcesExecuteInLocalFlagTrue(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	testStruct.SetupTerraform(t, nil)

	assert.True(testStruct.ExecutingInLocal)
}
