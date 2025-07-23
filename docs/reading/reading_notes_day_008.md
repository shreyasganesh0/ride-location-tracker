# JWT

## JSON Web Token
- introduced in RFC 7519
- way to transfer information securly as JSON object
- verified using signed secret (HMAC) or private/public key encryption (ECDSA or RSA)
- Uses
    - Authz 
        - user logs in
        - include jwt in each subsequent req
        - SSO uses JWT, stateless and small overhead
    - Info xchange
        - JWTS can be signed
        - can verify senders because of this
        - signatures use both content and payload
            - can verify tampering
- Structure
    - xxxx.yyyy.zzzz
    - Header
        - token type "typ" : "JWT"
        - signing algo "alg": "HS256"
        - base64url encoded
    - Payload
        - contains claims and additional data 
            - claim id are 3 chars long
        - 3 types
            - Registered claims:
                - predefined optional claims
                - iss(issuer)
                - exp(expiration time)
                - sub(subject)
            - Public claims:
                - recommended to be defined in the IANA JSON Web Token Registry
                - should be collision resistent
            - Private claims:
                - custom claims
                - used to share info between parties
        - base64url encoded
    - Signature
        - encoded header + encoded payload + secret algo and sign
- output is 3 basic base64url encoded string
- server generates JWT and returns to the user
- short lived since they are creds
- passed in the header of HTTP request Authorization: Bearer <token>
    - if used as stateless authz mechanism must not be greater then 8kb 
      since some servers will reject the header
    - use Auth0 Fine Grained Authorization if need to include info like
      user permissions
- Validation vs Verification
    - JWT validation
        - check structure and format and content of JWT
        - ensure it has the 3 parts header, payload and signature
        - format: veify encoding of each part and expected claims in payload
        - content: check if claims values are corred like expiration time(exp) 
          issued at(iat) and not before(nbf)
    - JWT verification
        - authenticy and integrity of token
        - Signature verification: verify the signature against header and payload
        - Issuer verification: check if claims match issuer
        - Audience check: aud claim matches audience

## Go JWT
- jwt.ParseRSAPrivateKeyFromPEM(signBytes);
    - uses the bytes of the private key to get the JWT sign key
- jwt.ParseRSAPublicKeyFromPEM(verifyBytes);
    - get the public key fromt the public key file bytes
    - use it to verify JWT
- using RSA assym auth to sign the jwt
- jwt.ParseWithClaims
