A module for extending Terratest options

# Using this Module

At the time of this Readme being updated, that is still a private mednax-it internal repository. This means that `go get` and `go mod tidy` cannot download and update without some help.

In your local, complete the following steps to install the module:
1. This requires you to have ssh set up for your git connection already. If you haven't done this before, see [this link](https://docs.github.com/en/authentication/connecting-to-github-with-ssh)
2. run the following bash commands:
   ```bash
   go env -w GOPRIVATE="github.com/mednax-it/*"
   git config --global url."ssh://git@github.com/".insteadOf "https://github.com/"
   ```

3. Add the following to your import statements: `import terratestPlus "github.com/mednax-it/terratest-plus"`

# Environment Variables
The following environment variables can be used, and are prioritized by this module

* `TF_source` the source directory for terraform. This is relative to wherever the terratest main file is being run from. So if the terratest main file isin the root, and terraform are in the `src` directory off that root, just `src/` is enough.
  * If the terratest files are in a parallel directory, then you can use the `../src` notation to go back up a directory.
* `TF_var_file` the path to the var file to be used, *relative to the terraform directory* - so if the terraform files are in `src` and the var files in `vars` and the file you want to use is `local.tfvars` the env variable should be set to `vars/local.tfvars`
* `TF_backend` the path to the backend file to use. Similar to vars, it is relative to the terraform source directory, so `backend/test.tfbackend` if the `backend` directory is in the `src` directory
* `TF_workspace` is used to set the name of the workspace. This is overwritten by Circles `CIRCLE_SHA1` if being run in a pipeline

* `LOG_TERRAFORM` set `=true` if you want to see verbose terraform logging out of Terratest. any other value will hide most of the terraform logs from of the output.

* `SKIP_terraform_init` will skip the terraform init - really only useful for local testing to speed up testing
* `SKIP_terraform_apply` will skip the terraform apply step - again mostly useful for local testing.


# Service Principle

This is currently set up to use either local if nothing exists or an Azure Service Principle. As such the following env values need to be set:

```
ARM_TENANT_ID
ARM_CLIENT_SECRET
ARM_CLIENT_ID
ARM_SUBSCRIPTION_ID
```

alternatively they can be provided in the var file as: (but be aware these are secrets and should not be committed to a repo. Safer to use the env values)

```
tenant_id
client_secret
client_id
subscription_id
```

# Modifying for your needs

Everyone's terraform is going to be different. This module can only do so much, so its intended to provide a few tools and to be extended to fit your specific terraform needs. Create a struct that is compromised of `terratestPlus.Deployment` and you will have all the basics to extend to your needs.
