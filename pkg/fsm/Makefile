gen:
	docker run --rm -it \
  		-w ${PWD} \
	  	-v ${PWD}:${PWD} \
	  	codemityio/type:latest \
	  	quicktype \
	  	-s schema dsl/schema.json \
	  	-o dsl/workflow_gen.go \
	  	-t Workflow \
	  	-l go \
  		--package dsl \
  		--just-types --just-types-and-package
