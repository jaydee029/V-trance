build:
	@cd user/cmd && go build -o user
	@cd trance-api/cmd && go build -o tranceapi
	@cd worker/cmd && go build -o worker

run: build
	@./scripts/rabbit.sh start && cd user/cmd && ./user & \
	cd trance-api/cmd && ./tranceapi & \
	cd worker/cmd && ./worker

clean:
	@./scripts/rabbit.sh stop
	@pkill -f ./user || true
	@pkill -f ./tranceapi || true
	@pkill -f ./worker || true
	@echo "Stopped all services."