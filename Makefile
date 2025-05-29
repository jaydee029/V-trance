build:
	@cd user/cmd && go build -o user
	@cd trance-api/cmd && go build -o tranceapi
	@cd worker/cmd && go build -o worker

run: build
	@./scripts/rabbit.sh start
	@bash -c 'cd user/cmd && ./user &' \
	&& sleep 7 \
	&& bash -c 'cd trance-api/cmd && ./tranceapi &' \
	&& sleep 5 \
	&& bash -c 'cd worker/cmd && ./worker &'
	@echo "All services started with delays."

clean:
	@./scripts/rabbit.sh stop 
	@pkill -f ./tranceapi
	@pkill -f ./user
	@pkill -f ./worker
	echo "services cleaned"
