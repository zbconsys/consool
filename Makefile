get-addresses:
	go run main.go -mode w3s-addresses -json-rpc https://rpc.devnet.linea.build -web3-signer http://localhost:9000

send-funds:
	go run main.go \
	-mode w3s-send \
	-json-rpc https://rpc.devnet.linea.build \
	-web3-signer http://localhost:9000 -public-key 0xbc968c24047c74d615eb4f0cd91037bd79e6349d8b1b26e850ae211fafc7886a985c6156a12511f490c10e39693dec077a830a1e9756c0644e2b5e3628360752 \
	-send-to-address 0x228466F2C715CbEC05dEAbfAc040ce3619d7CF0B
