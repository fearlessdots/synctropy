# Destination directory
ifeq ($(DESTDIR),)
	# If DESTDIR is not specified, use the default value.
	_destdir ?= /usr
else
	# If DESTDIR was specified, use its value.
	_destdir ?= ${DESTDIR}
endif

# Name of the binary to be created
BINARY_NAME=synctropy

# Path to the directory where the binary will be installed
INSTALL_PATH=${_destdir}/bin

# Documentation output directory
DOCS_OUT=./docs
DOCS_PATH=${_destdir}/share/doc/${BINARY_NAME}

# Autocompletion files
AUTOCOMPLETION_OUT=./autocompletion

BASH_AUTOCOMPLETION_FILE=synctropy
BASH_AUTOCOMPLETION_INSTALL=${_destdir}/share/bash-completion/completions

ZSH_AUTOCOMPLETION_FILE=_synctropy
ZSH_AUTOCOMPLETION_INSTALL=${_destdir}/share/zsh/site-functions

FISH_AUTOCOMPLETION_FILE=synctropy.fish
FISH_AUTOCOMPLETION_INSTALL=${_destdir}/share/fish/vendor_completions.d

# License
LICENSE_PATH=${_destdir}/share/licenses/${BINARY_NAME}

# Flags to pass to the go build command
GO_BUILD_FLAGS=-v
GO_BUILD_LDFLAGS=-w

.DEFAULT_GOAL := build

.PHONY: clean
clean:
	@echo "====> Removing binary"
	if [ -f ${BINARY_NAME} ]; then \
		rm ${BINARY_NAME}; \
	fi
	@echo "====> Removing autocompletion files"
	if [ -d ${AUTOCOMPLETION_OUT} ]; then \
		rm -rf ${AUTOCOMPLETION_OUT}; \
	fi
	@echo "====> Removing documentation"
	if [ -d ${DOCS_OUT} ]; then \
		rm -rf ${DOCS_OUT}; \
	fi

.PHONY: deps
deps:
	@echo "====> Installing dependencies"
	go get -v

.PHONY: build
build: deps
	@echo "====> Building binary"
	go build ${GO_BUILD_FLAGS} -ldflags "$(GO_BUILD_LDFLAGS)"
	strip ${BINARY_NAME}

	mkdir -p ${AUTOCOMPLETION_OUT}
	@echo "====> Building autocompletion file for Bash"
	./${BINARY_NAME} completion bash > ${AUTOCOMPLETION_OUT}/${BASH_AUTOCOMPLETION_FILE}
	@echo "====> Building autocompletion file for Zsh"
	./${BINARY_NAME} completion zsh > ${AUTOCOMPLETION_OUT}/${ZSH_AUTOCOMPLETION_FILE}
	@echo "====> Building autocompletion file for Fish"
	./${BINARY_NAME} completion fish > ${AUTOCOMPLETION_OUT}/${FISH_AUTOCOMPLETION_FILE}

	@echo "====> Building documentation"
	./${BINARY_NAME} docs generate -o ${DOCS_OUT}

.PHONY: install
install:
	# Binary
	@echo "====> Installing binary"
	mkdir -p ${INSTALL_PATH}
	cp ${BINARY_NAME} ${INSTALL_PATH}
	# Autocompletion
	@echo "====> Installing autocompletion files"
	if [ -d "${BASH_AUTOCOMPLETION_INSTALL}" ]; then \
		cp ${AUTOCOMPLETION_OUT}/${BASH_AUTOCOMPLETION_FILE} ${BASH_AUTOCOMPLETION_INSTALL}/; \
	fi
	if [ -d "${ZSH_AUTOCOMPLETION_INSTALL}" ]; then \
		cp ${AUTOCOMPLETION_OUT}/${ZSH_AUTOCOMPLETION_FILE} ${ZSH_AUTOCOMPLETION_INSTALL}/; \
	fi
	if [ -d "${FISH_AUTOCOMPLETION_INSTALL}" ]; then \
		cp ${AUTOCOMPLETION_OUT}/${FISH_AUTOCOMPLETION_FILE} ${FISH_AUTOCOMPLETION_INSTALL}/; \
	fi
	# Documentation
	@echo "====> Installing documentation"
	mkdir -p "${DOCS_PATH}"
	cp ${DOCS_OUT}/* ${DOCS_PATH}/
	cp ./README.md ${DOCS_PATH}/
	# License
	@echo "====> Installing license"
	mkdir -p "${LICENSE_PATH}"
	cp ./LICENSE ${LICENSE_PATH}/

.PHONY: uninstall
uninstall:
	# Binary
	@echo "====> Uninstalling binary"
	rm ${INSTALL_PATH}/${BINARY_NAME}
	# Autocompletion
	@echo "====> Uninstalling autocompletion files"
	if [ -f "${BASH_AUTOCOMPLETION_INSTALL}/${BASH_AUTOCOMPLETION_FILE}" ]; then \
		rm -f "${BASH_AUTOCOMPLETION_INSTALL}/${BASH_AUTOCOMPLETION_FILE}"; \
	fi
	if [ -f "${ZSH_AUTOCOMPLETION_INSTALL}/${ZSH_AUTOCOMPLETION_FILE}" ]; then \
		rm -f "${ZSH_AUTOCOMPLETION_INSTALL}/${ZSH_AUTOCOMPLETION_FILE}"; \
	fi
	if [ -f "${FISH_AUTOCOMPLETION_INSTALL}/${FISH_AUTOCOMPLETION_FILE}" ]; then \
		rm -f "${FISH_AUTOCOMPLETION_INSTALL}/${FISH_AUTOCOMPLETION_FILE}"; \
	fi
	# Documentation
	@echo "====> Uninstalling documentation"
	rm -rf ${DOCS_PATH}
	# License
	@echo "====> Uninstalling license"
	rm -rf "${LICENSE_PATH}"

##
### RELEASES
##

RELEASES_OS = linux darwin freebsd openbsd netbsd
GOX_OUTPUT_PATH = ./gox_output

.PHONY: releases_clean
releases_clean:
	@echo "Removing binaries generated by gox"
	rm -rf $(GOX_OUTPUT_PATH)

.PHONY: releases_build
releases_build:
	mkdir -p $(GOX_OUTPUT_PATH)
	# Build binaries
	@echo "Building binaries with gox"
	@echo "Output: $(GOX_OUTPUT_PATH)"
	@for GOOS in $(RELEASES_OS); do \
		echo ""; \
		echo "Building binaries for OS: $$GOOS"; \
		gox -os=$$GOOS -ldflags="-w" -output "$(GOX_OUTPUT_PATH)"'/{{.Dir}}_{{.OS}}_{{.Arch}}'; \
	done
