#!/usr/bin/env bash
if [ $# = 0 ] ; then
    echo "!! must specify the environment (dev/qa/prod)"
    echo "   ./setEnv.sh dev"
    return
fi

ENV_FILE="./backendConfigs/$1"
if [ ! -f ${ENV_FILE} ] ; then
	echo "File $ENV_FILE not found"
	return
fi

if [ $1 == "dev" ] ; then
    env=674346455231
elif [ $1 == "qa" ] ; then
    env=772404289823
elif [ $1 == "prod" ] ; then
    env=465292320167
else
    echo "unknown env: $1"
    return
fi

role="terraform"

rm -rf ./.terraform/

terraform init -backend-config=${ENV_FILE} -backend-config=role_arn="arn:aws:iam::${env}:role/${role}" -backend-config=profile=credo-auth

# attempt to set TF_VAR_environment to the appropriate env.   Will only work if script is sourced
export TF_VAR_environment=$1
export TF_VAR_role="${role}"
