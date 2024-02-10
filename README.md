# env

Manage environment variables for the whole team safely

## How to set up

Initialize

```sh
go run ./cmd/init/
```

Get profile key

```sh
go run ./cmd/get/ [profile name]
```

## Files

- `/config.json` configuration for all profiles. Encrypted using master password.
- `/main.json` contains keys for all profiles. Encrypted using master password.
- `/profiles/[profile name].env` env file containing profiles variables.

## How to request

1. HTTP Get request to the .env file hosted by GitHub Pages
2. Decrypt the response using the profiles key with AES-256-CFB
3. You get a enc file (key value pairs, seperated by `=`)

### Golang

```go
package main

import (
    "github.com/Frank-Mayer/env"
)

func main() {
    err := env.Import(
        "https://Frank-Mayer.github.io/env/dev.env",
        "GTK7b9zi8Xp482/mQqo/FGDbEiI16XWKnpvfA8KCpPU=",
    )
    if err != nil {
        panic(err)
    }

    // environment variables for dev now available
}
```

### Node.js

```ts
const fs = require("fs");
const crypto = require("crypto");

await importEnv(
  "https://Frank-Mayer.github.io/env/dev.env",
  "GTK7b9zi8Xp482/mQqo/FGDbEiI16XWKnpvfA8KCpPU=",
);

// Function to decrypt AES encrypted file
async function importEnv(url: string, key: string) {
  // Convert base64 key to Buffer
  const binKey = Buffer.from(key, "base64");

  // Check key length
  if (binKey.length !== 32) {
    throw new Error("Invalid key length. Key must be 256 bits (32 bytes).");
  }

  // Reconvert the key to base64
  const base64Key = binKey.toString("base64");
  if (base64Key !== key) {
    throw new Error("Invalid key encoding. Key must be base64 encoded.");
  }

  // Read the encrypted file
  const resp = await fetch(url);
  if (!resp.ok) {
    throw new Error(`Failed to fetch the file: ${resp.statusText}`);
  }
  const encryptedData = Buffer.from(await resp.arrayBuffer());

  // Extract the IV from the data
  const iv = encryptedData.subarray(0, 16); // Assuming the IV size is the same as AES block size (16 bytes)
  // Extract the actual ciphertext
  const ciphertext = encryptedData.subarray(16);

  // Initialize AES cipher with the provided key and IV
  const decipher = crypto.createDecipheriv("aes-256-cfb", binKey, iv);

  // Decrypt the data
  let decryptedData = decipher.update(ciphertext);
  decryptedData = Buffer.concat([decryptedData, decipher.final()]);

  // Convert the decrypted data to string
  const decryptedString = decryptedData.toString("utf-8");

  // Set the environment variables
  const lines = decryptedString.split("\n");
  for (const line of lines) {
    const i = line.indexOf("=");
    if (i > 0) {
      const key = line.substring(0, i);
      const value = line.substring(i + 1);
      console.debug("importing", key);
      process.env[key] = value;
    }
  }
}
```
