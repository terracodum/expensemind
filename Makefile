install:
	chmod +x $(PWD)/expensemind
	ln -sf $(PWD)/expensemind /usr/local/bin/expensemind

uninstall:
	rm -f /usr/local/bin/expensemind
