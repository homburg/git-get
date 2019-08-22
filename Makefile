SHELL_TEST=$(shell echo $SHELL)
test:
	echo $(SHELL)
	echo "TEST: " $(SHELL_TEST)
