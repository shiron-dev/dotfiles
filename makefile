.PHONY: sops-encrypt sops-decrypt sops-ci kics

sops-encrypt:
	@echo "Encrypting with SOPS..."; \
	if [ -n "$(FILE)" ]; then \
		if [ -f "$(FILE)" ] && [ "$${FILE##*.}" != "sops" ]; then \
			FILES="$(FILE)"; \
		elif [ -f "$(FILE)" ] && [ "$${FILE##*.}" = "sops" ]; then \
			base="$${FILE%.sops}"; \
			if [ -f "$$base" ]; then \
				FILES="$$base"; \
			else \
				echo "Error: plaintext $$base not found for $(FILE)" >&2; \
				exit 1; \
			fi; \
		elif [ -f "$(FILE).sops" ]; then \
			base="$(FILE)"; \
			if [ -f "$$base" ]; then \
				FILES="$$base"; \
			else \
				echo "Error: plaintext $$base not found (got $(FILE).sops)" >&2; \
				exit 1; \
			fi; \
		else \
			echo "Error: $(FILE) not found" >&2; \
			exit 1; \
		fi; \
	else \
		FILES="$$(find . -name "*.secrets.*" -type f ! -name "*.sops")"; \
	fi; \
	for file in $$FILES; do \
		echo "Encrypting $$file..."; \
		sops --output-type json --encrypt "$$file" > "$$file.sops"; \
	done

sops-decrypt:
	@echo "Decrypting with SOPS..."; \
	if [ -n "$(FILE)" ]; then \
		if [ -f "$(FILE)" ]; then \
			FILES="$(FILE)"; \
		elif [ -f "$(FILE).sops" ]; then \
			FILES="$(FILE).sops"; \
		else \
			echo "Error: $(FILE) or $(FILE).sops not found" >&2; \
			exit 1; \
		fi; \
	else \
		FILES="$$(find . -name "*.secrets.*.sops" -type f)"; \
	fi; \
	for file in $$FILES; do \
		echo "Decrypting $$file..."; \
		base="$${file%.sops}"; \
		ext="$${base##*.}"; \
		case "$$ext" in \
			yaml|yml) output_type="yaml" ;; \
			*) output_type="binary" ;; \
		esac; \
		if [ -f "$$base" ]; then \
			chmod +w "$$base"; \
		fi; \
		sops --decrypt --output-type "$$output_type" "$$file" > "$$base"; \
		chmod -w "$$base"; \
	done

sops-ci:
	@echo "Checking for unencrypted secrets tracked by git..."; \
	FILES="$$(find . -name '*.secrets.*' ! -name '*.secrets.*.sops' -type f)"; \
	EXIT=0; \
	for file in $$FILES; do \
		if git ls-files --error-unmatch "$$file" >/dev/null 2>&1; then \
			echo "Error: Unencrypted secrets file tracked by git: $$file" >&2; \
			EXIT=1; \
		fi; \
	done; \
	if [ $$EXIT -ne 0 ]; then \
		echo "One or more unencrypted secrets files are tracked by git. Please remove them from version control." >&2; \
		exit 1; \
	fi

kics:
	docker run -t -v $(PWD):/path checkmarx/kics:latest scan -p /path
