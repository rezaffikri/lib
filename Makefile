update:
	git tag v$(tag)
	git push origin v$(tag)

version:
	git describe --tags --abbrev=0