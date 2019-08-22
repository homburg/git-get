SHELL_TEST=$(shell echo $$SHELL)
test:
	env
	echo $(SHELL)
	echo "TEST: " $(SHELL_TEST)
