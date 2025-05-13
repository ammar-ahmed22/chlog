.PHONY: build install

build:
	@cd template-website && npm run build
	@echo "Copying dist to tpl/website..."
	@echo "Cleaning tpl/website..."
	@rm -rf tpl/website
	@mkdir -p tpl/website
	@cp -r template-website/dist/. tpl/website/
	@echo "Copied dist to tpl/website."


install: build
	@echo "Installing chlog..."
	@go install ./...
	@echo "Installing chlog complete."

