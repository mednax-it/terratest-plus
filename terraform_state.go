package terratestPlus

import (
	"regexp"

	"github.com/mednax-it/terratest-plus/deployment"
	"github.com/stretchr/testify/assert"
)

/*
FindAllResourceType finds based on a type like `azurerm_kubernetes_cluster`.

The resulting return will be a map of StateResource structs, allowing use of direct attribute calls, such as StateResource.
It will only returned managed resource types, as others are not necessary for testing.
Use FindALlResourceTypeByMode if data or other is needed.
*/
func (d *Deployment) FindAllResourceType(resourceType string) map[string]*deployment.StateResource {

	return d.FindAllResourceTypeByMode(resourceType, "managed")
}

/*
FindAllResourceTypeByMode finds based on a type like `asurerm_kubernetes_cluster` and a mode like `data`.
The resulting return will be a map of StateResource structs, allowing use of direct attribute calls, such as StateResource.
*/
func (d *Deployment) FindAllResourceTypeByMode(resourceType string, mode string) map[string]*deployment.StateResource {

	output := make(map[string]*deployment.StateResource)
	for v := range d.State.Resources {
		resource := d.State.Resources[v]
		if resource.Type == resourceType && resource.Mode == mode {
			full_address := d.cleanTerraformAddress(resource)
			output[full_address] = resource
		}
	}
	return output
}

/*
FindByName finds based on a name returning a map.

Note: if you use `this` as a terraform convention for several different types of resources, then this will return things it shouldn't.
If so use FindAllResourceType instead.
*/
func (d *Deployment) FindByName(resourceName string) map[string]*deployment.StateResource {
	resources := d.State.Resources
	output := make(map[string]*deployment.StateResource)
	for v := range resources {
		resource := resources[v]
		if resource.Name == resourceName {

			full_address := d.cleanTerraformAddress(resource)
			output[full_address] = resource
		}

	}
	return output
}

/*
GetInstanceNames loops through a map of StateResources and gets the name of Each INSTANCE in that map (not the name of the Terraform Resource!)

This is used for Name Convention checking - such as verifying that the resource name in Azure is going to be of a specific pattern.
*/
func (d *Deployment) GetInstanceNames(resources map[string]*deployment.StateResource) []string {

	all_names := make([]string, 0)
	for _, resource := range resources {
		for _, instance := range resource.Instances {
			all_names = append(all_names, instance.Attributes.Name.(string))
		}
	}

	return all_names
}

/*
GetInstanceAttributeValuesOfResource returns a slice of all the values of a specific Attribute for every instance in the state file.

It takes an Identifier, which can be either a terraform resource type (e.g. `azurerm_resource_group`) or a resource name.

It ignores any resource that is not Managed - so Data and Local resources are ignored and will not be found.

The return value is a slice of interfaces, so it will need to be defined depending on what type the attribute is (i.e `value.(string)` )

Will Cause the test to fail if it cannot find the attribute on the provided list of resources.

This is most useful for finding specific attributes not covered in StateResource and subsequent structs.

Note if you provide a name, if you have different resource types with the same name (such as the terraform convention `this`) this function will produce incorrect results.
*/
func (d *Deployment) GetInstanceAttributeValuesOfResource(identifier string, attribute_name string) []interface{} {
	all_values := make([]interface{}, 0)

	resources := d.RawState["resources"].([]interface{})
	for _, v := range resources {
		resource := v.(map[string]interface{})
		if resource["mode"] == "managed" && (resource["type"] == identifier || resource["name"] == identifier) {
			terraform_address := d.cleanTerraformAddress(resource)

			instances := resource["instances"].([]interface{})
			for _, inst_val := range instances {
				instance := inst_val.(map[string]interface{})
				attributes := instance["attributes"].(map[string]interface{})

				// verify the value exists before trying to append it
				if val, ok := attributes[attribute_name]; ok {
					all_values = append(all_values, val)
				} else {
					instance_name := terraform_address
					if name, ok := attributes["name"]; ok {
						instance_name = "." + name.(string)
					}

					assert.Failf(d.T, "Unable to find %s on %s", attribute_name, instance_name)
				}

			}
		}

	}

	return all_values
}

/*
CompileAllRegexGroups takes a regex string and returns a map of the parameters.
This Function requires the use of the regex functionally to name sub groups.

e.g.
`(?P<Year>\d{4})-(?P<Month>\d{2})-(?P<Day>\d{2})`
the `?P<name>` names the subgroup within the () and allows it to be pulled for this function.

NOTE: Does not return an error if the compile function fails! Verify your regex patterns.

(We do not use regexp.MustCompile because otherwise its a panic that could cause the testing framework to not clean up)
*/
func CompileAllRegexGroups(regEx, stringToMatch string) map[string]string {

	compRegEx, _ := regexp.Compile(regEx)

	match := compRegEx.FindStringSubmatch(stringToMatch)

	paramsMap := make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}

/*
CleanTerraformAddress takes a resource and returns a terraform full address of module.type.name format.

It will cause a FailNow if it cannot create the address (in order for cleanup on tests to continue!)
*/
func (d *Deployment) cleanTerraformAddress(resource interface{}) string {
	if _, ok := resource.(*deployment.StateResource); ok {
		value := resource.(*deployment.StateResource)
		return value.Module.(string) + "." + value.Type.(string) + "." + value.Name.(string)
	}

	if _, ok := resource.(map[string]interface{}); ok {
		value := resource.(map[string]interface{})
		return value["module"].(string) + "." + value["type"].(string) + "." + value["name"].(string)
	}

	assert.FailNow(d.T, "Could not create a clean Terraform Address! Something is wrong in the TerratestPlus functions")

	return ""
}
