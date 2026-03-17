.PHONY: build test bench clean

build:
	@mkdir -p bin
	docker build --target builder -t wdxtools-builder .
	docker create --name wdxtools-extract wdxtools-builder 2>/dev/null || true
	docker cp wdxtools-extract:/out/numcrn ./bin/numcrn
	docker rm wdxtools-extract
	@echo "Built: ./bin/numcrn"

test:
	docker build --target tester -t wdxtools-tester .

bench:
	docker build --target bencher -t wdxtools-bencher .

clean:
	rm -rf bin/
	docker rmi wdxtools-builder wdxtools-tester wdxtools-bencher 2>/dev/null || true
