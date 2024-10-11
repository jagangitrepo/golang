all = proff
.PHONY = all

proff:
	go build generators.go && generators.exe -cpuprofile=generators.prof && go tool pprof generators.exe generators.prof

bench:
	go test generators_test.go -bench=. -benchmem 