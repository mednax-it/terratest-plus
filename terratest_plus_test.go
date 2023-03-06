package terratestPlus

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassNilForOptionsCreatesDefault(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	expectedPath := "vars/local.tfvars"

	testStruct.SetupTerraform(t, nil)

	assert.Equal(expectedPath, testStruct.VarFilePath, "Nill Options passed to SetupTerraform did not result in default VarFilePath")
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

func TestMapVarInVarFile(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	variableName := "test_map"
	expectedValue := "test_map_value"

	testStruct.SetupTerraform(t, nil)

	testMapVar := testStruct.VarFileValues[variableName].(map[string]interface{})
	testValue := testMapVar["test_map_key"].(string)
	assert.Equalf(expectedValue, testValue, "Var File variable %s does not equal %s", variableName, expectedValue)
}

func TestArrayVarInVarFile(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	variableName := "test_array"
	expectedValue := "var1"

	testStruct.SetupTerraform(t, nil)

	testArrayVar := testStruct.VarFileValues[variableName].([]interface{})
	testValue := testArrayVar[0].(string)
	assert.Equalf(expectedValue, testValue, "Var File variable %s does not equal %s", variableName, expectedValue)
}

func TestDefaultVarFileOfLocalForcesInitTrue(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	testStruct.SetupTerraform(t, nil)

	assert.True(testStruct.RunInit)
}

func TestDefaultVarFileOfLocalForcesExecuteInLocalFlagTrue(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	testStruct.SetupTerraform(t, nil)

	assert.True(testStruct.ExecutingInLocal)
}

func TestBackendPathWhenNilOptionsDefaults(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	expectedPath := "backends/config.test_backend.tfbackend"

	testStruct.SetupTerraform(t, nil)

	assert.Equal(expectedPath, testStruct.BackendFilePath, "Nill Options passed to SetupTerraform did not result in default BackendPath")

}

func TestBackendPathWhenProvidedBlankOptionsDefaults(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	expectedPath := "backends/config.test_backend.tfbackend"
	options := new(SetupTerraformOptions)

	testStruct.SetupTerraform(t, options)

	assert.Equal(expectedPath, testStruct.BackendFilePath, "Blank Options passed to SetupTerraform did not result in default BackendPath")

}

func TestBackendWhenSetByEnvVariable(t *testing.T) {
	assert := assert.New(t)
	expectedPath := "backends/config.test_backend.tfbackend"
	os.Setenv("TF_backend", expectedPath)
	defer os.Unsetenv("TF_backend")

	testStruct := new(Deployment)
	options := new(SetupTerraformOptions)
	options.BackendDirectoryPath = "backend2/"

	testStruct.SetupTerraform(t, options)

	assert.Equalf(expectedPath, testStruct.BackendFilePath, "BackendPath of %s did not follow env value of  %s ", testStruct.VarFilePath, expectedPath)

}
