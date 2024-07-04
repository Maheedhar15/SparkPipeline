spark-create-change:
	@echo "==============Submitting change request"
	touch .env
	TEST_SLACK="https://hooks.slack.com/services/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"
	docker-compose -f 'docker-compose.wrapper.yml' run --rm go-spark-change
	