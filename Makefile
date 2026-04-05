.PHONY: build test bench bench-php clean

build:
	@mkdir -p bin
	docker build --target builder -t wdxtools-builder .
	docker create --name wdxtools-extract wdxtools-builder 2>/dev/null || true
	docker cp wdxtools-extract:/out/wdxtools ./bin/wdxtools
	docker rm wdxtools-extract
	@echo "Built: ./bin/wdxtools"

test:
	@docker build --target builder -t wdxtools-builder -q .
	@docker run --rm wdxtools-builder go test ./... -v

bench:
	@docker build --target builder -t wdxtools-builder -q .
	@docker run --rm wdxtools-builder go test ./... -bench=. -benchmem

bench-php:
	docker build -t wdxtools-php-bench bench/php/
	docker run --rm wdxtools-php-bench

clean:
	rm -rf bin/
	docker rmi wdxtools-builder 2>/dev/null || true
