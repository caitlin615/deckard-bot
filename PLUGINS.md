Deckard Plugins
============

**All information and requirements about plugins should be specified here**

AWS Credentials are configured with env `AWS_SECRET_ACCESS_KEY` and `AWS_ACCESS_KEY_ID`

| Plugin        | Command                    | Requirements                                                                                                                                       |
| ------------- | ------------------------   | ------------------------------------------------------------------------------------------------------------------|
| Cats          | `!cat`                     | None |
| Dice          | `!dice`                    | None |
| Tableflip     | `!tableflip` `!tablechill` | None |
| Write         | `!write`                   | <ul><li>env `HANDWRITINGIO_API_URL="url with authentication"`</li><li>env `S3_BUCKET="s3 bucket for storing images"`</li><li>AWS Credentials with access to `S3_BUCKET`</li></ul> |
| Principles    | `!principle`               | None |
