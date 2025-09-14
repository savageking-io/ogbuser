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
