# System Design

Main service
- [apis](./api.md)
- many stateless replicas

User service
- oauth integration
- login/logout
- user creation, modification, etc.

Database service
- abstractions for different DBMS + database types
  - postgresql, mysql, sqlite

hash service
- generate short id for url from long url
  - id: group identifier (region), local identifier, random hash + incrementing counter
  - ex. cluster1-0-ijanek-0 or similar
  - base64 hash? nice for obfuscation, but not necessary maybe?
  revisit


```mermaid
flowchart LR
  user(User)
  main[Main Service]
  auth[Authentication Service]
  db[Database Service]
  dbProvider1[(PostgreSQL)]
  dbProvider2[(sqlite)]
  dbProvider3[(MySQL)]
  authProvider1[Github]
  authProvider2[Azure Entra]
  authProvider3[Any OAuth2 Provider]
  user --> main
  main --> db
  main --> auth
  db --> auth
  db -.- dbProvider1
  db -.- dbProvider2
  db -.- dbProvider3
  auth -.- authProvider1
  auth -.- authProvider2
  auth -.- authProvider3

```
