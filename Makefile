BINARY = aurl
PREFIX ?= /usr/local

.PHONY: build docs install uninstall clean

build:
	go build -ldflags "-s -w" -o $(BINARY) .

docs:
	go run cmd/gendocs/main.go

install: build docs
	install -d $(PREFIX)/bin
	install -m 755 $(BINARY) $(PREFIX)/bin/$(BINARY)
	install -d $(PREFIX)/share/man/man1
	install -m 644 docs/aurl.1 $(PREFIX)/share/man/man1/aurl.1
	install -d $(PREFIX)/share/zsh/site-functions
	install -m 644 docs/completions/aurl.zsh $(PREFIX)/share/zsh/site-functions/_aurl
	install -d $(PREFIX)/share/bash-completion/completions
	install -m 644 docs/completions/aurl.bash $(PREFIX)/share/bash-completion/completions/aurl
	install -d $(PREFIX)/share/fish/vendor_completions.d
	install -m 644 docs/completions/aurl.fish $(PREFIX)/share/fish/vendor_completions.d/aurl.fish

uninstall:
	rm -f $(PREFIX)/bin/$(BINARY)
	rm -f $(PREFIX)/share/man/man1/aurl.1
	rm -f $(PREFIX)/share/zsh/site-functions/_aurl
	rm -f $(PREFIX)/share/bash-completion/completions/aurl
	rm -f $(PREFIX)/share/fish/vendor_completions.d/aurl.fish

clean:
	rm -f $(BINARY)
	rm -rf docs/*.1 docs/completions/
