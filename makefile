COVERPKG=$(shell go list ./... | grep -v '/mock$$' | paste -sd "," -)

test-coverage:
	@echo "Covering packages: $(COVERPKG)"
	go test -v -coverpkg="go test -v -coverpkg=$(COVERPKG) -coverprofile=profile.cov ./..."
	go tool cover -func=profile.cov
