## simple ldc api in golang
- [x] operators work
- [x] operators also used for dates (FormDate, DateVisited)
- [x] POST request handling with operators
- [x] limit/offset work
- [x] aws cognito group discrimination works
- [x] containerized 
- [x] CICD-ready
- [ ] unified logs with levels across the app
- [ ] exceptions to the routes (like tblProject)
- [ ] add aero data handling

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