spark-create-change:
	@echo "==============Submitting change request"
	touch .env
	TEST_SLACK=${TEST_SLACK	} 
	docker-compose -f 'docker-compose.wrapper.yml' run --rm go-spark-change
	