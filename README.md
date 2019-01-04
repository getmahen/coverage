
#coverage lambda
Enter a brief description of the lambda here.

##Environmental Variables 
List environmental variables here

# consul variables
Put consul variables used here

# vault secrets
Put names of vault secrets used here

# aws resources
Put what additional AWS resources are used ( S3, SSM params, etc)

# aws trigger
Document what trigger is used for this lambda

# Deployment 
To deploy this lambda to dev:

	make test
	make build
	make upload #(need DEV aws creds)

cd to `./infrastructure/terraform` and follow README.md in that directory