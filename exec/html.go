package exec

const htmlPage = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>steganography</title>
	</head>
	<body style="margin: auto; padding: auto; width: 100%; height: 100%">
		<div style="width: 450px; margin: auto; font-family: Courier New, monospace">
			<p>Injecting:</p> <hr/>
			<form action="/inject" method="post" enctype="multipart/form-data" target="_blank" style="border: 1px solid gray; padding: 20px">
				Carrier: <input type="file" name="carrier" /><br />
				Payload: <input type="file" name="payload" /><br />
				Private: <input type="file" name="private" /><br />
				SyncKey: <input type="file" name="encode-key" /><br />
				AES key: <input type="text" name="aes" /><br />
				<input type="submit" value="INJECT">
			</form>
			<p>Extracting:</p> <hr/>
			<form action="/extract" method="post" enctype="multipart/form-data" target="_blank" style="border: 1px solid gray; padding: 20px">
				Carrier: <input type="file" name="carrier" /><br />
				Public:&nbsp;&nbsp;<input type="file" name="public" /><br />
				SyncKey: <input type="file" name="encode-key" /><br />
				AES key: <input type="text" name="aes" /><br />
				<input type="submit" value="EXTRACT">
			</form>
			<p>Generate RSA keys:</p> <hr/>
			<form action="/generate" method="post" enctype="multipart/form-data" target="_blank" style="border: 1px solid gray; padding: 20px">
				<input type="submit" value="GENERATE">
			</form>
		</div>
	</body>
</html>
`
