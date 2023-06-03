.PHONY : wat-dev
wat-dev :
	while true; do inotifywait -e close_write ./wat/core.wat ; clear ; wat2wasm ./wat/core.wat ; done
