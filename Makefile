TAG=$(shell date +%Y%m%d%H%M)

release:
	git tag -a v1.1$(TAG) -m "..."
	git push