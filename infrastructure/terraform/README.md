# Terraform lambda deployment

	Prerequisites -  Setup credo-auth AWS profile (see https://credomobile.atlassian.net/wiki/spaces/MEP/blog/2018/02/22/331481104/Terraform+the+DevOps+Way)

	To initialize your project add any additional AWS resources to main.tf. Be sure to use tags and the sumo forwarder for any new resources.

To initialize your terraform directory:

	source ./setEnv.sh dev 
	terraform plan
	
	
To change environments you will need to call setEnv.sh again with the new env:

	source ./setEnv.sh {{env (dev/qa/prod)}}
