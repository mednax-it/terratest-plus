A module for extending Terratest options


# Environment Variables
The following environment variables can be used, and are prioritized by this module

* `TF_source` the source directory for terraform. This is relative to wherever the terratest main file is being run from. So if the terratest main file isin the root, and terraform are in the `src` directory off that root, just `src/` is enough.
  * If the terratest files are in a parallel directory, then you can use the `../src` notation to go back up a directory.
* `TF_var_file` the path to the var file to be used, *relative to the terraform directory* - so if the terraform files are in `src` and the var files in `vars` and the file you want to use is `local.tfvars` the env variable should be set to `vars/local.tfvars`
* `TF_backend` the path to the backend file to use. Similar to vars, it is relative to the terraform source directory, so `backend/config.test_backend.tfbackend` if the `backend` directory is in the `src` directory
* `TF_workspace` is used to set the name of the workspace. This is overwritten by Circles `CIRCLE_SHA1` if being run in a pipeline

* `SKIP_terraform_init` will skip the terraform init - really only useful for local testing to speed up testing
* `SKIP_terraform_apply` will skip the terraform apply step - again mostly useful for local testing.