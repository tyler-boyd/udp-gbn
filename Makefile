docker: Dockerfile
	docker build . -t networks

run_docker: docker
	docker run --rm -v$(shell pwd):/go/src/github.com/tyler-boyd/udp-gbn -it networks

run_both:
	# Data: 3333
	# Acks: 3334
	go run . -mode recv localhost 3333 3334 out.txt > recv.txt &
	sleep 1
	go run . -mode send localhost 3333 3334 f.txt > send.txt

run_emu:
	./nEmulator-linux386 9991 localhost 9994 9993 localhost 9992 1 0 0 &
	sleep 1
	go run . -mode recv localhost 9994 9993 out.txt &
	sleep 1
	go run . -mode send localhost 9991 9992 f.txt
	pkill nEmu
	echo "Success"

clean:
	pkill main || true; pkill nEmu || true
