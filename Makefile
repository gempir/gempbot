provision: 
	ansible-playbook -i ansible/hosts ansible/playbook.yml --ask-vault-pass ${ARGS}

bot:
	cd services/bot && go build

ingester:
	cd services/ingester && go build