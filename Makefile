PANDOC = $(shell which pandoc)
MAN_FILES = $(wildcard man/*.md)
MAN_PAGES = $(patsubst man/%.md,man/%,$(MAN_FILES))

man/%: man/%.md
	$(PANDOC) -s -t man $< -o $@

all: $(MAN_PAGES)
