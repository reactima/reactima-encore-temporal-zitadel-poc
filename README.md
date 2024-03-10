# reactima-encore-temporal-zitadel-poc

This is a demo playground for several enterprise ready technologies

- [Encore](https://github.com/encoredev/encore)
- [Temporal](https://temporal.io/) workflow engine
- [Zitadel](https://zitadel.com/) IDP (Indentity Provider) with features include multi-tenancy with branding customization, secure login, self-service, OpenID Connect, OAuth2.x, SAML2, Passwordless with FIDO2 (including Passkeys), OTP, U2F, and an unlimited audit trail.

## 🧾 Part 1  - Billing Service Demo

**Billing Service Demo** section is designed to showcase temporal workflow application in context of billing transactions.

## 🛠 Development Setup

To dive into this demo locally:

### 🌀 Launching Temporal Server

1. Grab the latest [Temporal](https://learn.temporal.io/getting_started/go/dev_environment/).
  
2. Start it up with the command:
   ```bash
   temporal server start-dev --db-filename temporal.db
   ```

3. 🌍 To access the Temporal dashboard, navigate to:
   ```
   http://127.0.0.1:8233
   ```

### 🎵 Working with Encore

For local development, ensure you've got [Encore](https://encore.dev/docs) ready.

1. First of all run test via Encore:
   ```bash
   encore test ./...
   ```

2. Run the following to kick off Encore:
   ```bash
   encore run
   ```

3. 🌍 This spins up the Encore UI in your default browser at:
   ```
   http://localhost:9400
   ```

   Use the UI to monitor and analyze API calls seamlessly.

## 🧾 Part 2  - Zitadel SDK Demo

**Zitadel SDK Demo** section is designed to showcase the following actions

- user creation
- updating user by email
- deactivating user
- changing password
- changing avatar
- creating project
- assigning user by email to project

## 🚀 Zitadel

### Zitadel Cloud

Sign up for cloud version https://zitadel.com/

### Zitadel Local/Selfhosted

Go through
- https://zitadel.com/docs/self-hosting/deploy/overview
- https://github.com/zitadel/zitadel/blob/main/cmd/defaults.yaml

For maximum control over Zitadel run the following on Ubuntu 20/22 server

```bash
make compile
docker build -f Dockerfile -t zitadel2 .
```

#### Run Zitadel

```bash
./zitadel start-from-init --config reactima-config/example-zitadel-config.yaml --masterkey "MasterkeyNeedsToHave32Dummy##"
```

#### Protobuf/Protoset interface

Generate protoset file

```bash
buf build -o zitadel.protoset
```

#### gRPC client

For testing gRPC calls to Zitadel install [grpcui](https://github.com/fullstorydev/grpcui)

```bash
grpcui -vv -plaintext -H 'authorization: Bearer XXX_TOKEN_FROM_ZITADEL_HERE' -protoset  /path/zitadel.protoset  -plaintext localhost:8080
```

Replace "localhost:8080" with Zitadel address. Use "example.com:443" for https cloud base version


# Headhunt me

Interested to discuss the above or looking for a partner to work on AI Agents, Data Mining, CRM, B2B Lead Generation and/or Outbound Marketing SaaS system related projects? Feel free to ping me to exchange ideas or request a consultation!

- +81 80 4756 2323 (WhatsApp, preferred)
[api.whatsapp.com/send?phone=818047562323](https://api.whatsapp.com/send?phone=818047562323)

- [1@ilya1.com](mailto:1@ilya1.com)
- [t.me/reactima](https://t.me/reactima)
- [linkedin.com/in/reactima](https://www.linkedin.com/in/reactima/)

