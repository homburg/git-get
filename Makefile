TAG=$(shell date +%Y%m%d%H%M)

release:
	git tag -a v1.2.$(TAG) -m "..."
	git push --follow-tags