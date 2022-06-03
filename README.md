# hideme - stealth data transfer application

Some time ago I read an interesting article about hidden data transfer. The author said that you could take a binary file and put every bit of it the lower bit of every byte of a PNG file.
"Okay" - I thought, - "But what if someone tries to read it? What if someone tries to replace that data with their own? How can I protect my data?".
So I decided to expand on that idea. And here is the result.

## data transfer

I want to transfer a file to my friend. I can only use public channels. Maybe file storage, maybe unsecured email, whatever. I can ...

```sh
    ./hideme inject --payload=./original.jpg --carrier=./carrier.png --out=./encoded.png
```

## data encryption

I want to encrypt my data, I have to protect my data

I can encrypt it with AES
```sh
./hideme inject --payload=./original.jpg --carrier=./carrier.png --out=./encoded.png --aes-key=af012453af01245305f76a0005f76a00
```

I can encrypt it with a key that is equal to (or greater than) the length of the original message
```sh
./hideme inject --payload=./original.jpg --carrier=./carrier.png --out=./encoded.png --encode-key=./crypt-key.jpg
```

Or both
```sh
./hideme inject --payload=./original.jpg --carrier=./carrier.png --out=./encoded.png --encode-key=./crypt-key.jpg --aes-key=af012453af01245305f76a0005f76a00
```

## digital signature

```sh
./hideme keys
./hideme inject --payload=./original.jpg --carrier=./carrier.png --out=./encoded.png --private=./rsa_key
```
