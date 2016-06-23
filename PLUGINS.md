Deckard Plugins
============

**All information and requirements about plugins should be specified here**

AWS Credentials are configured with env `AWS_SECRET_ACCESS_KEY` and `AWS_ACCESS_KEY_ID`

| Plugin        | Command                    | Requirements                                                                                                                                       |
| ------------- | ------------------------   | ------------------------------------------------------------------------------------------------------------------|
| Cats          | `!cat`                     | None |
| Dice          | `!dice`                    | None |
| Tableflip     | `!tableflip` `!tablechill` | None |
| Write         | `!write`                   | Plugin settings: <ul><li>`HandwritingAPIURL="url with authentication"`</li><li>`S3Bucket="s3 bucket for storing images"`</li><li>AWS Credentials with access to `S3Bucket`</li></ul> |
| Principles    | `!principle`               | None |
