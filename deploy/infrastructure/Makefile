.PHONY: dev, shutdown-dev
dev:
	cd terraform/env/Dev && \
	terraform init && \
	terraform apply -auto-approve && \
	./configure.bash
	ansible-playbook -i ansible/hosts.yaml ansible/playbook.yaml --limit dev

.PHONY: shutdown-dev
shutdown-dev:
	cd terraform/env/Dev && \
	terraform destroy -auto-approve


.PHONY: prod, shutdown-prod
prod:
	cd terraform/env/Prod && \
	terraform init && \
	terraform apply -auto-approve && \
	./configure.bash

.PHONY: shutdown-prod
shutdown-prod:
	cd terraform/env/Prod && \
	terraform destroy -auto-approve
