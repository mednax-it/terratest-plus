# These are needed to use the Service Provider with Terraform backends
export ARM_TENANT_ID=$AZURE_SP_TENANT
export ARM_CLIENT_SECRET=$AZURE_SP_PASSWORD
export ARM_CLIENT_ID=$AZURE_SP
export ARM_SUBSCRIPTION_ID=""24df1984-169b-47a7-95bf-08a1d9434cb2""


# the AZURE_ prefixed env variables are set in the CircleCI mednax-global context and are Secret.
export TF_VAR_tenant_id=$AZURE_SP_TENANT
export TF_VAR_client_secret=$AZURE_SP_PASSWORD
export TF_VAR_client_id=$AZURE_SP
export LOG_TERRAFORM=true
