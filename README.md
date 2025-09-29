# ogbuser



### Error Codes
| Code  | Error Text                   | Description                                                  |
|-------|------------------------------|--------------------------------------------------------------|
| 12000 | failed to parse request body | Malformed JSON received from REST service                    |
| 12001 | empty username               | Provided username/email were empty                           |
| 12002 | empty password               | Provided password was empty                                  |
| 12003 | user not found               | Username/password didn't match                               |
| 12004 | <dynamic>                    | Error occurred while loading user from DB                    |
| 12005 | <dynamic>                    | Failed to verify password due to encryption error            |
| 12006 | failed to initialize session | Failed to initialize session due to token generation problem |
| 12007 | failed to load groups        | Failed to load groups that a user is member of               |
| 12008 | failed to set permissions    | Failed to set permissions based on user's groups             |


### Permissions and Scopes
Each microservice defines their own scopes and user permissions. Globally
each permission has 3 access bits - Read, Write and Delete. Another important thing 
that defines permission is a domain - own, party, guild and global. 

Examples:
* OWN 1 READ 0 NOWRITE 0 NODELETE will allow user to read their own
data, but restrict writing (updating) or deleting.
* PARTY 1 READ 0 NOWRITE 0 NODELETE will be able to read data from an 
entire party