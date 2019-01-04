#Ansible management of vault secrets

#Operator instructions
Need to install pipenv (available from brew)

	pipenv install
	pipenv shell
	

* make sure dev vault credentials are setup for your shell (VAULT_ADDR and VAULT_TOKEN should be set)
* make sure vault certs are available at ${HOME}/.credo/vault/vault.crt (devenvironment start.sh should put them there)

from pipenv shell

	ansible-playbook vault-secrets.yml --vault-id @prompt --extra-vars "env=dev"

If you encounter TLS issues: Recommend using pyenv for python (system python uses outdated version of openssl)