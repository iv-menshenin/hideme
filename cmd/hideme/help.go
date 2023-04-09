package main

const helpInformation = `HIDEME - stealth data transfer application

If you want to hide a file
	./hideme inject \
		--payload=./original.jpg \
        --carrier=./carrier.png \
        --out=./encoded.png

If you want to encrypt hidden file
    ./hideme inject \
        --payload=./original.jpg \
        --carrier=./carrier.png \
        --out=./encoded.png \
        --encode-key=./crypt-key.jpg \
        --aes-key=af012453af01245305f76a0005f76a00

For detail information about injecting type:
	./hideme inject --help

For extraction you can use
    ./hideme extract \
        --input=./encoded.png \
        --public=./rsa_key.pub \
        --decode-key=./crypt-key.jpg \
        --aes-key=af012453af01245305f76a0005f76a00

For detail information about extraction type:
	./hideme extract --help

You can also generate signing keys. To do this, type:
	./hideme keys
`
