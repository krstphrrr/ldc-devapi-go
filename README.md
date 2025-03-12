## simple ldc api in golang
- [x] operators work
- [x] operators also used for dates (FormDate, DateVisited)
- [x] POST request handling with operators
- [x] limit/offset work
- [x] aws cognito group discrimination works
- [x] containerized 
- [x] CICD-ready
- [ ] unified logs with levels across the app
- [x] exceptions to the routes (like tblProject)
- [x] add aero data handling

### upticks
v1.0.8
- added user + endpoint tracking for prometheus

v1.0.7
- fixes schema name duplication on POST request
- added a root path to display last updated on and app version number.

v1.0.6
- allowed enum on `tblProject`

v1.0.5
- enabled cors middleware on main

v1.0.4 
- geoindicators static schema at `/internal/schemas/geoIndicators.go` still referenced `"Precipitation_Yearly_Maximum_Daily_2_YR"`

v1.0.3 
- geoindicators static schema at `/internal/schemas/geoIndicators.go` still referenced `wkb_geometry`

v1.0.2
- removed columns and columntypes db request

v1.0.1
- exception to middleware routing
- aerodata handling

### Manual building
- update manually: `internal/version/version.go`, record uptick in `README.md`, and update docker image version on `docker-compose.yml` file. App version should correspond to docker image version.
- depending on host environment, export build date environmental variable:
```sh
export BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
```
- run `docker compose build` to build the image.
- deploy the stack: `docker stack deploy -c docker-compose.yml app-stack`. ensure that an external `ldc-go-net` attachable overlay network exists. (`docker network create ldc-go-net --attachable -d overlay`)


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