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

var TestNameThree string = "testNameWithNumberIndex"
var TestAttributeName string = "AttributeName"
var TestAttributeNameTwo string = "AttributeNameAsWell"
var TestAttributeValue string = "TestAttributeValue"
var TestAttributeValueTwo string = "TestAttributeValueTwo"
var index1 interface{} = "Index1"
var indexNumber interface{} = 0

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
		IndexKey:   &index1,
	}
	resourceInstanceTwo := deployment.StateResourceInstance{
		Attributes: &attributesTwo,
		IndexKey:   nil,
	}

	resourceInstanceThree := deployment.StateResourceInstance{
		Attributes: &attributes,
		IndexKey:   &indexNumber,
	}

	resource := deployment.StateResource{
		Module:    TestModule,
		Mode:      "managed",
		Type:      TestTypeOne,
		Name:      TestNameOne,
		Provider:  "resourceProvider",
		Instances: []*deployment.StateResourceInstance{&resourceInstance, &resourceInstanceTwo},
	}

	resourceSameType := deployment.StateResource{
		Module:    TestModule,
		Mode:      "managed",
		Type:      TestTypeOne,
		Name:      TestNameTwo,
		Provider:  "resourceProvider",
		Instances: []*deployment.StateResourceInstance{&resourceInstance},
	}

	resourceWithNumberIndex := deployment.StateResource{
		Module:    TestModule,
		Mode:      "managed",
		Type:      TestTypeOne,
		Name:      TestNameThree,
		Provider:  "resourceProvider",
		Instances: []*deployment.StateResourceInstance{&resourceInstanceThree},
	}

	state := deployment.TerraformState{
		Resources: []*deployment.StateResource{&resource, &resourceSameType, &resourceWithNumberIndex},
	}

	return &state
}

func SetupMockRawState() map[string]interface{} {
	rawMockAttributes := map[string]interface{}{
		TestAttributeName:    TestAttributeValue,
		TestAttributeNameTwo: TestAttributeValueTwo,
	}
	rawMockInstance := map[string]interface{}{
		"attributes": rawMockAttributes,
	}
	rawMockResource := map[string]interface{}{
		"module":    TestModule,
		"mode":      "managed",
		"type":      TestTypeOne,
		"name":      TestNameOne,
		"instances": []interface{}{rawMockInstance, rawMockInstance},
	}
	rawMockState := map[string]interface{}{
		"resources": []interface{}{rawMockResource},
	}
	return rawMockState
}

func TestFindAllResourceTypeFindsTheAppropriateResources(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	testStruct.State = SetupMockState()
	expectedLength := 3

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
	keyShouldBe := TestModule + "." + TestTypeOne + "." + TestNameOne

	foundResources := testStruct.FindByName(TestNameOne)

	_, found := foundResources[keyShouldBe]
	assert.True(found, "Key [%s] was not found in the output", keyShouldBe)

}

func TestGetInstanceNamesReturnsTheCorrectNumberOfNames(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	testStruct.State = SetupMockState()
	resources := testStruct.FindAllResourceType(TestTypeOne)
	expectedLength := 4

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

func TestCompileAllRegexGroups(t *testing.T) {
	assert := assert.New(t)
	regexString := `module.platforms\[(?P<platform_id>\w*)\]`
	expectedPlatformID := "TestPlatformID"
	stringToMatch := "module.platforms[" + expectedPlatformID + "]"

	parameterMap := CompileAllRegexGroups(regexString, stringToMatch)

	assert.Equalf(expectedPlatformID, parameterMap["platform_id"], "The value for [platform_id] was %s and did not match expected %s", parameterMap["platform_id"], expectedPlatformID)
}

func TestGetInstanceAttributeValuesOfResourceReturnsAsExpectedAmmountWithResourceType(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	testStruct.RawState = SetupMockRawState()
	expectedLength := 2

	foundAttributes := testStruct.GetInstanceAttributeValuesOfResource(TestTypeOne, TestAttributeNameTwo)

	assert.Equalf(expectedLength, len(foundAttributes), "The length of returned value (%s) was longer than expected %s", len(foundAttributes), expectedLength)

}

func TestGetInstanceAttributeValuesOfResourceReturnsAsExpectedAmmountWithResourceName(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	testStruct.RawState = SetupMockRawState()
	expectedLength := 2

	foundAttributes := testStruct.GetInstanceAttributeValuesOfResource(TestNameOne, TestAttributeNameTwo)

	assert.Equalf(expectedLength, len(foundAttributes), "The length of returned value (%s) was longer than expected %s", len(foundAttributes), expectedLength)

}

func TestIndexKeysForStringsAsIndexFromMaps(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	testStruct.State = SetupMockState()
	resources := testStruct.FindByName(TestNameTwo)

	for key, resource := range resources {
		for _, instance := range resource.Instances {
			assert.Equalf(index1, *instance.IndexKey, "Index key (%s) in %s resource does not equal expected %s", *instance.IndexKey, key, index1)
		}
	}

}

func TestIndexKeysForNumbersAsIndexFromList(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	testStruct.State = SetupMockState()
	resources := testStruct.FindByName(TestNameThree)

	for key, resource := range resources {
		for _, instance := range resource.Instances {
			assert.Equalf(indexNumber, *instance.IndexKey, "Index key (%s) in %s resource does not equal expected %s", *instance.IndexKey, key, indexNumber)
		}
	}

}

func TestGetInstanceAttributeValuesReturnsAttributes(t *testing.T) {
	assert := assert.New(t)
	testStruct := new(Deployment)
	testStruct.RawState = SetupMockRawState()

	foundAttributes := testStruct.GetInstanceAttributeValuesOfResource(TestNameOne, TestAttributeNameTwo)

	assert.Containsf(foundAttributes, TestAttributeValueTwo, "The found attributes do not contain the expected one")

}

// Go does not have a default "Expect Fail" option, so I would love a test right here that Tests the FailNow of the function above but alas, it isnt possible with default
