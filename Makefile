BINARY = arc
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
	install -m 644 docs/arc.1 $(PREFIX)/share/man/man1/arc.1
	install -d $(PREFIX)/share/zsh/site-functions
	install -m 644 docs/completions/arc.zsh $(PREFIX)/share/zsh/site-functions/_arc
	install -d $(PREFIX)/share/bash-completion/completions
	install -m 644 docs/completions/arc.bash $(PREFIX)/share/bash-completion/completions/arc
	install -d $(PREFIX)/share/fish/vendor_completions.d
	install -m 644 docs/completions/arc.fish $(PREFIX)/share/fish/vendor_completions.d/arc.fish

uninstall:
	rm -f $(PREFIX)/bin/$(BINARY)
	rm -f $(PREFIX)/share/man/man1/arc.1
	rm -f $(PREFIX)/share/zsh/site-functions/_arc
	rm -f $(PREFIX)/share/bash-completion/completions/arc
	rm -f $(PREFIX)/share/fish/vendor_completions.d/arc.fish

clean:
	rm -f $(BINARY)
	rm -rf docs/*.1 docs/completions/
