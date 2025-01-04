## simple ldc api in golang
- no fancy param parsing yet 
- limit and offset good to go
- no dates
- no operators
- no post routes

### to use, include config.yaml at root level with this structure:
```yaml
server:
  port: 8080

database:
  host: x.x.x.x
  port: xxxx
  name: dbname
  tenants:
    public:
      user: user1
      password: xxx
    publication:
      user: user2
      password: xxx
    legal:
      user: user3
      password: xxx

awsCognito:
  userPoolId: x
  clientId: x
  tokenType: id 
```