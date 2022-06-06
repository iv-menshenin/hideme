# hideme - stealth data transfer application

Some time ago I read an interesting article about hidden data transfer. The author said that you could take a binary file and put every bit of it the lower bit of every byte of a PNG file.
"Okay" - I thought, - "But what if someone tries to read it? What if someone tries to replace that data with their own? How can I protect my data?".
So I decided to expand on that idea. And here is the result.

## data transfer

I want to transfer a file to my friend. I can only use public channels. Maybe file storage, maybe unsecured email, whatever.
I have to take the PNG image and my payload file and pass them as parameters to the `inject` command.

```shell
    ./hideme inject \
        --payload=./original.jpg \
        --carrier=./carrier.png \
        --out=./encoded.png
```

## data encryption

I want to encrypt my data, I have to protect my data.

Let's say that I have agreed in advance with my friend to use a particular encryption key.
I can encrypt data with AES by specifying the `aes-key` parameter.
```shell
    ./hideme inject \
        --payload=./original.jpg \
        --carrier=./carrier.png \
        --out=./encoded.png \
        --aes-key=af012453af01245305f76a0005f76a00
```

I can encrypt it with a key that is equal to (or greater) than length of the original message.
Suppose I make an agreement with a friend of mine that at a certain time of a certain day I will send him a file over the Internet.
This file will be the encryption key. Some time later, on a different day, I send my encrypted message.
I can encrypt my message by specifying the `encode-key` parameter. 
```shell
    ./hideme inject \
        --payload=./original.jpg \
        --carrier=./carrier.png \
        --out=./encoded.png \
        --encode-key=./crypt-key.jpg
```

The tool supports double encryption. I can use both approaches.
```shell
    ./hideme inject \
        --payload=./original.jpg \
        --carrier=./carrier.png \
        --out=./encoded.png \
        --encode-key=./crypt-key.jpg \
        --aes-key=af012453af01245305f76a0005f76a00
```

## digital signature

Let's imagine that my friend needs to ensure that the decoded data is not edited by the man in the middle (MITM).
I can sign my payload with private async key and give my friend the public key so that he can verify the signature.
Now, if the MITM is able to change the message (let's assume that he revealed our keys).
Then without knowing my private key he will not be able to fool my friend.
```shell
    ./hideme keys &&
    ./hideme inject \
        --payload=./original.jpg \
        --carrier=./carrier.png \
        --out=./encoded.png \
        --private=./rsa_key
```

## file extraction

How can my friend get the hidden file and verify its digital signature?
He just has to specify a public key for the program to perform the verification of digital signature.
```shell
    ./hideme extract \
        --input=./encoded.png \
        --public=./rsa_key.pub \
        --decode-key=./crypt-key.jpg \
        --aes-key=af012453af01245305f76a0005f76a00

    2022/06/03 16:36:01 the signature is verified well
```
If you see this message at the bottom (I mean "the signature is verified well"), then the digital signature has been successfully verified.
Otherwise, the digital signature is not valid. Information about the presence of a digital signature in case of failure is not disclosed for security reasons.

## help

You can get help with application launch options. To do this, simply type:
```shell
    ./hideme help
```

Well, in order to help the author of the project. Just learn programming, learn new techniques, practice making your applications stable and secure.
We have so much bad code so far