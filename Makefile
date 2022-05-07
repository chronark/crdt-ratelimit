
protoc:
	protoc \
    	-I proto \
    	-I vendor/protoc-gen-validate \
    	-o /dev/null \
    	$(find proto -name '*.proto')
