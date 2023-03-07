package terratestPlus

import (
	"testing"

	"github.com/mednax-it/terratest-plus/deployment"
	"github.com/stretchr/testify/assert"
)

var TestTypeOne string = "resourceType"
var TestModule string = "resourceModule"
var TestNameOne string = "testNameOne"
var TestNameTwo string = "testNameButDifferent"
var TestAttributeName string = "AttributeName"
var TestAttributeNameTwo string = "AttributeNameAsWell"

func SetupMockState() *deployment.TerraformState {
	tags := map[string]interface{}{
		"Tag1": "tag1Value",
	}
	attributes := deployment.StateResourceAttributes{
		Id:       "AttributeID",
		Location: "AttributeLocation",
		Name:     TestAttributeName,
		Tags:     tags,
	}
	attributesTwo := deployment.StateResourceAttributes{
		Id:       "AttributeID",
		Location: "AttributeLocation",
		Name:     TestAttributeNameTwo,
		Tags:     tags,
	}

	resourceInstance := deployment.StateResourceInstance{
		Attributes: &attributes,
	}
	resourceInstanceTwo := deployment.StateResourceInstance{
		Attributes: &attributesTwo,
	}

	resource := deployment.StateResource{
		Module:    TestModule,
		Mode:      "resourceMode",
		Type:      TestTypeOne,
		Name:      TestNameOne,
		Provider:  "resourceProvider",
		Instances: []*deployment.StateResourceInstance{&resourceInstance, &resourceInstanceTwo},
	}

	resourceSameType := deployment.StateResource{
		Module:    TestModule,
		Mode:      "resourceMode",
		Type:      TestTypeOne,
		Name:      TestNameTwo,
		Provider:  "resourceProvider",
		Instances: []*deployment.StateResourceInstance{&resourceInstance},
	}

	state := deployment.TerraformState{
		Resources: []*deployment.StateResource{&resource, &resourceSameType},
	}

	return &state
}

func TestFindAllResourceTypeFindsTheAppropriateResources(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	testStruct.State = SetupMockState()
	expectedLength := 2

	foundResources := testStruct.FindAllResourceType(TestTypeOne)

	assert.Equalf(expectedLength, len(foundResources), "Incorrect number [%s] of resources returned (should be %s)", len(foundResources), expectedLength)
}

func TestFindAllResourceTypeBuildsKeysBasedOnResourceData(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	testStruct.State = SetupMockState()

	foundResources := testStruct.FindAllResourceType(TestTypeOne)

	keyShouldBe := TestModule + "." + TestTypeOne + "." + TestNameOne
	_, found := foundResources[keyShouldBe]
	assert.True(found, "Key [%s] was not found in the output", keyShouldBe)

}

func TestFindByNameReturnsTheAppropriateResources(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	testStruct.State = SetupMockState()
	expectedLength := 1

	foundResources := testStruct.FindByName(TestNameOne)

	assert.Equalf(expectedLength, len(foundResources), "Incorrect number [%s] of resources returned (should be %s)", len(foundResources), expectedLength)

}

func TestFindByNameBuildsKeysBasedOnResourceData(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	testStruct.State = SetupMockState()

	foundResources := testStruct.FindByName(TestNameOne)

	keyShouldBe := TestModule + "." + TestTypeOne + "." + TestNameOne
	_, found := foundResources[keyShouldBe]
	assert.True(found, "Key [%s] was not found in the output", keyShouldBe)

}

func TestGetInstanceNamesReturnsTheCorrectNumberOfNames(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	testStruct.State = SetupMockState()
	resources := testStruct.FindAllResourceType(TestTypeOne)
	expectedLength := 3

	foundNames := testStruct.GetInstanceNames(resources)

	assert.Equalf(expectedLength, len(foundNames), "Incorrect number [%s] of names returned (should be %s)", len(foundNames), expectedLength)

}

func TestGetInstanceNamesReturnsCorrectName(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	testStruct.State = SetupMockState()
	resources := testStruct.FindByName(TestNameTwo)

	foundNames := testStruct.GetInstanceNames(resources)

	assert.Containsf(foundNames, TestAttributeName, "Expected Name of %s was not found in %s", TestAttributeName, foundNames[0])

}
