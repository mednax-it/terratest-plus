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

func SetupMockState() *deployment.TerraformState {
	tags := map[string]interface{}{
		"Tag1": "tag1Value",
	}
	attributes := deployment.StateResourceAttributes{
		Id:       "AttributeID",
		Location: "AttributeLocation",
		Name:     "AttributeName",
		Tags:     tags,
	}

	resourceInstance := deployment.StateResourceInstance{
		Attributes: &attributes,
	}

	resource := deployment.StateResource{
		Module:    TestModule,
		Mode:      "resourceMode",
		Type:      TestTypeOne,
		Name:      TestNameOne,
		Provider:  "resourceProvider",
		Instances: []*deployment.StateResourceInstance{&resourceInstance},
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
