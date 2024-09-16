# uaausersimport


## Problem
When cloudfoundry integrates with external userstores (E.g. ldap/saml). Currently it does not have any way to assign user roles (E.g. org and space). So users have to login to the CF env first, then operator assigns the roles to them. When they login again, users can see the org and spaces.

This may not be applicable for the invitation model for operators, since they have to ask users to login first

## Resolution

Prepopulate the users with uaa api and cloudcontroller api. So they can have all the roles before user logins.

Prerequisites:

Use uaac to create a client id, who has cloudcontroller.admin and scim.write

* Target your UAA
  ```
  uaac target https://uaa.<systemdomain>/
  ```
* Login using Admin Client Credentials
  The Admin Client secret can be retrieved from Ops Manager. Navigate to Ops Manager -> Pivotal Elastic Runtime -> Credentials -> Admin Client Credentials
  ```
  uaac token client get admin -s <secret>
  ```
* add a new client

  * uaac client add -i
  * Client name: bulkimport
  * New client secret: <secret>
  * Verify new client secret: <secret>
  * scope (list): Press Enter
  * authorized grant types (list):  client_credentials
  * authorities (list):  cloud_controller.admin,scim.write
  * access token validity (seconds):  Press Enter
  * refresh token validity (seconds): Press Enter
  * redirect uri (list): Press Enter
  * autoapprove (list): Press Enter
  * signup redirect url (url):  Press Enter


Steps (What this progam is doing?):

1. Get token from uaa
2. Add user to the uaa
   * sample users yaml file: [sample file](config/fixtures/users.yml)
   * Sample user config

   ```
   origin: ldap
   - uid: jcalabrese@pivotal.io
     externalid: uid=jcalabrese,ou=People,dc=homelab,dc=io
     emails:
     - jcalabrese@pivotal.io
     orgs:
       - name: org1
         roles:
         - OrgManager
         - OrgAuditor
         spaces:
           - name: space1
             roles:
             - SpaceManager
             - SpaceDeveloper
             - SpaceAuditor
           - name: space2
             roles:
             - SpaceManager
             - SpaceAuditor
       - name: org2
         roles:
         - OrgAuditor
         spaces:
           - name: space1
             roles:
             - SpaceAuditor
           - name: space2
             roles:
             - SpaceAuditor
   ```

3. Add user to the cloudcontroler
4. Associate user roles with the orgs
5. Associate user roles with the spaces

   * Functional Programming (In [main.go](main.go))

   ```
      token.GetToken.MapUsers().AddUAAUsers().AddCCUsers().MapOrgs().MapSpaces(ctx)
   ```

## How to run

* Install go

* Get the binary

```
go get -u github.com/pivotalservices/uaausersimport

```
* Target the cf environment

```
export CF_ENVIRONMENT=environment.yml (change to your environment.yml)
```

* Target the users file

```
export USERS_CONFIG_FILE=config/fixtures/users.yml (change to your user files)
```
* Enable http traffic dump, optional:

  ```
  export DEBUG_HTTP=true
  ```

* Run

  ```
  uaausersimport
  ```

## Future work

* Create an interface (web/command line) help client generate formatted file
* Cross compile the code
